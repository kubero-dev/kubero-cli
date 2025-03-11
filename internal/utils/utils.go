package utils

import (
	"math/rand"
	"os/exec"
	"time"
)

func GenerateRandomString(length int, chars string) string {
	if chars == "" {
		chars = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789!+?._-%"
	}
	var letterRunes = []rune(chars)

	b := make([]rune, length)
	rand.New(rand.NewSource(time.Now().UnixNano()))

	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

func CheckBinary(binary string) bool {
	_, err := exec.LookPath(binary)
	return err == nil
}

func CheckAllBinaries(binaries ...string) {
	for _, binary := range binaries {
		if !CheckBinary(binary) {
			panic(binary + " binary is not installed")
		}
	}
}
