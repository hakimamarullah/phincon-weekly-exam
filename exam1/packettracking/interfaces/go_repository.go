package interfaces

type GoRepository[T any] interface {
	Add(data T)
	GetAll() []T
	FindById(id string) (bool, *T)
}
