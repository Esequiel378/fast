// Package validator
package validator

type Validator interface {
	// ValidateStruct validates a struct using the validate struct tag.
	// The input should be a pointer to a struct, and it should have exported fields.
	ValidateStruct(input any) error
	// Translate translates the error returned by ValidateStruct to a slice of errors.
	// that can be used to return to the user.
	Translate(err error) []Error
}
