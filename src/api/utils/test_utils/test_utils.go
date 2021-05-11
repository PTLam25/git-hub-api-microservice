package test_utils

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"net/http/httptest"
)

func GetMockedContext(request *http.Request, response *httptest.ResponseRecorder) *gin.Context {
	// создали фейк контекст gin-gonic для теста, который будет менять созданный response
	c, _ := gin.CreateTestContext(response)
	// добавляем в контекст объект request
	c.Request = request

	return c
}
