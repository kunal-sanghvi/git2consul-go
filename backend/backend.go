package backend

type ConfigBackend interface {
	Populate(key string, value string) error
}
