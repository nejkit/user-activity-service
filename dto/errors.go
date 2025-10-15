package dto

type ErrorResponseDto struct {
	Code    ErrorCode `json:"code"`
	Message string    `json:"message,omitempty"`
}

type ErrorCode int

const (
	ErrorCodeNone ErrorCode = iota
	ErrorCodeInvalidRequest
	ErrorCodeAccountNotFound
)
