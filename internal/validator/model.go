package validator

// Error is the serializer for a single validation error.
type Error struct {
	Field   string `json:"field,omitempty"`
	Message string `json:"message,omitempty"`
}

// ParsingErrorSerializer is the serializer for parsing errors.
// Generally, this is returned when there is an error
// parsing the request body or query parameters.
//
// Example:
//
//	c.Status(fiber.StatusBadRequest).JSON(common.ParsingErrorSerializer{
//	  Error: "error parsing input data",
//	})
type ParsingErrorSerializer struct {
	Error string `json:"error"`
}

// InternalErrorSerializer is the serializer for internal errors.
// Generally, this is returned when there is an internal error
// in the application.
//
// Example:
//
//	c.Status(fiber.StatusInternalServerError).JSON(common.InternalErrorSerializer{
//	  ErrorID: h.logger.Error(err),
//	})
type InternalErrorSerializer struct {
	ErrorID string `json:"error_id"`
}

// ValidationErrorSerializer is the serializer for validation errors.
// Generally, this is returned when there is an error
// validating the input/output data.
//
// Example:
//
//	c.Status(fiber.StatusBadRequest).JSON(common.ValidationErrorSerializer{
//	  Errors: h.validator.Translate(err),
//	})
type ValidationErrorSerializer struct {
	Errors []Error `json:"errors"`
}
