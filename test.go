package main

import (
	"fmt"
	"reflect"
)

type myint = int

// Test jasd
func Test() {
	x := struct {
		Foo string
		Bar int
	}{"123123", 2}

	for idx, v := range IterStructFieldValue(&x) {
		fmt.Println(idx, v)
	}

	fmt.Println(x)
}

// IterStructFieldValue IterStructFieldValue
func IterStructFieldValue(u interface{}) []interface{} {
	val := reflect.ValueOf(u).Elem()
	length := val.NumField()
	v := make([]interface{}, length)

	for i := 0; i < length; i++ {
		valueField := val.Field(i)
		v[i] = valueField.Interface()
	}

	return v
}
