package fs

type FS interface {
	CreateStorage() error
	CleanUp() error
	GetRootDir() string
}
