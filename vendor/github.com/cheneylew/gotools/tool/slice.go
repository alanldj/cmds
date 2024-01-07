package tool

import (
	"strings"
	"reflect"
	"github.com/JohnCGriffin/yogofn"
)

//数组字符串包含，即为true
func ArrContain(arr []string, search string) bool {
	for _, item := range arr {
		if strings.Contains(item, search) {
			return true
		}
	}

	return false
}

//数组字符串完全相对，即为true
func ArrIn(arr []string, search string) bool {
	for _, item := range arr {
		if item == search {
			return true
		}
	}

	return false
}

func ArrToInterfaces(arr interface{}) ([]interface{},error) {
	v := reflect.ValueOf(arr)
	if v.Kind() != reflect.Slice {
		return nil, ErrorMSG("arr必须为slice类型")
	}
	var res []interface{}
	switch reflect.TypeOf(arr).Kind() {
	case reflect.Slice:
		s := reflect.ValueOf(arr)
		for i := 0; i < s.Len(); i++ {
			res = append(res, s.Index(i).Interface())
		}
	}
	return res, nil
}

func ArrFilter(function interface{}, collection interface{}) interface{} {
	return yogofn.Filter(function, collection)
}

func ArrMap(function interface{}, collection ...interface{}) interface{} {
	return yogofn.Map(function, collection...)
}

func ArrAny(function interface{}, collection ...interface{}) bool {
	return yogofn.Any(function, collection...)
}

func ArrEvery(function interface{}, collection ...interface{}) bool {
	return yogofn.Every(function, collection...)
}

func ArrReduce(collection interface{} , binary interface{}, init ...interface{}) interface{} {
	return yogofn.Reduce(binary, collection, init...)
}

func MergeStrings(a ...[]string) []string {
	var res []string
	for _, i2 := range a {
		for _, i3 := range i2 {
			res = append(res, i3)
		}
	}
	return res
}

func MergeObjects(a ...[]interface{}) []interface{} {
	var res []interface{}
	for _, i2 := range a {
		for _, i3 := range i2 {
			res = append(res, i3)
		}
	}
	return res
}