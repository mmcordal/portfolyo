package viewmodel

type Validation interface {
	Validate() []error
}

type SuccessResponse struct {
	Success   bool        `json:"success"`
	Message   string      `json:"message,omitempty"`
	Data      interface{} `json:"data"`
	DataCount int         `json:"data_count"`
}

type ErrorResponse struct {
	Success      bool   `json:"success"`
	ErrorMessage string `json:"error_message"`
}
