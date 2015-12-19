package util

import (
	"fmt"
	"math/rand"
	"time"

	"code.google.com/p/go-uuid/uuid"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func RandIntSuffix(str, sep string) string {
	return fmt.Sprintf("%s%s%d", str, rand.Int())
}

func RandString() string {
	return fmt.Sprintf("%s-%s", uuid.New(), rand.Int())
}

func RandBytes() []byte {
	return []byte(RandString())
}
