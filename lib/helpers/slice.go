package helpers

func SafeGet[T any](s []T, i int) (*T, error) {
	if i >= len(s) {
		return nil, Error{"Index out of bounds", ""}
	}

	return &s[i], nil
}

func SafeGetRange[T any](s []T, i, j int) ([]T, error) {
	if i > len(s) || j > len(s) {
		return nil, Error{"Index out of bounds", ""}
	}

	return s[i:j], nil
}
