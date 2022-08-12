package model

type Paginated[T any] struct {
	Items      []T
	TotalCount int
}
