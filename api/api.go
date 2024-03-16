package api

import (
	"encoding/json"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"strconv"
	"strings"
	"time"

	api_models "filmoteka/api/models"
	db_models "filmoteka/db"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-pg/pg/v10"
)

func StartAPI(pgdb *pg.DB) *chi.Mux {
	r := chi.NewRouter()
	r.Use(middleware.Logger, middleware.RequestID, middleware.Recoverer, middleware.WithValue("DB", pgdb))

	r.Route("/films", func(r chi.Router) {
		r.Get("/", getFilms)
		r.Post("/", createFilm)
		r.Put("/{filmID}", updateFilm)
		r.Delete("/{filmID}", deleteFilm)
	})
	r.Route("/actors", func(r chi.Router) {
		r.Get("/", getActors)
		r.Post("/", createActor)
		r.Put("/{actorID}", updateActor)
		r.Delete("/{actorID}", deleteActor)
	})

	r.Get("/healthcheck", func(w http.ResponseWriter, r *http.Request) {
		_, ok := r.Context().Value("DB").(*pg.DB)
		if !ok {
			msg := "could not get the DB from context"
			slog.Error(msg)
			w.Write([]byte(msg))
			return
		}
		w.Write([]byte("OK"))
		w.WriteHeader(http.StatusOK)
	})

	slog.Info("Success start API routes")
	return r
}

