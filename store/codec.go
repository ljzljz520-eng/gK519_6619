package store

import (
	"encoding/json"
	"fmt"
)

func encode(value any) ([]byte, error) {
	data, err := json.Marshal(value)
	if err != nil {
		return nil, fmt.Errorf("encode record: %w", err)
	}
	return data, nil
}

func decode(data []byte, target any) error {
	if len(data) == 0 {
		return fmt.Errorf("record is empty")
	}
	if err := json.Unmarshal(data, target); err != nil {
		return fmt.Errorf("decode record: %w", err)
	}
	return nil
}

func cloneBytes(data []byte) []byte {
	if data == nil {
		return nil
	}
	result := make([]byte, len(data))
	copy(result, data)
	return result
}
