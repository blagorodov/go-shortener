package repositories

type LinksRepository interface {
	Put(string) string
	Get(string) (string, bool)
}
