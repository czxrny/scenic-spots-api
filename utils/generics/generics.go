package generics

func DereferenceAll[T any](in []*T) []T {
	out := make([]T, 0, len(in))
	for _, ptr := range in {
		if ptr != nil {
			out = append(out, *ptr)
		}
	}
	return out
}
