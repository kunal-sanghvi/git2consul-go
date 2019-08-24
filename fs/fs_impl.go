package fs

import "os"

type localFS struct {
	rootDirectory string
}

func (lfs *localFS) CreateStorage() error {
	return os.Mkdir(lfs.rootDirectory, os.FileMode(0755))
}

func (lfs *localFS) CleanUp() error {
	return os.RemoveAll(lfs.rootDirectory)
}

func (lfs *localFS) GetRootDir() string {
	return lfs.rootDirectory
}

func NewLocalFS(rootDir string) FS {
	return &localFS{
		rootDirectory: rootDir,
	}
}
