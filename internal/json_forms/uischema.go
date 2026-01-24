package json_forms

import (
	"encoding/json"
	"errors"
	"strings"

	"github.com/bitmagnet-io/bitmagnet/pkg/json_schema"
)

type UISchema Element

func UnmarshalUISchema(data []byte) (UISchema, error) {
	element, err := unmarshalElement(data)
	if err != nil {
		return nil, err
	}

	return UISchema(element), nil
}

type baseElement struct {
	Type string `json:"type"`
}

func unmarshalElement(data []byte) (Element, error) {
	var base baseElement
	if err := json.Unmarshal(data, &base); err != nil {
		return nil, err
	}

	switch base.Type {
	case "Control":
		var control Control
		if err := json.Unmarshal(data, &control); err != nil {
			return nil, err
		}

		return control, nil
	case string(LayoutTypeVertical), string(LayoutTypeHorizontal), string(LayoutTypeGroup):
		var layout Layout
		if err := json.Unmarshal(data, &layout); err != nil {
			return nil, err
		}

		return layout, nil
	case "Categorization":
		var categorization Categorization
		if err := json.Unmarshal(data, &categorization); err != nil {
			return nil, err
		}

		return categorization, nil
	default:
		return nil, errors.New("unknown element type: " + base.Type)
	}
}

type Element interface {
	element()
}

type Control struct {
	Type    ControlStr `json:"type"`
	Scope   string     `json:"scope"`
	Label   *string    `json:"label,omitempty"`
	Options *Options   `json:"options,omitempty"`
	Rule    *Rule      `json:"rule,omitempty"`
}

func (Control) element() {}

type ControlStr string

func (c *ControlStr) UnmarshalJSON(data []byte) error {
	var str string
	if err := json.Unmarshal(data, &str); err != nil {
		return err
	}

	if str != string(ControlControl) {
		return errors.New("invalid control type")
	}

	*c = ControlControl

	return nil
}

const (
	ControlControl ControlStr = "Control"
)

type Options struct {
	Detail           Detail  `json:"detail,omitempty"`
	ShowSortButtons  *bool   `json:"showSortButtons,omitempty"`
	ElementLabelProp *string `json:"elementLabelProp,omitempty"`
	Format           *Format `json:"format,omitempty"`
	ReadOnly         *bool   `json:"readOnly,omitempty"`
}

func (o *Options) UnmarshalJSON(data []byte) error {
	type Alias Options

	aux := &struct {
		Detail json.RawMessage `json:"detail,omitempty"`
		*Alias
	}{
		Alias: (*Alias)(o),
	}

	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	if len(aux.Detail) > 0 {
		d, err := unmarshalDetail(aux.Detail)
		if err != nil {
			return err
		}

		o.Detail = d
	}

	return nil
}

func unmarshalDetail(data []byte) (Detail, error) {
	if string(data) == "null" {
		//nolint:nilnil
		return nil, nil
	}

	var detailStr DetailStr
	if err := json.Unmarshal(data, &detailStr); err == nil {
		return detailStr, nil
	}

	if element, err := unmarshalElement(data); err == nil {
		return DetailElement{Element: element}, nil
	}

	return nil, errors.New("invalid detail value")
}

type Detail interface {
	detail()
}

type DetailStr string

func (d *DetailStr) UnmarshalJSON(data []byte) error {
	var str string
	if err := json.Unmarshal(data, &str); err != nil {
		return err
	}

	switch str {
	case string(DetailDefault), string(DetailGenerated), string(DetailRegistered):
		*d = DetailStr(str)
		return nil
	default:
		return errors.New("invalid detail value")
	}
}

func (DetailStr) detail() {}

const (
	DetailDefault    DetailStr = "DEFAULT"
	DetailGenerated  DetailStr = "GENERATED"
	DetailRegistered DetailStr = "REGISTERED"
)

type DetailElement struct {
	Element
}

func (d DetailElement) MarshalJSON() ([]byte, error) {
	return json.Marshal(d.Element)
}

func (DetailElement) detail() {}

type Format string

func (f *Format) UnmarshalJSON(data []byte) error {
	var str string
	if err := json.Unmarshal(data, &str); err != nil {
		return err
	}

	if str != string(FormatRadio) {
		return errors.New("invalid format value")
	}

	*f = FormatRadio

	return nil
}

const (
	FormatRadio Format = "radio"
)

type Layout struct {
	Type     LayoutType `json:"type"`
	Label    *string    `json:"label,omitempty"`
	Elements []Element  `json:"elements"`
	Rule     *Rule      `json:"rule,omitempty"`
}

