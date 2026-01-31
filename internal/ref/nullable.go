package ref

import (
	"fmt"
	"io"
)

type Nullable struct {
	Ref
	Valid bool
}

func NewNullable(ref Ref) Nullable {
	return Nullable{
		Ref:   ref,
		Valid: true,
	}
}

func (n Nullable) Value() (Ref, bool) {
	return n.Ref, n.Valid
}

func (n Nullable) MarshalGQL(w io.Writer) {
	if !n.Valid {
		_, _ = fmt.Fprintf(w, `null`)
		return
	}

	n.Ref.MarshalGQL(w)
}

func (n *Nullable) UnmarshalGQL(v any) error {
	vStr, ok := v.(string)
	if !ok {
		return fmt.Errorf("expected string, got %T", v)
	}

	if vStr == "" || vStr == "null" {
		n.Valid = false
		return nil
	}

	parsed, err := Parse(vStr)
	if err != nil {
		return err
	}

	*n = NewNullable(parsed)

	return nil
}
