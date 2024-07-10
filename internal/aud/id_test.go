// Copyright 2024 Igor Zibarev
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
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestID_String(t *testing.T) {
	f := func(id ID, want string) {
		t.Helper()

		assert.Equal(t, want, id.String())
	}

	f(MustParseID("6ba7b810-9dad-11d1-80b4-00c04fd430c8"), "6ba7b810-9dad-11d1-80b4-00c04fd430c8")
	f(ID{}, "")
}

func TestID_Scan(t *testing.T) {
	f := func(src interface{}, want ID) {
		t.Helper()

		var id ID
		assert.NoError(t, id.Scan(src))
		assert.Equal(t, want, id)
	}

	f("6ba7b810-9dad-11d1-80b4-00c04fd430c8", MustParseID("6ba7b810-9dad-11d1-80b4-00c04fd430c8"))
	f(nil, ID{})
}

func TestParseID(t *testing.T) {
	f := func(s string, want ID, checkErr assert.ErrorAssertionFunc) {
		t.Helper()

		got, err := ParseID(s)
		assert.Equal(t, want, got)
		checkErr(t, err)
	}

	f("6ba7b810-9dad-11d1-80b4-00c04fd430c8", MustParseID("6ba7b810-9dad-11d1-80b4-00c04fd430c8"), assert.NoError)
	f("", ID{}, assert.Error)
	f("invalid", ID{}, assert.Error)
}
