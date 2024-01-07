package tool

import (
	"fmt"
	"strconv"
	"encoding/base64"
	"encoding/json"
)

func ToString(i interface{}) string {
	return fmt.Sprintf("%v", i)
}

func ToInt8(str string) int8 {
	i,err := strconv.Atoi(str)
	if err != nil {
		return 0
	}
	return int8(i)
}

func ToInt32(str string) int32 {
	it, err := strconv.Atoi(str)
	if err != nil {
		return 0
	}
	return int32(it)
}

func ToInt(str string) int {
	it, err := strconv.Atoi(str)
	if err != nil {
		return 0
	}
	return int(it)
}

func ToInt64(str string) int64 {
	it, err := strconv.Atoi(str)
	if err != nil {
		return 0
	}
	return int64(it)
}

func ToFloat64(str string) float64 {
	it,err := strconv.ParseFloat(str, 64)
	if err != nil {
		return 0
	}
	return float64(it)
}

func Base64Encode(str string) string {
	return base64.StdEncoding.EncodeToString([]byte(str))
}

func Base64Decode(str string) (string, error) {
	b, err := base64.StdEncoding.DecodeString(str)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

func ToJsonString(i interface{}) (string, error) {
	b, e := json.Marshal(i)
	return string(b), e
}