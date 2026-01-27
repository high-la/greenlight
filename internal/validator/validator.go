package validator

import (
	"regexp"
	"slices"
)

// Declare a regular expression for sanity checking the format of email ads
var (
	EmailRX = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")
)

// Define a new Validator type which contains a map of validation errors.
type Validator struct {
	Errors map[string]string
}

// New is a helper which creates anew Validator instance with an empty errors map.
func New() *Validator {
	return &Validator{Errors: make(map[string]string)}
}

// Valid returns tru if the errors map doesn't contain any entries.
func (v *Validator) Valid() bool {
	return len(v.Errors) == 0
}

// AddError adds an erros message to the map (so long as no entry already exists for
// the given key).
func (v *Validator) AddError(key, message string) {
	if _, exists := v.Errors[key]; !exists {
		v.Errors[key] = message
	}
}

// Check adds an error to the map only if a validation check is not ok.
func (v *Validator) Check(ok bool, key, message string) {
	if !ok {
		v.AddError(key, message)
	}
}

// Generic function which returns true if a specific value is in a lis of
// permitted values.
func PermittedValue[T comparable](value T, permittedValues ...T) bool {

	return slices.Contains(permittedValues, value)
}

// Matches returns true if a string value matches a specific regex pattern.
func Matches(value string, rx *regexp.Regexp) bool {

	return rx.MatchString(value)
}

// Generic function which returns true if all values in a slice are unique.
func Unique[T comparable](values []T) bool {

	uniqueValues := make(map[T]bool)

	for _, value := range values {
		uniqueValues[value] = true
	}

	return len(values) == len(uniqueValues)
}
