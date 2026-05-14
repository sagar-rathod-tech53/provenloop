package utils

import (
	"fmt"
	"math/rand"
	"time"
)

func GenerateOTP(length int) string {

	rand.Seed(time.Now().UnixNano())

	otp := ""

	for i := 0; i < length; i++ {
		otp += fmt.Sprintf("%d", rand.Intn(10))
	}

	return otp
}
