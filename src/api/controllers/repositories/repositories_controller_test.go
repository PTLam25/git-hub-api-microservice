package repositories

import (
	"encoding/json"
	"github.com/PTLam25/git-hub-api-microservice/src/api/client/restclient"
	"github.com/PTLam25/git-hub-api-microservice/src/api/domain/repositories"
	"github.com/PTLam25/git-hub-api-microservice/src/api/utils/errors"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
)

func TestMain(m *testing.M) {
	// первая функция, которая запускается при начале тестирования.
	// наподобие setup в Питоне

	// включаем режим тестирования для restclient
	restclient.StartMockUps()
	// run test и exit после заверщения
	os.Exit(m.Run())
}

func TestCreateRepoInvalidJsonRequest(t *testing.T) {
	// 1) Инициализация
	// создаем объект Response для теста
	response := httptest.NewRecorder()
	// создали фейк контекст gin-gonic для теста, который будет менять созданный response
	c, _ := gin.CreateTestContext(response)
	// добавляем в контекст объект request
	request, _ := http.NewRequest(http.MethodPost, "/repositories", strings.NewReader(``))
	c.Request = request

	// 2) вызываем функция для теста, которая делает запрос и поменяет данные response
	CreateRepo(c)

	// 3) Валидация
	// валидация статуса ответа в заголовке
	assert.EqualValues(t, http.StatusBadRequest, response.Code)
	// валидация тела ответа
	apiError, err := errors.NewApiErrorFromBytes(response.Body.Bytes())
	assert.Nil(t, err)
	assert.NotNil(t, apiError)
	assert.EqualValues(t, http.StatusBadRequest, apiError.Status())
	assert.EqualValues(t, "invalid json body", apiError.Message())
}

func TestCreateRepoErrorFromGitHub(t *testing.T) {
	// 1) Инициализация
	// создаем объект Response для теста
	response := httptest.NewRecorder()
	// создали фейк контекст gin-gonic для теста, который будет менять созданный response
	c, _ := gin.CreateTestContext(response)
	// добавляем в контекст объект request
	request, _ := http.NewRequest(http.MethodPost, "/repositories", strings.NewReader(`{"name": "testing"}`))
	c.Request = request

	restclient.FlushMockUps()
	restclient.AddMockUp(restclient.Mock{
		Url:        "https://api.github.com/user/repos",
		HttpMethod: http.MethodPost,
		Response: &http.Response{
			StatusCode: http.StatusUnauthorized,
			Body:       ioutil.NopCloser(strings.NewReader(`{"message": "Requires authentication","documentation_url": "https://developer.github.com/v3/repos/#create"}`)),
		},
	})
	// 2) вызываем функция для теста, которая делает запрос и поменяет данные response
	CreateRepo(c)

	// 3) Валидация
	// валидация статуса ответа в заголовке
	assert.EqualValues(t, http.StatusUnauthorized, response.Code)
	// валидация тела ответа
	apiError, err := errors.NewApiErrorFromBytes(response.Body.Bytes())
	assert.Nil(t, err)
	assert.NotNil(t, apiError)
	assert.EqualValues(t, http.StatusUnauthorized, apiError.Status())
	assert.EqualValues(t, "Requires authentication", apiError.Message())
}

func TestCreateRepoNoError(t *testing.T) {
	// 1) Инициализация
	// создаем объект Response для теста
	response := httptest.NewRecorder()
	// создали фейк контекст gin-gonic для теста, который будет менять созданный response
	c, _ := gin.CreateTestContext(response)
	// добавляем в контекст объект request
	request, _ := http.NewRequest(http.MethodPost, "/repositories", strings.NewReader(`{"name": "testing"}`))
	c.Request = request

	restclient.FlushMockUps()
	restclient.AddMockUp(restclient.Mock{
		Url:        "https://api.github.com/user/repos",
		HttpMethod: http.MethodPost,
		Response: &http.Response{
			StatusCode: http.StatusCreated,
			Body:       ioutil.NopCloser(strings.NewReader(`{"id": 123}`)),
		},
	})
	// 2) вызываем функция для теста, которая делает запрос и поменяет данные response
	CreateRepo(c)

	// 3) Валидация
	// валидация статуса ответа в заголовке
	assert.EqualValues(t, http.StatusCreated, response.Code)
	// валидация тела ответа
	var result repositories.CreateResponse
	err := json.Unmarshal(response.Body.Bytes(), &result)
	assert.Nil(t, err)
	assert.EqualValues(t, 123, result.Id)
	assert.EqualValues(t, "", result.Name)
	assert.EqualValues(t, "", result.Owner)
}
