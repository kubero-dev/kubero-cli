package utils

import (
	"testing"
)

var (
	utils = NewUtils()
)

func TestGenerateRandomString(t *testing.T) {
	length := 10
	result := utils.GenerateRandomString(length, "")
	if len(result) != length {
		t.Errorf("Expected string length of %d, but got %d", length, len(result))
	}
}

func TestCheckBinary(t *testing.T) {
	binary := "go"
	if !utils.CheckBinary(binary) {
		t.Errorf("Expected to find '%s' binary, but it was not found", binary)
	}
}
