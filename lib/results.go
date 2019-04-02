package lib

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"path"
)

// структура файла с переменными
type Result struct {
	Version     string                 `json:"version"`
	LastUpdated int64                  `json:"lastUpdated"`
	Variables   map[string]interface{} `json:"variables"`
}

func getResultsPath(workDir string) string {
	return path.Join(workDir, "smart-env.vars.json")
}

// сохранение результатов в файл
func SaveResultsToFile(result Result, workDir string) {
	resultJson, _ := json.Marshal(result)
	if err := ioutil.WriteFile(getResultsPath(workDir), resultJson, 0775); err != nil {
		log.Fatal(err)
	}
}

func ClearPrevResults(workDir string) {
	filePath := getResultsPath(workDir)
	if _, err := os.Stat(filePath); !os.IsNotExist(err) {
		if removeErr := os.Remove(filePath); removeErr != nil {
			log.Fatal(removeErr)
		}
	}
}
