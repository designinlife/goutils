package goutils

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestBase64Encoding(t *testing.T) {
	assert.Equal(t, "Z2l0aHViLmNvbQ==", Base64Encoding("github.com"))
}

func TestBase64Decoding(t *testing.T) {
	assert.Equal(t, "github.com", Base64Decoding("Z2l0aHViLmNvbQ=="))
}

func TestBase58Encoding(t *testing.T) {
	assert.Equal(t, "6oyQPEHvqTRYo6", Base58Encoding("github.com"))
}

func TestBase58Decoding(t *testing.T) {
	assert.Equal(t, "github.com", Base58Decoding("6oyQPEHvqTRYo6"))
}
