package db

import (
	"fmt"
	"time"

	"github.com/go-pg/pg/v10"
)

type Actor struct {
	ID    int64     `json:"id"`
	Name  string    `json:"name"`
	Sex   string    `json:"sex" validate:"oneof=male female"`
	Birth time.Time `json:"birthday" validate:"datetime"`
	Films []Film    `json:"films" pg:"many2many:film_to_actors"`
}

func GetActors(db *pg.DB) ([]*Actor, error) {
	actors := make([]*Actor, 0)

	err := db.Model(&actors).Relation("Films").
		Select()

	return actors, err
}

func CreateActor(db *pg.DB, req *Actor) (*Actor, error) {
	_, err := db.Model(req).Insert()
	if err != nil {
		return nil, err
	}

	actor := &Actor{}

	err = db.Model(actor).
		Relation("Films").
		Where("actor.id = ?", req.ID).
		Select()

	return actor, err
}

func UpdateActor(db *pg.DB, req *Actor) (*Actor, error) {
	fmt.Println(req)
	res, err := db.Model(req).
		Where("actor.id = ?", req.ID).
		UpdateNotZero()
	fmt.Println(res, err)
	if err != nil {
		return nil, err
	}

	actor := &Actor{}

	err = db.Model(actor).
		Relation("Films").
		Where("actor.id = ?", req.ID).
		Select()

	return actor, err
}

func DeleteActor(db *pg.DB, actorID int64) error {
	actor := &Actor{}
	film2actor := &FilmToActor{}

	err := db.Model(actor).
		Relation("Films").
		Where("actor.id = ?", actorID).
		Select()
	if err != nil {
		return err
	}

	_, err = db.Model(actor).WherePK().Delete()
	_, err = db.Model(film2actor).Where("actor_id = ?", actor.ID).Delete()

	return err
}
