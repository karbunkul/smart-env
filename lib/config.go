package lib

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"
)

// структура конфигурационного файла
type Config struct {
	Version   string `yaml:"version"`
	Variables map[string]struct {
		ValueType   string                 `yaml:"valueType"`
		Constraints map[string]interface{} `yaml:"constraints"`
	} `yaml:"variables"`
	Stages map[string]map[string]string `yaml:"stages"`
}

// ищем путь к файлу конфигурации
func FindConfFile(dir string) (string, error) {
	formats := [3]string{"yaml", "yml", "json"}
	for _, format := range formats {
		file := path.Join(dir, "smart-env."+format)
		if _, err := os.Stat(file); !os.IsNotExist(err) {
			return file, nil
		}
	}
	os.Exit(-1)
	return "", nil
}

// загрузка конфигурационного файла
func LoadConfig(path string) (Config, error) {
	data, _ := ioutil.ReadFile(path)
	var config Config
	ext := filepath.Ext(path)[1:]

	switch strings.ToLower(ext) {
	case "yaml", "yml":
		err := yaml.Unmarshal(data, &config)
		if err != nil {
			log.Fatalf("error: %v", err)
		}
		break
	}
	return config, nil
}