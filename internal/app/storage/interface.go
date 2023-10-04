package storage

type Storage interface {
	Put(string) (string, error)
	Get(string) (string, bool)
}
