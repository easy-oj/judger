package tools

import (
	"fmt"

	"github.com/satori/go.uuid"
)

func GenUUID() string {
	return fmt.Sprintf("%x", uuid.Must(uuid.NewV4()).Bytes())
}
