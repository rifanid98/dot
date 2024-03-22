package helper

import (
	"dot/core"
	"dot/pkg/util"
	"encoding/json"
	"reflect"
	"strconv"
)

var log = util.NewLogger()

func DataToBytes(data any) []byte {
	bts, _ := json.Marshal(data)
	return bts
}

func StringToBytes(data string) []byte {
	return []byte(data)
}

func DataToString(data any) string {
	if data == nil {
		return ""
	}
	if data == "" {
		return ""
	}
	if reflect.TypeOf(data).String() == "string" {
		return data.(string)
	}
	return string(DataToBytes(data))
}

func StringToStruct(data string, dest any) *core.CustomError {
	err := json.Unmarshal(StringToBytes(data), dest)
	if err != nil {
		return &core.CustomError{
			Code:    core.INTERNAL_SERVER_ERROR,
			Message: err.Error(),
		}
	}

	return nil
}

func BytesToData(bts []byte, dest any) *core.CustomError {
	err := json.Unmarshal(bts, dest)
	if err != nil {
		return &core.CustomError{
			Code:    core.INTERNAL_SERVER_ERROR,
			Message: err.Error(),
		}
	}

	return nil
}

func DataToInt(data any) int64 {
	if data == nil {
		return 0
	}

	str := DataToString(data)
	if str == "" {
		return 0
	}

	i, err := strconv.ParseInt(str, 10, 64)
	if err != nil {
		return 0
	}

	return i
}

func GetTagNamebyName(name string, valueStruct any) string {
	field, ok := reflect.TypeOf(valueStruct).Elem().FieldByName(name)
	if !ok {
		return "default_name_by_system"
	}

	return string(field.Tag.Get("json"))
}
