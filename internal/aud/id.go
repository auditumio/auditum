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
	"database/sql/driver"

	"github.com/gofrs/uuid/v5"
)

// ID is a unique identifier.
type ID uuid.UUID

func (id ID) String() string {
	return uuid.UUID(id).String()
}

func (id ID) IsEmpty() bool {
	return uuid.UUID(id).IsNil()
}

func (id ID) Value() (driver.Value, error) {
	if id.IsEmpty() {
		return nil, nil
	}

	return uuid.UUID(id).Value()
}

func (id *ID) Scan(src interface{}) error {
	return (*uuid.UUID)(id).Scan(src)
}

// NewID returns a new identifier.
func NewID() (ID, error) {
	id, err := uuid.NewV7()
	if err != nil {
		return ID{}, err
	}
	return ID(id), nil
}

// MustNewID returns a new identifier, panicking on error.
func MustNewID() ID {
	return ID(uuid.Must(uuid.NewV7()))
}

func ParseID(s string) (ID, error) {
	id, err := uuid.FromString(s)
	if err != nil {
		return ID{}, err
	}

	return ID(id), nil
}

func MustParseID(s string) ID {
	return ID(uuid.Must(uuid.FromString(s)))
}
