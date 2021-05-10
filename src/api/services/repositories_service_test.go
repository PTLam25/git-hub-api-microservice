package services

import (
	"github.com/PTLam25/git-hub-api-microservice/src/api/client/restclient"
	"github.com/PTLam25/git-hub-api-microservice/src/api/domain/repositories"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"testing"
)

func TestMain(m *testing.M) {
	restclient.StartMockUps()
	os.Exit(m.Run())
}

func TestCreateRepoInvalidInputName(t *testing.T) {
	// текст кейс когда неправильные данные от клиента (пустая строка в name)
	// 1) инициализация
	request := repositories.CreateRepoRequest{}

	// 2) вызов функции для теста
	result, err := RepositoryService.CreateRepo(request)

	// 3) валидация
	assert.Nil(t, result)
	assert.NotNil(t, err)
	assert.EqualValues(t, http.StatusBadRequest, err.Status())
	assert.EqualValues(t, "invalid repository name", err.Message())
}

func TestCreateRepoErrorFromGitHub(t *testing.T) {
	// текст кейс когда ошибка от Github
	// 1) инициализация
	restclient.FlushMockUps()
	restclient.AddMockUp(restclient.Mock{
		Url:        "https://api.github.com/user/repos",
		HttpMethod: http.MethodPost,
		Response: &http.Response{
			StatusCode: http.StatusUnauthorized,
			Body:       ioutil.NopCloser(strings.NewReader(`{"message": "Requires authentication","documentation_url": "https://developer.github.com/docs"}`)),
		},
	})
	request := repositories.CreateRepoRequest{Name: "testing"}

	// 2) вызов функции для теста
	result, err := RepositoryService.CreateRepo(request)

	// 3) валидация
	assert.Nil(t, result)
	assert.NotNil(t, err)
	assert.EqualValues(t, http.StatusUnauthorized, err.Status())
	assert.EqualValues(t, "Requires authentication", err.Message())
}

func TestCreateRepoNoError(t *testing.T) {
	// текст кейс когда все ОК
	// 1) инициализация
	restclient.FlushMockUps()
	restclient.AddMockUp(restclient.Mock{
		Url:        "https://api.github.com/user/repos",
		HttpMethod: http.MethodPost,
		Response: &http.Response{
			StatusCode: http.StatusCreated,
			Body:       ioutil.NopCloser(strings.NewReader(`{"id": 123}`)),
		},
	})

	request := repositories.CreateRepoRequest{Name: "testing"}

	// 2) вызов функции для теста
	result, err := RepositoryService.CreateRepo(request)

	// 3) валидация
	assert.Nil(t, err)
	assert.NotNil(t, result)
	assert.EqualValues(t, 123, result.Id)
	assert.EqualValues(t, "", result.Name)
	assert.EqualValues(t, "", result.Owner)
}
