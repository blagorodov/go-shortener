package storage

type Repository interface {
	Put(string) string
	Get(string) (string, bool)
}
