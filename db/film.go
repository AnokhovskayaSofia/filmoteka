package db

import (
	"time"

	"github.com/go-pg/pg/v10"
)

type Film struct {
	ID          int       `json:"id"`
	Name        string    `json:"name" validate:"min=1,max=150"`
	Description string    `json:"description" validate:"max=1000"`
	Date        time.Time `json:"date"`
	Rate        int       `json:"rate" validate:"gte=0,lte=10"`
	Actors      []Actor   `json:"actors" pg:"many2many:film_to_actors"`
}

type FilmToActor struct {
	FilmID  int
	ActorID int
}

func GetFilms(db *pg.DB, sortBy string, filter []string) ([]*Film, error) {
	films := make([]*Film, 0)
	var err error
	if cap(filter) > 0 {
		err = db.Model(&films).Relation("Actors").Where("? like '%' || ? || '%'", pg.Ident(filter[0]), filter[1]).Order(sortBy).
			Select()
		err = db.Model(&films).Relation("Actors").Where("? like '%' || ? || '%'", pg.Ident(filter[0]), filter[1]).Order(sortBy).
			Select()
	} else {
		err = db.Model(&films).Relation("Actors").Order(sortBy).
			Select()
	}

	return films, err
}

func CreateFilm(db *pg.DB, req *Film, req_actors []int) (*Film, error) {
	_, err := db.Model(req).Insert()

	for _, actor_id := range req_actors {
		req := FilmToActor{req.ID, actor_id}
		_, err = db.Model(&req).Insert()
	}

	if err != nil {
		return nil, err
	}

	film := &Film{}
	err = db.Model(film).
		Relation("Actors").
		Where("film.id = ?", req.ID).
		Select()

	return film, err
}

func UpdateFilm(db *pg.DB, req *Film) (*Film, error) {
	_, err := db.Model(req).
		WherePK().
		Update()
	if err != nil {
		return nil, err
	}

	film := &Film{}

	err = db.Model(film).
		Relation("Actors").
		Where("film.id = ?", req.ID).
		Select()

	return film, err
}

func DeleteFilm(db *pg.DB, filmID int64) error {
	film := &Film{}
	film2actor := &FilmToActor{}

	err := db.Model(film).
		Relation("Actors").
		Where("film.id = ?", filmID).
		Select()
	if err != nil {
		return err
	}

	_, err = db.Model(film).WherePK().Delete()
	_, err = db.Model(film2actor).Where("film_id = ?", film.ID).Delete()

	return err
}
