package models

type ErrorResponse struct {
	ErrorCode    int  `json:"error_code"`
	ErrorMessage any  `json:"error_message"`
	Success      bool `json:"success"`
}

func NewErrorResponse(code int, message error) *ErrorResponse {
	return &ErrorResponse{
		ErrorCode:    code,
		ErrorMessage: message,
		Success:      false,
	}
}

func NewErrorsResponse(code int, message string) *ErrorResponse {
	return &ErrorResponse{
		ErrorCode:    code,
		ErrorMessage: message,
		Success:      false,
	}
}

type SuccessResponse struct {
	Success bool `json:"success"`
	Data    any  `json:"data"`
}

func NewSuccessResponse(data any) *SuccessResponse {
	return &SuccessResponse{
		Success: true,
		Data:    data,
	}
}
