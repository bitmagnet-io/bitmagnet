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

func (r Nullable) MarshalGQL(w io.Writer) {
	if !r.Valid {
		_, _ = fmt.Fprintf(w, `null`)
		return
	}
	r.Ref.MarshalGQL(w)
}

func (r *Nullable) UnmarshalGQL(v any) error {
	vStr, ok := v.(string)
	if !ok {
		return fmt.Errorf("expected string, got %T", v)
	}

	if vStr == "" || vStr == "null" {
		r.Valid = false
		return nil
	}

	parsed, err := Parse(vStr)
	if err != nil {
		return err
	}

	*r = NewNullable(parsed)

	return nil
}
