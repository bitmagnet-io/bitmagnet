package model

import (
	"fmt"
	"gorm.io/gorm"
	"regexp"
)

var tagNameRegex = regexp.MustCompile(`^[a-z0-9]+(-[a-z0-9]+)*$`)

const tagNameMaxLength = 20

func ValidateTagName(name string) error {
	if !tagNameRegex.MatchString(name) {
		return fmt.Errorf("invalid tag name: '%s' (must be kebab-case)", name)
	}
	if len(name) > tagNameMaxLength {
		return fmt.Errorf("invalid tag name: '%s' (must be no longer than %d characters)", name, tagNameMaxLength)
	}
	return nil
}

func (t *TorrentTag) BeforeCreate(*gorm.DB) error {
	if validateErr := ValidateTagName(t.Name); validateErr != nil {
		return validateErr
	}
	return nil
}
