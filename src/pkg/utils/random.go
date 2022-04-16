package utils

import (
	"github.com/google/uuid"
	"strings"
)

func RandomString() string {
	newUUID, _ := uuid.NewUUID()
	return strings.Replace(newUUID.String(), "-", "", -1)
}
