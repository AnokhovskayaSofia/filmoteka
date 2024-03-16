package db

import (
	"filmoteka/config"
	"log/slog"

	"github.com/go-pg/pg/v10"
	"github.com/go-pg/pg/v10/orm"
)

func init() {
	// Register many to many model so ORM can better recognize m2m relation.
	// This should be done before dependant models are used.
	orm.RegisterTable((*FilmToActor)(nil))
}

func StartDB(cnf *config.Config) (*pg.DB, error) {
	var (
		opts *pg.Options
		err  error
	)

	opts = &pg.Options{
		Addr:     cnf.PostgresDB.Addr,
		User:     cnf.PostgresDB.User,
		Password: cnf.PostgresDB.Password,
		Database: cnf.PostgresDB.Database,
	}

	db := pg.Connect(opts)
	err = createManyToManyTables(db)

	// data_time, _ := time.Parse("2001-02-02", "2001-02-02")
	// values := []interface{}{
	// 	&db_models.Actor{ID: 1, Name: "Actor1", Sex: "female", Birth: data_time},
	// 	&db_models.Actor{ID: 2, Name: "Actor2", Sex: "female", Birth: data_time},
	// 	&db_models.Film{ID: 1, Name: "Film1", Description: "Desk film1", Date: data_time, Rate: 5},
	// 	&FilmToActor{FilmID: 1, ActorID: 1},
	// 	&FilmToActor{FilmID: 1, ActorID: 2},
	// }
	// for _, v := range values {
	// 	_, err := db.Model(v).Insert()
	// 	if err != nil {
	// 		panic(err)
	// 	}
	// }
	slog.Info("Success start PostgrasDB")
	return db, err
}

func createManyToManyTables(db *pg.DB) error {
	models := []interface{}{
		(*Film)(nil),
		(*Actor)(nil),
		(*FilmToActor)(nil),
	}
	for _, model := range models {
		err := db.Model(model).CreateTable(&orm.CreateTableOptions{
			Temp:        false,
			IfNotExists: true,
		})
		if err != nil {
			return err
		}
	}
	return nil
}
