package response

type APIResponse struct {
    Success bool        `json:"success"`
    Message string      `json:"message"`
    Data    interface{} `json:"data,omitempty"`
    Error   string      `json:"error,omitempty"`
}

func Success(message string, data interface{}) APIResponse {
    return APIResponse{
        Success: true,
        Message: message,
        Data:    data,
    }
}

func Error(message, error string) APIResponse {
    return APIResponse{
        Success: false,
        Message: message,
        Error:   error,
    }
}