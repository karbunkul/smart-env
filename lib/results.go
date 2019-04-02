package lib

import (
	"encoding/json"
	"io/ioutil"
	"log"
)

// структура файла с переменными
type Result struct {
	Version     string                 `json:"version"`
	LastUpdated int64                  `json:"lastUpdated"`
	Variables   map[string]interface{} `json:"variables"`
}

// сохранение результатов в файл
func SaveResultsToFile(result Result, filePath string) {
	resultJson, _ := json.Marshal(result)
	if err := ioutil.WriteFile(filePath, resultJson, 0775); err != nil {
		log.Fatal(err)
	}
}
