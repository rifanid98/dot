package books

import (
	"dot/core/v1/entity"
)

type Get struct {
	Id string `json:"id" validate:"required" example:"65c1de91056ae9755c64ffba" param:"id"`
}

type Delete struct {
	Id string `json:"id" validate:"required" example:"65c1de91056ae9755c64ffba" param:"id"`
}

type List struct {
	Page  int64 `json:"page" validate:"omitempty,min=1" example:"1" query:"page"`
	Limit int64 `json:"limit" validate:"omitempty,min=1" example:"1" query:"limit"`
}

type Book struct {
	Id     string `json:"id" validate:"omitempty" swaggerignore:"true" param:"id" query:"id"`
	Name   string `json:"name" validate:"required"`
	Author string `json:"author" validate:"required"`
}

func (a *Book) Entity() *entity.Book {
	return &entity.Book{
		Id:     a.Id,
		Name:   a.Name,
		Author: a.Author,
	}
}

type BookPartial struct {
	Id     string `json:"id" validate:"required" swaggerignore:"true" param:"id" query:"id"`
	Name   string `json:"name" validate:""`
	Author string `json:"author" validate:""`
}

func (a *BookPartial) Entity() *entity.Book {
	return &entity.Book{
		Id:     a.Id,
		Name:   a.Name,
		Author: a.Author,
	}
}
