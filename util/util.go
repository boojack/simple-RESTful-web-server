package util

import (
	"fmt"
	"net/http"
	"reflect"
	"strconv"
	"time"

	guuid "github.com/google/uuid"
)

// ParseInt parse string to int
func ParseInt(s string) (int, error) {
	num, err := strconv.Atoi(s)

	if err != nil {
		fmt.Println(err.Error())
	}

	return num, err
}

// GenUUID generate google uuid
func GenUUID() string {
	return guuid.New().String()
}

// GetNowTime get now time
func GetNowTime() time.Time {
	return time.Now()
}

// GetKeyValueFromCookie get key value
func GetKeyValueFromCookie(key string, r *http.Request) (string, error) {
	cookie, err := r.Cookie(key)

	if err == nil {
		return cookie.Value, nil
	}
	return "", err
}

// IterStructFieldAddr use for figure out the full struct field "addrs"
func IterStructFieldAddr(u interface{}) []interface{} {
	val := reflect.ValueOf(u).Elem()
	length := val.NumField()
	v := make([]interface{}, length)

	for i := 0; i < length; i++ {
		valueField := val.Field(i)
		v[i] = valueField.Addr().Interface()
	}

	return v
}

// IterStructFieldValue use for figure out the full struct field "values"
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
