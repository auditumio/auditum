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

package auditumv1alpha1

import (
	"fmt"
	"unicode/utf8"
)

func validateProjectDisplayName(src string) error {
	const (
		minChars = 3
		maxChars = 64
	)

	if src == "" {
		return fmt.Errorf(`must not be empty`)
	}

	chars := utf8.RuneCountInString(src)
	if chars < minChars {
		return fmt.Errorf(`must not be shorter than %d characters`, minChars)
	}
	if chars > maxChars {
		return fmt.Errorf(`must not be longer than %d characters`, maxChars)
	}

	return nil
}

func validateProjectExternalID(src string) error {
	const (
		minChars = 3
		maxChars = 64
	)

	if src == "" {
		return nil
	}

	chars := utf8.RuneCountInString(src)
	if chars < minChars {
		return fmt.Errorf(`must not be shorter than %d characters`, minChars)
	}
	if chars > maxChars {
		return fmt.Errorf(`must not be longer than %d characters`, maxChars)
	}

	return nil
}
