package utils


type APIResponse struct {
	Success bool `json:"success"`
	Message string `json:"message"`
	Data any `json:"data,omitempty"`
	Error *ErrorInfo `json:"error,omitempty"`
}


type ErrorInfo struct {
	Code string `json:"code,omitempty"`
	Details string `json:"details,omitempty"`
}


func SuccessResponse(message string , data any) APIResponse{
	return APIResponse{
		Success: true,
		Message: message,
		Data: data,
	}
}



func ErrorResponse(message string,err error) APIResponse {
	response := APIResponse{
		Success: false,
		Message: message,
	}

	if err != nil{
		response.Error = &ErrorInfo{
			Details: err.Error(),
		}
	}


	return  response
}