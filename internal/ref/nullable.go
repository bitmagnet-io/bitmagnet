package ref

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
