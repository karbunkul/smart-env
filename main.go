package main

import (
	"fmt"
	"github.com/urfave/cli"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"
)

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
	app.EnableBashCompletion = true
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "env, e",
			Value: "dev",
			Usage: "загрузка переменных из окружения",
		},
		cli.StringFlag{
			Name:  "config-dir, d",
			Value: "",
			Usage: "путь к директории с файлом конфигурации по умолчанию текущая директория",
		},
		cli.StringFlag{
			Name:  "config-name, f",
			Value: "smart-env",
			Usage: "имя файла конфигурации без расширения файла",
		},
	}
	// главное действие
	app.Action = cliMainAction
	return app
}

// ищем путь к файлу конфигурации
func findConfFile(dir string, name string) (string, error) {
	formats := [3]string{"yaml", "yml", "json"}

	for _, format := range formats {
		file := path.Join(dir, name+"."+format)
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
		fi, _ := os.Stat(absPath)
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
func getConfigName(param string) (string, error) {
	return param, nil
}

func cliMainAction(c *cli.Context) error {
	configDir, _ := getConfigDir(c.String("config-dir"))
	configName, _ := getConfigName(c.String("config-name"))

	configPath, _ := findConfFile(configDir, configName)

	data, _ := loadConfig(configPath)
	fmt.Println(data.Variables)

	return nil
}

type Config struct {
	Version   string
	Variables []struct {
		Name   string
		Schema interface{}
	}
	Stages interface{}
}

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
