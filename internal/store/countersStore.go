package store

type CountersStore interface {
	Store

	IncrCounter(userId string) (uint64, error)
	DecrCounter(userId string) (uint64, error)
	GetTotalMessages(userId string) (uint64, error)
}
