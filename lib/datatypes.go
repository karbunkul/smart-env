package lib

import (
	"errors"
	"github.com/xeipuuv/gojsonschema"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

// проверка значений
func CheckVariables(config Config) (Result, error) {
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
		if variable.Constraints != nil {
			if status, err := ValidateConstraints(variable.Constraints, convertedValue); status == true {
				vars[name] = convertedValue
			} else {
				log.Fatal(err)
			}
		} else {
			vars[name] = convertedValue
		}
	}
	result := Result{
		Version:     ApiVersion,
		Variables:   vars,
		LastUpdated: time.Now().Unix(),
	}
	return result, nil
}

// проверка ограничений
func ValidateConstraints(constraints map[string]interface{}, value interface{}) (bool, error) {
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
