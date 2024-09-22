package util

import (
	nanoid "github.com/matoous/go-nanoid/v2"

)

func GenerateId() string {
	return nanoid.MustGenerate("0123456789abcdefghijklmnopqrstuvwxyz", 21)
}
