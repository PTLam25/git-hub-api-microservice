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
	status  int    `json:"status"`
	message string `json:"message"`
	// omitempty - не показывать поле, если его нет
	error string `json:"error,omitempty"`
}

func (ae *apiError) Status() int {
	return ae.status
}

func (ae *apiError) Message() string {
	return ae.message
}

func (ae *apiError) Error() string {
	return ae.error
}

func NewInternalServerError(message string) ApiError {
	return &apiError{
		status:  http.StatusInternalServerError,
		message: message,
	}
}

func NewNotFoundApiError(message string) ApiError {
	return &apiError{
		status:  http.StatusNotFound,
		message: message,
	}
}

func NewBadRequestError(message string) ApiError {
	return &apiError{
		status:  http.StatusBadRequest,
		message: message,
	}
}
