package model

import (
	"dot/core/v1/entity"
	"gorm.io/gorm"
)

type Creator struct {
	gorm.Model
	Id   string `gorm:"id,omitempty"`
	Name string `gorm:"name"`

	// one-to-many
	Books []Book
}

func (doc *Creator) Bind(creator *entity.Creator) *Creator {
	return &Creator{
		Id:   creator.Id,
		Name: creator.Name,
	}
}

func (doc *Creator) Entity() *entity.Creator {
	return &entity.Creator{
		Id:   doc.Id,
		Name: doc.Name,
	}
}

type Creators []Creator

func (authrs Creators) Bind(creators []entity.Creator) Creators {
	for i := range creators {
		authr := new(Creator).Bind(&creators[i])
		if creators[i].Id != "" {
			authr.Id = creators[i].Id
		}
		authrs = append(authrs, *authr)
	}
	return authrs
}

func (authrs Creators) Generics() []any {
	var data []any
	for i := range authrs {
		data = append(data, authrs[i])
	}
	return data
}

func (authrs Creators) Entities() []entity.Creator {
	var es []entity.Creator
	for i := range authrs {
		es = append(es, *authrs[i].Entity())
	}
	return es
}
