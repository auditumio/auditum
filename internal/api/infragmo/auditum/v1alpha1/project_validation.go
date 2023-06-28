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
