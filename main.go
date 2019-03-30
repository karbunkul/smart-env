package main

import (
	"github.com/urfave/cli"
	"log"
	"os"
	"path/filepath"
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

	return "", nil
}

// разбор флага config-dir
func getConfigDir(param string) (string, error) {
	if param == "" {
		// если директория не задана то по умолчанию текущая директория
		param, _ = os.Getwd()
	}
	// полный путь из относительного
	if path, _ := filepath.Abs(param); filepath.IsAbs(filepath.Dir(path)) == true {
		param = filepath.Dir(path)
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

	println(configDir, configName)

	return nil
}
