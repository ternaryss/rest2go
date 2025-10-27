package rest2go

import "fmt"

const (
	DefaultPage     int = 1
	DefaultPageSize int = 20
)

type Pagination struct {
	Page     int
	Pages    int
	Size     int
	Pageable int
	Limit    int
	Offset   int
}

func NewPagination(page, size, all int) (Pagination, error) {
	pages := 0

	if size > 0 && all > 0 {
		pages = all / size

		if all%size != 0 {
			pages = pages + 1
		}
	}

	if pages == 0 || page <= 0 || size <= 0 || page > pages {
		return Pagination{}, fmt.Errorf("invalid pagination [pages = %d, page = %d, size = %d]", pages, page, size)
	}

	return Pagination{
		Page:     page,
		Pages:    pages,
		Size:     size,
		Pageable: all,
		Limit:    size,
		Offset:   (page - 1) * size,
	}, nil
}

type PageDto[T any] struct {
	Page     int `json:"page"`
	Pages    int `json:"pages"`
	Size     int `json:"size"`
	Pageable int `json:"pageable"`
	Content  []T `json:"content"`
}

func NewPageDto[T any](pagination Pagination, content []T) PageDto[T] {
	return PageDto[T]{
		Page:     pagination.Page,
		Pages:    pagination.Pages,
		Size:     pagination.Size,
		Pageable: pagination.Pageable,
		Content:  content,
	}
}

func EmptyPageDto[T any]() PageDto[T] {
	return PageDto[T]{}
}