func (l *Layout) UnmarshalJSON(data []byte) error {
	type Alias Layout

	aux := &struct {
		Elements []json.RawMessage `json:"elements"`
		*Alias
	}{
		Alias: (*Alias)(l),
	}

	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	l.Elements = make([]Element, len(aux.Elements))
	for i, elemData := range aux.Elements {
		elem, err := unmarshalElement(elemData)
		if err != nil {
			return err
		}

		l.Elements[i] = elem
	}

	return nil
}

func (Layout) element() {}

type LayoutType string

func (l *LayoutType) UnmarshalJSON(data []byte) error {
	var str string
	if err := json.Unmarshal(data, &str); err != nil {
		return err
	}

	switch str {
	case string(LayoutTypeVertical), string(LayoutTypeHorizontal), string(LayoutTypeGroup):
		*l = LayoutType(str)
		return nil
	default:
		return errors.New("invalid layout type")
	}
}

const (
	LayoutTypeVertical   LayoutType = "VerticalLayout"
	LayoutTypeHorizontal LayoutType = "HorizontalLayout"
	LayoutTypeGroup      LayoutType = "Group"
)

type Categorization struct {
	Type     CategorizationStr `json:"type"`
	Elements []Category        `json:"elements"`
	Rule     *Rule             `json:"rule,omitempty"`
}

func (c *Categorization) UnmarshalJSON(data []byte) error {
	type Alias Categorization

	aux := &struct {
		Elements []json.RawMessage `json:"elements"`
		*Alias
	}{
		Alias: (*Alias)(c),
	}

	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	c.Elements = make([]Category, len(aux.Elements))
	for i, elemData := range aux.Elements {
		var category Category
		if err := json.Unmarshal(elemData, &category); err != nil {
			return err
		}

		c.Elements[i] = category
	}

	return nil
}

func (Categorization) element() {}

type CategorizationStr string

func (c *CategorizationStr) UnmarshalJSON(data []byte) error {
	var str string
	if err := json.Unmarshal(data, &str); err != nil {
		return err
	}

	if str != string(CategorizationCategorization) {
		return errors.New("invalid categorization type")
	}

	*c = CategorizationCategorization

	return nil
}

const (
	CategorizationCategorization CategorizationStr = "Categorization"
)

type Category struct {
	Type     CategoryStr `json:"type"`
	Label    *string     `json:"label,omitempty"`
	Elements []Element   `json:"elements"`
	Rule     *Rule       `json:"rule,omitempty"`
}

func (c *Category) UnmarshalJSON(data []byte) error {
	type Alias Category

	aux := &struct {
		Elements []json.RawMessage `json:"elements"`
		*Alias
	}{
		Alias: (*Alias)(c),
	}

	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	c.Elements = make([]Element, len(aux.Elements))
	for i, elemData := range aux.Elements {
		elem, err := unmarshalElement(elemData)
		if err != nil {
			return err
		}

		c.Elements[i] = elem
	}

	return nil
}

func (Category) element() {}

type CategoryStr string

func (c *CategoryStr) UnmarshalJSON(data []byte) error {
	var str string
	if err := json.Unmarshal(data, &str); err != nil {
		return err
	}

	if str != string(CategoryCategory) {
		return errors.New("invalid category type")
	}

	*c = CategoryCategory

	return nil
}

const (
	CategoryCategory CategoryStr = "Category"
)

type Rule struct {
	Effect    RuleEffect `json:"effect"`
	Condition Condition  `json:"condition"`
}

type RuleEffect string

func (r *RuleEffect) UnmarshalJSON(data []byte) error {
	var str string
	if err := json.Unmarshal(data, &str); err != nil {
		return err
	}

	upper := strings.ToUpper(str)

	switch upper {
	case string(RuleEffectHide), string(RuleEffectShow), string(RuleEffectEnable), string(RuleEffectDisable):
		*r = RuleEffect(upper)
		return nil
	default:
		return errors.New("invalid rule effect")
	}
}

const (
	RuleEffectHide    RuleEffect = "HIDE"
	RuleEffectShow    RuleEffect = "SHOW"
	RuleEffectEnable  RuleEffect = "ENABLE"
	RuleEffectDisable RuleEffect = "DISABLE"
)

type Condition struct {
	Scope             string                 `json:"scope"`
	Schema            json_schema.JSONSchema `json:"schema"`
	FailWhenUndefined *bool                  `json:"failWhenUndefined,omitempty"`
}
