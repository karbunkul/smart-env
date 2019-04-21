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

type CastedValue struct {
	Value interface{} `json:"value"`
	Type  string      `json:"type"`
}

func castToType(value string, castType string) (*CastedValue, error) {
	switch strings.ToLower(castType) {
	case "number", "int":
		if castedValue, err := convertToInt(value); err != nil {
			return nil, err
		} else {
			return &CastedValue{Value: castedValue, Type: "number"}, nil
		}
	case "float":
		if castedValue, err := convertToFloat(value); err != nil {
			return nil, err
		} else {
			return &CastedValue{Value: castedValue, Type: "float"}, nil
		}
	case "bool", "boolean":
		if castedValue, err := convertToBool(value); err != nil {
			return nil, err
		} else {
			return &CastedValue{Value: castedValue, Type: "boolean"}, nil
		}
	default:

		return &CastedValue{Value: value, Type: "string"}, nil
	}
}

// проверка значений
func CheckVariables(config Config) (*Result, error) {
	vars := map[string]ResultVariable{}
	for name, variable := range config.Variables {
		value := os.Getenv(name)
		if strings.Trim(value, "") == "" {
			log.Fatal(errors.New("env " + name + " is empty"))
		}
		castedValue, err := castToType(value, variable.CastTo)
		if err != nil {
			return nil, err
		}
		if variable.Constraints != nil {
			if status, err := ValidateConstraints(variable.Constraints, castedValue.Value); status == true {
				vars[name] = ResultVariable{
					Type:  castedValue.Type,
					Value: castedValue.Value,
				}
			} else {
				log.Fatal(err)
			}
		} else {
			vars[name] = ResultVariable{
				Type:  castedValue.Type,
				Value: castedValue.Value,
			}
		}
	}
	result := Result{
		Version:     ApiVersion,
		Variables:   vars,
		LastUpdated: time.Now().Unix(),
	}
	return &result, nil
}

// проверка ограничений
func ValidateConstraints(constraints map[string]interface{}, value interface{}) (bool, error) {
	schemaLoader := gojsonschema.NewGoLoader(constraints)
	valueLoader := gojsonschema.NewGoLoader(value)
	result, err := gojsonschema.Validate(schemaLoader, valueLoader)
	if err != nil {
		return false, errors.New(err.Error())
	}

	if result.Valid() {
		return true, nil
	} else {
		return false, errors.New(result.Errors()[0].String())
	}
}

// преобразовать строку в число
func convertToInt(value string) (*int64, error) {
	result, err := strconv.ParseInt(value, 10, 64)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// преобразовать строку в число с плавающей точкой
func convertToFloat(value string) (*float64, error) {
	result, err := strconv.ParseFloat(value, 10)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// преобразовать строку в булево значение
func convertToBool(value string) (*bool, error) {
	result, err := strconv.ParseBool(value)
	if err != nil {
		return nil, err
	}
	return &result, nil
}