func getFilms(w http.ResponseWriter, r *http.Request) {
	// sortBy is expected to look like field.orderdirection i. e. id.asc
	sortBy := r.URL.Query().Get("sortBy")
	filter := r.URL.Query().Get("filter")
	var splits []string
	if filter != "" {
		splits = strings.Split(filter, ".")
	}
	if sortBy == "" || sortBy == "rate" {
		// rate ASC is the default sort query
		sortBy = "rate DESC"
	}
	//get the db from context
	pgdb, ok := r.Context().Value("DB").(*pg.DB)
	//if we can't get the db let's handle the error
	//and send an adequate response
	if !ok {
		res := &api_models.FilmsResponse{
			Success: false,
			Error:   "could not get the DB from context",
			Films:   nil,
		}
		err := json.NewEncoder(w).Encode(res)
		//if there's an error with encoding handle it
		if err != nil {
			log.Printf("error sending response %v\n", err)
		}
		//return a bad request and exist the function
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	//call models package to access the database and return the comments
	films, err := db_models.GetFilms(pgdb, sortBy, splits)
	if err != nil {
		res := &api_models.FilmsResponse{
			Success: false,
			Error:   err.Error(),
			Films:   nil,
		}
		log.Printf("error getting films %v\n", err)
		err := json.NewEncoder(w).Encode(res)
		if err != nil {
			log.Printf("error sending response %v\n", err)
		}
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	//positive response
	res := &api_models.FilmsResponse{
		Success: true,
		Error:   "",
		Films:   films,
	}
	//encode the positive response to json and send it back
	err = json.NewEncoder(w).Encode(res)
	if err != nil {
		log.Printf("error encoding films: %v\n", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func createFilm(w http.ResponseWriter, r *http.Request) {
	//get the request body and decode it
	req := &api_models.CreateFilmRequest{}
	err := json.NewDecoder(r.Body).Decode(req)
	//if there's an error with decoding the information
	//send a response with an error
	if err != nil {
		res := &api_models.FilmResponse{
			Success: false,
			Error:   err.Error(),
			Film:    nil,
		}
		err = json.NewEncoder(w).Encode(res)
		//if there's an error with encoding handle it
		if err != nil {
			log.Printf("error sending response %v\n", err)
		}
		//return a bad request and exist the function
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	//get the db from context
	pgdb, ok := r.Context().Value("DB").(*pg.DB)
	//if we can't get the db let's handle the error
	//and send an adequate response
	if !ok {
		res := &api_models.FilmResponse{
			Success: false,
			Error:   "could not get the DB from context",
			Film:    nil,
		}
		err := json.NewEncoder(w).Encode(res)
		//if there's an error with encoding handle it
		if err != nil {
			log.Printf("error sending response %v\n", err)
		}
		//return a bad request and exist the function
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	//if we can get the db then
	datetime, err := time.Parse("2006-01-02", req.Date)
	if err != nil {
		fmt.Println("Got datetime ", datetime)
	}
	film, err := db_models.CreateFilm(pgdb, &db_models.Film{
		Name:        req.Name,
		Description: req.Description,
		Date:        datetime,
		Rate:        req.Rate,
	}, req.Actors)
	if err != nil {
		res := &api_models.FilmResponse{
			Success: false,
			Error:   err.Error(),
			Film:    nil,
		}
		err = json.NewEncoder(w).Encode(res)
		//if there's an error with encoding handle it
		if err != nil {
			log.Printf("error sending response %v\n", err)
		}
		//return a bad request and exist the function
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	//everything is good
	//let's return a positive response
	res := &api_models.FilmResponse{
		Success: true,
		Error:   "",
		Film:    film,
	}
	err = json.NewEncoder(w).Encode(res)
	if err != nil {
		log.Printf("error encoding after creating comment %v\n", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func updateFilm(w http.ResponseWriter, r *http.Request) {
	//get the data from the request
	req := &api_models.UpdateFilmRequest{}
	//decode the data
	err := json.NewDecoder(r.Body).Decode(req)
	if err != nil {
		res := &api_models.FilmResponse{
			Success: false,
			Error:   err.Error(),
			Film:    nil,
		}
		err = json.NewEncoder(w).Encode(res)
		//if there's an error with encoding handle it
		if err != nil {
			log.Printf("error sending response %v\n", err)
		}
		//return a bad request and exist the function
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	pgdb, ok := r.Context().Value("DB").(*pg.DB)
	if !ok {
		res := &api_models.FilmResponse{
			Success: false,
			Error:   "could not get the DB from context",
			Film:    nil,
		}
		err := json.NewEncoder(w).Encode(res)
		//if there's an error with encoding handle it
		if err != nil {
			log.Printf("error sending response %v\n", err)
		}
		//return a bad request and exist the function
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	//get the commentID to know what comment to modify
	filmID := chi.URLParam(r, "filmID")
	//we get a string but we need to send an int so we convert it
	intFilmID, err := strconv.Atoi(filmID)
	datetime, err := time.Parse("2006-01-02", req.Date)
	if err != nil {
		res := &api_models.FilmResponse{
			Success: false,
			Error:   err.Error(),
			Film:    nil,
		}
		err = json.NewEncoder(w).Encode(res)
		//if there's an error with encoding handle it
		if err != nil {
			log.Printf("error sending response %v\n", err)
		}
		//return a bad request and exist the function
		// w.WriteHeader(http.StatusBadRequest)
	}

	//update the comment
	film, err := db_models.UpdateFilm(pgdb, &db_models.Film{
		ID:          intFilmID,
		Name:        req.Name,
		Description: req.Description,
		Date:        datetime,
		Rate:        req.Rate,
	})
	if err != nil {
		res := &api_models.FilmResponse{
			Success: false,
			Error:   err.Error(),
			Film:    nil,
		}
		err = json.NewEncoder(w).Encode(res)
		//if there's an error with encoding handle it
		if err != nil {
			log.Printf("error sending response %v\n", err)
		}
		//return a bad request and exist the function
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	//return successful response
	res := &api_models.FilmResponse{
		Success: true,
		Error:   "",
		Film:    film,
	}
	//send the encoded response to responsewriter
	err = json.NewEncoder(w).Encode(res)
	if err != nil {
		log.Printf("error encoding comments: %v\n", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	//send a 200 response
	w.WriteHeader(http.StatusOK)
}

func deleteFilm(w http.ResponseWriter, r *http.Request) {

	//get the db from ctx
	pgdb, ok := r.Context().Value("DB").(*pg.DB)
	if !ok {
		res := &api_models.FilmResponse{
			Success: false,
			Error:   "could not get the DB from context",
			Film:    nil,
		}
		err := json.NewEncoder(w).Encode(res)
		//if there's an error with encoding handle it
		if err != nil {
			log.Printf("error sending response %v\n", err)
		}
		//return a bad request and exist the function
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	//get the commentID
	filmID := chi.URLParam(r, "filmID")
	intFilmID, err := strconv.ParseInt(filmID, 10, 64)
	if err != nil {
		res := &api_models.FilmResponse{
			Success: false,
			Error:   err.Error(),
			Film:    nil,
		}
		err = json.NewEncoder(w).Encode(res)
		//if there's an error with encoding handle it
		if err != nil {
			log.Printf("error sending response %v\n", err)
		}
		//return a bad request and exist the function
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	//delete comment
	err = db_models.DeleteFilm(pgdb, intFilmID)
	if err != nil {
		res := &api_models.FilmResponse{
			Success: false,
			Error:   err.Error(),
			Film:    nil,
		}
		err = json.NewEncoder(w).Encode(res)
		//if there's an error with encoding handle it
		if err != nil {
			log.Printf("error sending response %v\n", err)
		}
		//return a bad request and exist the function
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	//send a 200 response
	w.WriteHeader(http.StatusOK)
}

func getActors(w http.ResponseWriter, r *http.Request) {
	//get the db from context
	pgdb, ok := r.Context().Value("DB").(*pg.DB)
	//if we can't get the db let's handle the error
	//and send an adequate response
	if !ok {
		res := &api_models.ActorsResponse{
			Success: false,
			Error:   "could not get the DB from context",
			Actors:  nil,
		}
		err := json.NewEncoder(w).Encode(res)
		//if there's an error with encoding handle it
		if err != nil {
			log.Printf("error sending response %v\n", err)
		}
		//return a bad request and exist the function
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	//call models package to access the database and return the comments
	actors, err := db_models.GetActors(pgdb)
	if err != nil {
		res := &api_models.ActorsResponse{
			Success: false,
			Error:   err.Error(),
			Actors:  nil,
		}
		log.Printf("error getting actors %v\n", err)
		err := json.NewEncoder(w).Encode(res)
		if err != nil {
			log.Printf("error sending response %v\n", err)
		}
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	//positive response
	res := &api_models.ActorsResponse{
		Success: true,
		Error:   "",
		Actors:  actors,
	}
	//encode the positive response to json and send it back
	err = json.NewEncoder(w).Encode(res)
	if err != nil {
		log.Printf("error encoding actors: %v\n", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func createActor(w http.ResponseWriter, r *http.Request) {
	//get the request body and decode it
	req := &api_models.CreateActorRequest{}
	err := json.NewDecoder(r.Body).Decode(req)
	//if there's an error with decoding the information
	//send a response with an error
	if err != nil {
		res := &api_models.ActorResponse{
			Success: false,
			Error:   err.Error(),
			Actor:   nil,
		}
		err = json.NewEncoder(w).Encode(res)
		//if there's an error with encoding handle it
		if err != nil {
			log.Printf("error sending response %v\n", err)
		}
		//return a bad request and exist the function
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	//get the db from context
	pgdb, ok := r.Context().Value("DB").(*pg.DB)
	//if we can't get the db let's handle the error
	//and send an adequate response
	if !ok {
		res := &api_models.ActorResponse{
			Success: false,
			Error:   "could not get the DB from context",
			Actor:   nil,
		}
		err := json.NewEncoder(w).Encode(res)
		//if there's an error with encoding handle it
		if err != nil {
			log.Printf("error sending response %v\n", err)
		}
		//return a bad request and exist the function
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	//if we can get the db then
	birthday, err := time.Parse("2006-01-02", req.Birth)
	if err != nil {
		fmt.Println("Got datetime ", birthday)
	}
	actor, err := db_models.CreateActor(pgdb, &db_models.Actor{
		Name:  req.Name,
		Sex:   req.Sex,
		Birth: birthday,
	})
	if err != nil {
		res := &api_models.ActorResponse{
			Success: false,
			Error:   err.Error(),
			Actor:   nil,
		}
		err = json.NewEncoder(w).Encode(res)
		//if there's an error with encoding handle it
		if err != nil {
			log.Printf("error sending response %v\n", err)
		}
		//return a bad request and exist the function
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	//everything is good
	//let's return a positive response
	res := &api_models.ActorResponse{
		Success: true,
		Error:   "",
		Actor:   actor,
	}
	err = json.NewEncoder(w).Encode(res)
	if err != nil {
		log.Printf("error encoding after creating comment %v\n", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func updateActor(w http.ResponseWriter, r *http.Request) {
	//get the data from the request
	req := &api_models.UpdateActorRequest{}
	//decode the data
	err := json.NewDecoder(r.Body).Decode(req)
	if err != nil {
		res := &api_models.ActorResponse{
			Success: false,
			Error:   err.Error(),
			Actor:   nil,
		}
		err = json.NewEncoder(w).Encode(res)
		//if there's an error with encoding handle it
		if err != nil {
			log.Printf("error sending response %v\n", err)
		}
		//return a bad request and exist the function
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	pgdb, ok := r.Context().Value("DB").(*pg.DB)
	if !ok {
		res := &api_models.ActorResponse{
			Success: false,
			Error:   "could not get the DB from context",
			Actor:   nil,
		}
		err := json.NewEncoder(w).Encode(res)
		//if there's an error with encoding handle it
		if err != nil {
			log.Printf("error sending response %v\n", err)
		}
		//return a bad request and exist the function
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	//get the commentID to know what comment to modify
	actorID := chi.URLParam(r, "actorID")
	//we get a string but we need to send an int so we convert it
	intActorID, err := strconv.ParseInt(actorID, 10, 64)
	datetime, err := time.Parse("2006-01-02", req.Birth)

	//update the comment
	actor, err := db_models.UpdateActor(pgdb, &db_models.Actor{
		ID:    intActorID,
		Name:  req.Name,
		Sex:   req.Sex,
		Birth: datetime,
	})
	if err != nil {
		res := &api_models.ActorResponse{
			Success: false,
			Error:   err.Error(),
			Actor:   nil,
		}
		err = json.NewEncoder(w).Encode(res)
		//if there's an error with encoding handle it
		if err != nil {
			log.Printf("error sending response %v\n", err)
		}
		//return a bad request and exist the function
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	//return successful response
	res := &api_models.ActorResponse{
		Success: true,
		Error:   "",
		Actor:   actor,
	}
	//send the encoded response to responsewriter
	err = json.NewEncoder(w).Encode(res)
	if err != nil {
		log.Printf("error encoding comments: %v\n", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	//send a 200 response
	w.WriteHeader(http.StatusOK)
}

func deleteActor(w http.ResponseWriter, r *http.Request) {

	//get the db from ctx
	pgdb, ok := r.Context().Value("DB").(*pg.DB)
	if !ok {
		res := &api_models.ActorResponse{
			Success: false,
			Error:   "could not get the DB from context",
			Actor:   nil,
		}
		err := json.NewEncoder(w).Encode(res)
		//if there's an error with encoding handle it
		if err != nil {
			log.Printf("error sending response %v\n", err)
		}
		//return a bad request and exist the function
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	//get the commentID
	actorID := chi.URLParam(r, "actorID")
	intActorID, err := strconv.ParseInt(actorID, 10, 64)
	if err != nil {
		res := &api_models.ActorResponse{
			Success: false,
			Error:   err.Error(),
			Actor:   nil,
		}
		err = json.NewEncoder(w).Encode(res)
		//if there's an error with encoding handle it
		if err != nil {
			log.Printf("error sending response %v\n", err)
		}
		//return a bad request and exist the function
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	//delete comment
	err = db_models.DeleteActor(pgdb, intActorID)
	if err != nil {
		res := &api_models.ActorResponse{
			Success: false,
			Error:   err.Error(),
			Actor:   nil,
		}
		err = json.NewEncoder(w).Encode(res)
		//if there's an error with encoding handle it
		if err != nil {
			log.Printf("error sending response %v\n", err)
		}
		//return a bad request and exist the function
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	//send a 200 response
	w.WriteHeader(http.StatusOK)
}
