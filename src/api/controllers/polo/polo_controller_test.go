package polo

import (
	"github.com/PTLam25/git-hub-api-microservice/src/api/utils/test_utils"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestConstants(t *testing.T) {
	// тестируем константы в файле. Это важно проверять,
	//так как при изменения константы программа может поменять свое поведение
	assert.EqualValues(t, "polo", polo)
}

func TestPolo(t *testing.T) {
	// 1) инициализация
	response := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodGet, "/marco", nil)
	c := test_utils.GetMockedContext(request, response)

	// 2) вызов функции для теста
	Marco(c)

	// 3) валидация
	assert.EqualValues(t, http.StatusOK, response.Code)
	assert.EqualValues(t, "polo", response.Body.String())
}
