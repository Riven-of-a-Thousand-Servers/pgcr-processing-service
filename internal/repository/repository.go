package repository

type Repository[T any] interface {
	save(entity T) (T, error)
}
