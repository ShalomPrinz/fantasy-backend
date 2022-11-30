package utils

import (
	"log"

	"github.com/mitchellh/mapstructure"
)

func MapToStruct[T any](_map map[string]any) T {
	var _struct T
	err := mapstructure.Decode(_map, &_struct)
	if err != nil {
		log.Fatalf("Couldn't convert given map to struct. %v", err)
	}
	return _struct
}

func StructToMap[T any](_struct T) map[string]any {
	var _map map[string]any
	err := mapstructure.Decode(_struct, &_map)
	if err != nil {
		log.Fatalf("Couldn't convert given struct to map. %v", err)
	}
	return _map
}

func GetStringFromStruct[T any](_struct T, field string) string {
	_map := StructToMap(_struct)
	str, properCast := _map[field].(string)
	if !properCast {
		log.Fatalf("Couldn't convert field %v value to string in struct %v", field, _struct)
	}
	return str
}
