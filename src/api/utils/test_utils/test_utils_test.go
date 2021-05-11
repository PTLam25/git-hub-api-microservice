package test_utils

import (
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetMockedContext(t *testing.T) {
	// 1) Инициализация
	request := httptest.NewRequest(http.MethodGet, "http://localhost:8000/something", nil)
	request.Header = http.Header{"X-Mock": {"true"}}
	response := httptest.NewRecorder()

	// 2) Вызов функции для теста
	c := GetMockedContext(request, response)

	// 3) Валидация
	assert.EqualValues(t, http.MethodGet, c.Request.Method)
	assert.EqualValues(t, "8000", c.Request.URL.Port())
	assert.EqualValues(t, "/something", c.Request.URL.Path)
	assert.EqualValues(t, "http", c.Request.URL.Scheme)
	assert.EqualValues(t, 1, len(c.Request.Header))
	assert.EqualValues(t, "true", c.GetHeader("x-mock"))
	assert.EqualValues(t, "true", c.GetHeader("X-Mock"))
}
