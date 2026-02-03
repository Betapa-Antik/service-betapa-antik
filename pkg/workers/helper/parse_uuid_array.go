package helper

import (
	"encoding/json"
	"fmt"

	"github.com/google/uuid"
)

func ParseUUIDArray(raw string) ([]uuid.UUID, error) {
	if raw == "" {
		return nil, nil
	}

	var ids []string
	if err := json.Unmarshal([]byte(raw), &ids); err != nil {
		return nil, fmt.Errorf("format id harus array JSON")
	}

	var result []uuid.UUID
	for _, s := range ids {
		id, err := uuid.Parse(s)
		if err != nil {
			return nil, fmt.Errorf("uuid tidak valid: %s", s)
		}
		result = append(result, id)
	}

	return result, nil
}
