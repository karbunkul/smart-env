package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/urfave/cli"
	"github.com/xeipuuv/gojsonschema"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

const FlagDir = "dir"
const FlagOutputDir = "output"
const FlagStage = "stage"

func main() {
	app := initApp()

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

// инициализация утилиты
func initApp() *cli.App {
	app := cli.NewApp()
	app.Name = "smart-env"
	app.Usage = "Утилита проверки переменных окружения"
	app.Author = "Alexander Pokhodyun (karbunkul)"
	app.Email = "karbunkul@yourtask.ru"
	app.Version = "0.0.1"
	app.Copyright = "(c) Alexander Pokhodyun 2019"
	app.EnableBashCompletion = true
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  FlagDir + ", d",
			Value: "",
			Usage: "path for search config file, default value current work directory",
		},
		cli.StringFlag{
			Name:  FlagOutputDir + ", o",
			Value: "",
			Usage: "output directory, default value from " + FlagDir,
		},
		cli.StringFlag{
			Name:  FlagStage + ", s",
			Value: "",
			Usage: "current stage for env values, don't use in production env",
		},
	}
	// главное действие
	app.Action = cliMainAction
	return app
}

// ищем путь к файлу конфигурации
func findConfFile(dir string) (string, error) {
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

// разбор флага config-dir
func getConfigDir(param string) (string, error) {
	if param == "" {
		// если директория не задана то по умолчанию текущая директория
		param, _ = os.Getwd()
	}
	// полный путь из относительного
	if absPath, _ := filepath.Abs(param); filepath.IsAbs(filepath.Dir(absPath)) == true {
		fi, err := os.Stat(absPath)
		if err != nil {
			log.Fatal(err)
		}
		if fi.Mode().IsDir() {
			param = absPath
		} else {
			param = path.Dir(absPath)
		}
	}
	// проверяем существует ли директория
	if _, err := os.Stat(param); !os.IsNotExist(err) {
		return param, nil
	} else {
		log.Println(err)
		return "", err
	}
}

// разбор флага config-name
func getStageName(param string) (string, error) {
	return param, nil
}

// главное действие утилиты
func cliMainAction(c *cli.Context) error {
	configDir, _ := getConfigDir(c.String(FlagDir))
	outputDir := getConfigOutputDir(c.String(FlagOutputDir), configDir)
	configPath, _ := findConfFile(configDir)
	stageName, _ := getStageName(c.String(FlagStage))
	fmt.Println(stageName)
	config, _ := loadConfig(configPath)
	result, _ := checkVariables(config)
	resultJson, _ := json.Marshal(result)
	outputFile := path.Join(outputDir, "smart-env.vars.json")
	if err := ioutil.WriteFile(outputFile, resultJson, 0775); err != nil {
		log.Fatal(err)
	}
	return nil
}

func getConfigOutputDir(param string, configDir string) string {
	if strings.Trim(param, "") == "" {
		param = configDir
	}
	return param
}

// преобразовать строку в число
func convertToInt(value string) int64 {
	result, err := strconv.ParseInt(value, 10, 64)
	if err != nil {
		log.Fatal(err)
	}
	return result
}

// преобразовать строку в число с плавающей точкой
func convertToFloat(value string) float64 {
	result, err := strconv.ParseFloat(value, 10)
	if err != nil {
		log.Fatal(err)
	}
	return result
}

// преобразовать строку в булево значение
func convertToBool(value string) bool {
	result, err := strconv.ParseBool(value)
	if err != nil {
		log.Fatal(err)
	}
	return result
}

// проверка значений
func checkVariables(config Config) (Result, error) {
	vars := map[string]interface{}{}
	for name, variable := range config.Variables {
		value := os.Getenv(name)
		if strings.Trim(value, "") == "" {
			log.Fatal(errors.New("env " + name + " is empty"))
		}
		var convertedValue interface{}
		switch strings.ToLower(variable.ValueType) {
		case "number", "int":
			convertedValue = convertToInt(value)
			break
		case "float":
			convertedValue = convertToFloat(value)
			break
		case "bool", "boolean":
			convertedValue = convertToBool(value)
			break
		default:
			convertedValue = value
			break
		}
		if status, err := validateConstraints(variable.Constraints, convertedValue); status == true {
			vars[name] = convertedValue
		} else {
			log.Fatal(err)
		}
	}
	result := Result{
		Version:     "1.0",
		Variables:   vars,
		LastUpdated: time.Now().Unix(),
	}
	return result, nil
}

// проверка ограничений
func validateConstraints(constraints map[string]interface{}, value interface{}) (bool, error) {
	schemaLoader := gojsonschema.NewGoLoader(constraints)
	valueLoader := gojsonschema.NewGoLoader(value)
	result, err := gojsonschema.Validate(schemaLoader, valueLoader)
	if result != nil {
		if result.Valid() {
			return true, nil
		} else {
			log.Fatal(result.Errors())
			return false, nil
		}
	} else {
		return false, err
	}
}

// структура файла с переменными
type Result struct {
	Version     string                 `json:"version"`
	LastUpdated int64                  `json:"lastUpdated"`
	Variables   map[string]interface{} `json:"variables"`
}

// структура конфигурационного файла
type Config struct {
	Version   string `yaml:"version"`
	Variables map[string]struct {
		ValueType   string                 `yaml:"valueType"`
		Constraints map[string]interface{} `yaml:"constraints"`
	} `yaml:"variables"`
	Stages map[string]map[string]string `yaml:"stages"`
}

// загрузка конфигурационного файла
func loadConfig(path string) (Config, error) {
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
