package options

import (
	"fmt"
)

// format: en_US
type Locale string

const LocaleDefault Locale = "en_US"

func (l Locale) String() string {
	return string(l)
}

func (l Locale) Validate() error {
	if l == "" {
		return fmt.Errorf("invalid locale: empty")
	}

	if len(l) != 5 {
		return fmt.Errorf("invalid locale: isn't 5 characters long")
	}

	if !(('a' <= l[0] && l[0] <= 'z') && ('a' <= l[1] && l[1] <= 'z')) {
		return fmt.Errorf("invalid locale: first two characters must be lowercase ASCII letters")
	}

	if !(('A' <= l[3] && l[3] <= 'Z') && ('A' <= l[4] && l[4] <= 'Z')) {
		return fmt.Errorf("invalid locale: last two characters must be uppercase ASCII letters")
	}

	if l[2] != '_' {
		return fmt.Errorf("invalid locale: third character must be underscore")
	}

	return nil
}

func StringToLocale(s string) (Locale, error) {
	l := Locale(s)
	if err := l.Validate(); err != nil {
		return "", err
	}

	return l, nil
}
