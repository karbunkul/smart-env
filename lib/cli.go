package lib

import (
	"fmt"
	"github.com/joho/godotenv"
	"github.com/urfave/cli"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"
)

const FlagWorkDir = "dir"
const FlagOutputDir = "output"
const FlagEnvFile = "env"

// инициализация утилиты
func InitApp() *cli.App {
	app := cli.NewApp()
	app.Name = "smart-env"
	app.Usage = "Утилита проверки переменных окружения"
	app.Author = "Alexander Pokhodyun (karbunkul)"
	app.Email = "karbunkul@yourtask.ru"
	app.Copyright = "(c) Alexander Pokhodyun 2019"
	app.EnableBashCompletion = true
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  FlagWorkDir + ", d",
			Value: "",
			Usage: "path for search config file, default value from option --" + FlagWorkDir,
		},
		cli.StringFlag{
			Name:  FlagOutputDir + ", o",
			Value: "",
			Usage: "output directory, default value from option --" + FlagWorkDir,
		},
		cli.StringFlag{
			Name:  FlagEnvFile + ", e",
			Value: "",
			Usage: "file path for env file, default find .env from option --" + FlagWorkDir,
		},
	}
	return app
}

// разбор флага dir
func GetWorkDir(param string) string {
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
		return param
	} else {
		log.Fatal(err)
	}
	return ""
}

func GetConfigOutputDir(param string, workDir string) string {
	if strings.Trim(param, "") == "" {
		param = workDir
	}
	return param
}

func LoadFromEnvFile(filepath string, workDir string) {
	if filepath == "" {
		envFilePath := path.Join(workDir, ".env")
		if _, err := os.Stat(envFilePath); err == nil {
			log.Println("load from enf file: " + envFilePath)
			if err := godotenv.Load(envFilePath); err != nil {
				log.Fatal(err)
			}
		}
	} else {
		if path.IsAbs(filepath) != true {
			// если путь к файлу относительный то изменяем его на абсолютный на основе значения workDir
			filepath = path.Join(workDir, path.Base(filepath))
		}
		if err := godotenv.Load(filepath); err != nil {
			log.Fatal(err)
		}
	}
}

// генерируем конфигурационный файл
func GenerateConfigFile(workDir string) {
	fn := path.Join(workDir, "smart-env.yaml")
	if _, err := os.Stat(fn); os.IsNotExist(err) {
		variables := map[string]ConfigVariable{
			"FOO": {
				CastTo: "string",
			},
		}
		config := Config{
			Version:   ApiVersion,
			Variables: variables,
		}
		// преобразуем структуру в yaml
		data, err := yaml.Marshal(config)
		if err != nil {
			log.Fatalf("error: %v", err)
		}
		if err := ioutil.WriteFile(path.Join(fn), nil, 0775); err != nil {
			log.Fatal(err)
		}
		file, err := os.OpenFile(fn, os.O_WRONLY|os.O_APPEND, 0775)
		if err != nil {
			log.Fatalf("failed opening file: %s", err)
		}
		defer file.Close()

		_, err = file.WriteString("# This file was generated by smart-env see more by link https://github.com/karbunkul/smart-env\n\n")
		_, err = file.Write(data)
		if err != nil {
			log.Fatalf("failed writing to file: %s", err)
		}
	} else {
		fmt.Println("config file already exist")
	}
}
