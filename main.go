package main

import (
	"./lib"
	"github.com/urfave/cli"
	"log"
	"os"
)

func main() {
	app := lib.InitApp()
	app.Version = version()
	app.Action = cliMainAction
	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

// главное действие утилиты
func cliMainAction(c *cli.Context) error {
	// параметры утилиты
	workDir := lib.GetWorkDir(c.String(lib.FlagWorkDir))
	outputDir := lib.GetConfigOutputDir(c.String(lib.FlagOutputDir), workDir)
	// ищем и грузим конфигурационный файл
	configPath := lib.FindConfFile(workDir)
	config, _ := lib.LoadConfig(configPath)
	// загрузка переменных окружения с файла
	lib.LoadFromEnvFile(c.String(lib.FlagEnvFile), workDir)
	// удаляем предыдущие результаты если файл существует
	lib.ClearPrevResults(outputDir)
	// конвертируем и проверяем переменные на ограничения
	result, _ := lib.CheckVariables(config)
	// сохраняем результат
	lib.SaveResultsToFile(result, outputDir)
	return nil
}

func version() string {
	return "0.0.1"
}
