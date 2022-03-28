package utils

import (
	"strings"

	"github.com/satori/go.uuid"
)

// GetUUID1 get uuid v1
// NewV1 returns UUID based on current timestamp and MAC address
func GetUUID1() (string, error) {
	id := uuid.NewV1()
	return id.String(), nil
}

// GetUUID4 get uuid v4
// NewV4 returns random generated UUID.
func GetUUID4() (string, error) {
	id := uuid.NewV4()
	return id.String(), nil
}

// GenerateRequestID generate request id
func GenerateRequestID() string {
	id, err := GetUUID1()
	if err != nil {
		id, err = GetUUID4()
	}
	return strings.Replace(id, "-", "", -1)
}
