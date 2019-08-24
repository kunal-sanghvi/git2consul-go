package parser

import "context"

type Parser interface {
	ParseFromFile() (*Config, error)
	ParseConfigFile(ctx context.Context, branch, repo, filePath string) error
	ParseConfigMap(ctx context.Context, branch, repo, filePath string, configMap map[string]interface{})
}

type Repo struct {
	Name     string   `json:"name"`
	URL      string   `json:"url"`
	Branches []string `json:"branches"`
}

type Config struct {
	ConsulHost string `json:"host"`
	Repos []*Repo `json:"repos"`
}
