package aud

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
)

var enc = base64.RawURLEncoding

type Cursor interface {
	Empty() bool
}

func DecodePageToken(token string, cursor any) error {
	if token == "" {
		return nil
	}

	v, err := enc.DecodeString(token)
	if err != nil {
		return fmt.Errorf("decode token string: %v", err)
	}

	if err := json.Unmarshal(v, cursor); err != nil {
		return fmt.Errorf("unmarshal token from json: %v", err)
	}

	return nil
}

func EncodePageToken(cursor Cursor) (string, error) {
	if cursor.Empty() {
		return "", nil
	}

	v, err := json.Marshal(cursor)
	if err != nil {
		return "", fmt.Errorf("marshal token to json: %v", err)
	}

	return enc.EncodeToString(v), nil
}
