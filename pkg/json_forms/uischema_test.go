package json_forms_test

import (
	"encoding/json"
	"testing"

	"github.com/bitmagnet-io/bitmagnet/pkg/json_forms"
	"github.com/stretchr/testify/require"
)

func TestUISchema(t *testing.T) {
	t.Parallel()

	for _, test := range []struct {
		input    string
		expected json_forms.Element
	}{
		{
			input: `{
  "type": "Control",
  "scope": "#/properties/name"
}`,
			expected: json_forms.Control{
				Type:  "Control",
				Scope: "#/properties/name",
			},
		},
		{
			input: `{
  "type": "Control",
  "scope": "#/properties/name",
  "label": "First name",
  "options": {
    "detail" : "DEFAULT"
  }
}`,
			expected: json_forms.Control{
				Type:  "Control",
				Scope: "#/properties/name",
				Label: func() *string { s := "First name"; return &s }(),
				Options: &json_forms.Options{
					Detail: json_forms.DetailDefault,
				},
			},
		},
		{
			input: `{
  "type": "HorizontalLayout",
  "elements": [
    {
      "type": "Control",
      "scope": "#/properties/name"
    }
  ]
}`,
			expected: json_forms.Layout{
				Type: json_forms.LayoutTypeHorizontal,
				Elements: []json_forms.Element{
					json_forms.Control{
						Type:  "Control",
						Scope: "#/properties/name",
					},
				},
			},
		},
		{
			input: `{
  "type": "Control",
  "scope": "#/properties/name",
  "label": "First name",
  "options": {
    "detail" : {
      "type": "HorizontalLayout",
      "elements": [
        {
          "type": "Control",
          "scope": "#/properties/name"
        }
      ]
    }
  }
}`,
			expected: json_forms.Control{
				Type:  "Control",
				Scope: "#/properties/name",
				Label: func() *string { s := "First name"; return &s }(),
				Options: &json_forms.Options{
					Detail: json_forms.DetailElement{
						Element: json_forms.Layout{
							Type: json_forms.LayoutTypeHorizontal,
							Elements: []json_forms.Element{
								json_forms.Control{
									Type:  "Control",
									Scope: "#/properties/name",
								},
							},
						},
					},
				},
			},
		},
	} {
		t.Run(test.input, func(t *testing.T) {
			t.Parallel()

			schema, err := json_forms.UnmarshalUISchema([]byte(test.input))
			require.NoError(t, err)
			require.Equal(t, test.expected, schema)

			jsonEncoded, err := json.Marshal(schema)
			t.Log(string(jsonEncoded))
			require.NoError(t, err)
			require.JSONEq(t, test.input, string(jsonEncoded))
		})
	}
}
