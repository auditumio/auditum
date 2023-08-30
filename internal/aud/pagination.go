// Copyright 2023 Igor Zibarev
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

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
