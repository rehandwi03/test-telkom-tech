package util

import (
	"math/rand"
	"strconv"
	"time"
)

func RandString() string {
	randomNumber := rand.New(rand.NewSource(time.Now().UnixNano()))
	return strconv.Itoa(randomNumber.Int())
}
