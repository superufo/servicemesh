package common

import (
	"fmt"
	"reflect"
	"unsafe"
)

// StructToMap  结构体转map
func StructToMap(object interface{}) map[string]interface{} {
	t := reflect.TypeOf(object)
	v := reflect.ValueOf(object)
	res := make(map[string]interface{})

	if t != nil {
		for i := 0; i < t.NumField(); i++ {
			res[string(t.Field(i).Name)] = v.Field(i).Interface()
		}
	}

	return res
}

// StructToMapNoEmpty  todo 优化这个函数的写法  结构体转map
func StructToMapNoEmpty(object interface{}) map[string]interface{} {
	t := reflect.TypeOf(object)
	v := reflect.ValueOf(object)
	res := make(map[string]interface{})

	fmt.Printf("StructToMap %+v %+v", t, v)
	for i := 0; i < t.NumField(); i++ {
		val := fmt.Sprintf("%v", v.Field(i))
		if val != "" && val != "0" {
			res[string(t.Field(i).Name)] = v.Field(i).Interface()
		}
	}
	return res
}

func StrToBytes(s string) []byte {
	x := (*[2]uintptr)(unsafe.Pointer(&s))
	h := [3]uintptr{x[0], x[1], x[1]}
	return *(*[]byte)(unsafe.Pointer(&h))
}

func BytesToStr(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}
