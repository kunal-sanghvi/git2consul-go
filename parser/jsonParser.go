package parser

import (
	"encoding/json"
	"fmt"
	"git2consul-go/backend"
	"golang.org/x/net/context"
	"io/ioutil"
	"log"
	"strings"
)

type jsonParser struct {
	rootDirectory string
	configFile string
	consulHost string
}

var (
	configBackend backend.ConfigBackend
)

func (jsonConf *jsonParser) ParseFromFile() (*Config, error) {
	conf := new(Config)
	data, err := ioutil.ReadFile(jsonConf.configFile)
	if err != nil {
		return conf, err
	}
	err = json.Unmarshal(data, conf)
	if err != nil {
		return conf, err
	}
	jsonConf.consulHost = conf.ConsulHost
	configBackend = backend.NewConsulBackend(conf.ConsulHost)
	return conf, nil
}

func (jsonConf *jsonParser) ParseConfigFile(ctx context.Context, branch, repo, filePath string) error {
	configMap := make(map[string]interface{})
	data, err := ioutil.ReadFile(jsonConf.rootDirectory + repo + "/" + filePath)
	if err != nil {
		return err
	}
	err = json.Unmarshal(data, &configMap)
	if err != nil {
		return err
	}
	jsonConf.ParseConfigMap(ctx, branch, repo, strings.Split(filePath, ".")[0], configMap)
	return nil
}

func (jsonConf *jsonParser) ParseConfigMap(ctx context.Context, branch, repo, filePath string, configMap map[string]interface{}) {
	for key, value := range configMap {
		switch value.(type) {
		case string, bool, float64, []interface{}:
			consulKey := fmt.Sprintf("%s/%s/%s/%s", repo, branch, filePath, key)
			consulValue := fmt.Sprintf("%v", value)
			if err := configBackend.Populate(consulKey, consulValue); err != nil {
				log.Printf("Couldn't populate key %s with value %s. ERROR: %v", consulKey, consulValue, err)
			}
			//log.Printf("%s %s", consulKey, consulValue)
		default:
			jsonConf.ParseConfigMap(ctx, branch, repo,filePath + "/" + key, determineType(branch, filePath + key, value))
		}
	}
}

func determineType(branch, key string, i interface{}) map[string]interface{} {
	defer func() map[string]interface{} {
		x := recover()
		if x != nil {
			log.Printf("ERROR: %s %s %v", branch, key, i)
			log.Printf("ERROR: Failed to convert %v", x)
		}
		return map[string]interface{}{
			key: fmt.Sprintf("%v", i),
		}
	}()
	return i.(map[string]interface{})
}

func NewJSONParser(rootDirectory, file string) Parser {
	return &jsonParser{
		rootDirectory: rootDirectory,
		configFile: file,
	}
}