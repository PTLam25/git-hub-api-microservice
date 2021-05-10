package errors

// Интерфейс ошибок приложении, которые можно вернуть пользователю
import (
	"net/http"
)

type ApiError interface {
	Status() int
	Message() string
	Error() string
}

type apiError struct {
	AStatus  int    `json:"status"`
	AMessage string `json:"message"`
	// omitempty - не показывать поле, если его нет
	AnError string `json:"error,omitempty"`
}

func (ae *apiError) Status() int {
	return ae.AStatus
}

func (ae *apiError) Message() string {
	return ae.AMessage
}

func (ae *apiError) Error() string {
	return ae.AnError
}

func NewInternalServerError(message string) ApiError {
	return &apiError{
		AStatus:  http.StatusInternalServerError,
		AMessage: message,
	}
}

func NewApiError(statusCode int, message string) ApiError {
	return &apiError{
		AStatus:  statusCode,
		AMessage: message,
	}
}

func NewNotFoundApiError(message string) ApiError {
	return &apiError{
		AStatus:  http.StatusNotFound,
		AMessage: message,
	}
}

func NewBadRequestError(message string) ApiError {
	return &apiError{
		AStatus:  http.StatusBadRequest,
		AMessage: message,
	}
}
