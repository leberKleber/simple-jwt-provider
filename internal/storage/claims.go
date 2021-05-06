package storage

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
)

type Claims map[string]interface{}

// Scan scan value into Claims
func (j *Claims) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return fmt.Errorf("failed to unmarshal Claims value: %s", value)
	}

	err := json.Unmarshal(bytes, &j)
	return err
}

// Value return json value as byte slice
func (j Claims) Value() (driver.Value, error) {
	return json.Marshal(j)
}
