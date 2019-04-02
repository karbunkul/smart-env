package main

import (
	"./lib"
	"github.com/urfave/cli"
	"log"
	"os"
	"path"
)

func main() {
	app := lib.InitApp()
	app.Action = cliMainAction
	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

// главное действие утилиты
func cliMainAction(c *cli.Context) error {
	// параметры утилиты
	configDir, _ := lib.GetConfigDir(c.String(lib.FlagDir))
	outputDir := lib.GetConfigOutputDir(c.String(lib.FlagOutputDir), configDir)
	// ищем и грузим конфигурационный файл
	configPath, _ := lib.FindConfFile(configDir)
	config, _ := lib.LoadConfig(configPath)
	// конвертируем и проверяем переменные на ограничения
	result, _ := lib.CheckVariables(config)
	// сохраняем результат
	outputFile := path.Join(outputDir, "smart-env.vars.json")
	lib.SaveResultsToFile(result, outputFile)
	return nil
}
