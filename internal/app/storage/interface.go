package storage

type Storage interface {
	Put(string) string
	Get(string) (string, bool)
}
