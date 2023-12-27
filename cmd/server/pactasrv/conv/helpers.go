package conv

import "time"

func ptr[T any](t T) *T {
	return &t
}

func strPtr[T ~string](t T) *string {
	return ptr(string(t))
}

func ifNil[T any](t *T, fallback T) T {
	if t == nil {
		return fallback
	}
	return *t
}

func timeToNilable(t time.Time) *time.Time {
	if t.IsZero() {
		return nil
	}
	return &t
}

func stringToNilable[T ~string](t T) *string {
	if t == "" {
		return nil
	}
	s := string(t)
	return &s
}

func convAll[I any, O any](is []I, f func(I) (O, error)) ([]O, error) {
	os := make([]O, len(is))
	for i, v := range is {
		o, err := f(v)
		if err != nil {
			return nil, err
		}
		os[i] = o
	}
	return os, nil
}

func dereferenceAll[T any](ts []*T) []T {
	result := make([]T, len(ts))
	for i, t := range ts {
		result[i] = *t
	}
	return result
}
