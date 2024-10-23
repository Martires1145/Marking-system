package util

import (
	"fmt"
	"math/rand"
	"time"
)

func RandToken() string {
	rand.Seed(time.Now().Unix())
	return fmt.Sprintf("%.4d", rand.Uint32()%100000)
}
