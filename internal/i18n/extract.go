package i18n

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/bitmagnet-io/bitmagnet/internal/slice"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
	"gopkg.in/yaml.v3"
)

func Write(messages []*i18n.Message) error {
	for _, lang := range languages {
		filePath := filepath.Join(".", "i18n", fmt.Sprintf("%s.yaml", lang))

		var current map[string]fileMessage

		currentBytes, err := os.ReadFile(filePath)
		if err == nil {
			err = yaml.Unmarshal(currentBytes, &current)
		}

		if err != nil {
			current = make(fileMessages)
		}

		v, err := marshalValue(messages, current, lang == language.English)
		if err != nil {
			return err
		}

		content, err := yaml.Marshal(v)
		if err != nil {
			return fmt.Errorf("failed to marshal %s strings: %w", lang, err)
		}

		err = os.WriteFile(filePath, content, 0o666)
		if err != nil {
			return err
		}
	}

	return nil
}

type (
	fileMessage struct {
		Description *string `yaml:"description,omitempty"`
		Hash        *string `yaml:"hash,omitempty"`
		Zero        *string `yaml:"zero,omitempty"`
		One         *string `yaml:"one,omitempty"`
		Two         *string `yaml:"two,omitempty"`
		Few         *string `yaml:"few,omitempty"`
		Many        *string `yaml:"many,omitempty"`
		Other       *string `yaml:"other,omitempty"`
	}

	fileMessages map[string]fileMessage
)

func marshalValue(messages []*i18n.Message, current fileMessages, sourceLanguage bool) (fileMessages, error) {
	messageTemplates := slice.Map(messages, i18n.NewMessageTemplate)

	v := make(fileMessages, len(messageTemplates))
	for _, template := range messageTemplates {
		if template == nil {
			return nil, errors.New("template is empty")
		}

		if _, ok := v[template.ID]; ok {
			return nil, fmt.Errorf("duplicate message ID: %s", template.ID)
		}

		var m fileMessage

		curr := current[template.ID]

		if template.Description != "" {
			m.Description = &template.Description
		}

		if !sourceLanguage && template.Hash != "" {
			m.Hash = &template.Hash
		}

		for pluralForm, template := range template.PluralTemplates {
			var src string
			if sourceLanguage {
				src = template.Src
			}

			switch string(pluralForm) {
			case "zero":
				if !sourceLanguage && curr.Zero != nil {
					src = *curr.Zero
				}

				m.Zero = &src
			case "one":
				if !sourceLanguage && curr.One != nil {
					src = *curr.One
				}

				m.One = &src
			case "two":
				if !sourceLanguage && curr.Two != nil {
					src = *curr.Two
				}

				m.Two = &src
			case "few":
				if !sourceLanguage && curr.Few != nil {
					src = *curr.Few
				}

				m.Few = &src
			case "many":
				if !sourceLanguage && curr.Many != nil {
					src = *curr.Many
				}

				m.Many = &src
			case "other":
				if !sourceLanguage && curr.Other != nil {
					src = *curr.Other
				}

				m.Other = &src
			}
		}

		v[template.ID] = m
	}

	return v, nil
}
