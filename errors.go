package voyage

import "errors"

var (
	ErrInputRequired     = errors.New("input is required")
	ErrInputTooLarge     = errors.New("input length must be less than or equal to 128")
	ErrDocumentsRequired = errors.New("documents are required")
	ErrDocumentsTooLarge = errors.New("documents length must be less than or equal to 1000")
	ErrQueryRequired     = errors.New("query is required")
	ErrModelRequired     = errors.New("model is required")
)
