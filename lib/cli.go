package lib

import (
	"github.com/urfave/cli"
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"
)

const FlagDir = "dir"
const FlagOutputDir = "output"
const FlagEnvFile = "env"

// инициализация утилиты
func InitApp() *cli.App {
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
			Name:  FlagEnvFile + ", e",
			Value: "",
			Usage: "file path for env file, default find .env from " + FlagDir,
		},
	}
	return app
}

// разбор флага config-dir
func GetConfigDir(param string) (string, error) {
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
func GetStageName(param string) (string, error) {
	return param, nil
}

func GetConfigOutputDir(param string, configDir string) string {
	if strings.Trim(param, "") == "" {
		param = configDir
	}
	return param
}
