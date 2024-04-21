package helper

import "github.com/google/uuid"

func GetToken() string {
	return uuid.New().String()
}
