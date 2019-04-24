package main

import (
	"./lib"
	"fmt"
	"github.com/urfave/cli"
	"log"
	"os"
)

const appVersion = "0.0.3"

func main() {
	app := lib.InitApp()
	app.Version = appVersion
	app.Action = cliMainAction
	app.Commands = []cli.Command{
		{
			Name:    "init",
			Aliases: []string{"i"},
			Usage:   "create configuration file in current directory",
			Action:  initCommand,
		},
		{
			Name:   "env",
			Usage:  "create env file in current directory",
			Action: envCommand,
		},
	}
	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

// главное действие утилиты
func cliMainAction(c *cli.Context) error {
	workDir := lib.GetWorkDir(c.GlobalString(lib.FlagWorkDir))
	// ищем и грузим конфигурационный файл
	configPath := lib.FindConfFile(workDir)
	config, _ := lib.LoadConfig(configPath)
	// загрузка переменных окружения с файла
	lib.LoadFromEnvFile(c.String(lib.FlagEnvFile), workDir)
	// удаляем предыдущие результаты если файл существует
	outputDir := lib.GetConfigOutputDir(c.GlobalString(lib.FlagOutputDir), workDir)
	lib.ClearPrevResults(outputDir)
	// конвертируем и проверяем переменные на ограничения
	if result, err := lib.CheckVariables(config); err != nil {
		log.Fatal(err)
	} else {
		// сохраняем результат
		lib.SaveResultsToFile(result, outputDir)
	}
	return nil
}

func initCommand(c *cli.Context) error {
	workDir := lib.GetWorkDir(c.GlobalString(lib.FlagWorkDir))
	force := c.GlobalBool(lib.FlagForce)
	lib.GenerateConfigFile(workDir, force)
	return nil
}

func envCommand(c *cli.Context) error {
	workDir := lib.GetWorkDir(c.GlobalString(lib.FlagWorkDir))
	configPath := lib.FindConfFile(workDir)
	config, _ := lib.LoadConfig(configPath)

	if values, err := lib.GenerateEnvFile(config); err != nil {
		log.Fatal(err)
	} else {
		fmt.Println(values)
	}

	return nil
}
