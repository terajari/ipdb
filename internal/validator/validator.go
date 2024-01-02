package validator

import "slices"

type Validator struct {
	Errors map[string]string
}

func New() *Validator {
	return &Validator{Errors: make(map[string]string)}
}

func (v *Validator) Valid() bool {
	return len(v.Errors) == 0
}

func (v *Validator) AddError(key, message string) {
	if _, exists := v.Errors[key]; !exists {
		v.Errors[key] = message
	}
}

func (v *Validator) Check(ok bool, key, message string) {
	if !ok {
		v.AddError(key, message)
	}
}

func Unique[T comparable](data ...T) bool {
	keys := make(map[T]bool)
	for _, entry := range data {
		if _, value := keys[entry]; value {
			return false
		}
		keys[entry] = true
	}
	return true
}

func PermitedValues[T comparable](value T, permitedValues ...T) bool {
	return slices.Contains(permitedValues, value)
}
