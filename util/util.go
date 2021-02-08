package util

import (
	"fmt"
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
