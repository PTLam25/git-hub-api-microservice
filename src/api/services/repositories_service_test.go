package services

import (
	"github.com/PTLam25/git-hub-api-microservice/src/api/client/restclient"
	"github.com/PTLam25/git-hub-api-microservice/src/api/domain/repositories"
	"github.com/PTLam25/git-hub-api-microservice/src/api/utils/errors"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"sync"
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

func TestCreateRepoConcurrentInvalidRequest(t *testing.T) {
	// 1) Инициализация
	request := repositories.CreateRepoRequest{}
	output := make(chan repositories.CreateRepositoriesResult)
	service := repoService{}

	// 2) Вызов функции для теста в корутине
	go service.createRepoConcurrent(request, output)
	result := <-output

	// 3) Валидация
	assert.NotNil(t, result)
	assert.Nil(t, result.Response)
	assert.NotNil(t, result.Error)
	assert.EqualValues(t, http.StatusBadRequest, result.Error.Status())
	assert.EqualValues(t, "invalid repository name", result.Error.Message())
}

func TestCreateRepoConcurrentErrorFromGithub(t *testing.T) {
	// 1) Инициализация
	// текст кейс когда все ОК
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
	output := make(chan repositories.CreateRepositoriesResult)
	service := repoService{}

	// 2) Вызов функции для теста в корутине
	go service.createRepoConcurrent(request, output)
	result := <-output

	// 3) Валидация
	assert.NotNil(t, result)
	assert.Nil(t, result.Response)
	assert.NotNil(t, result.Error)
	assert.EqualValues(t, http.StatusUnauthorized, result.Error.Status())
	assert.EqualValues(t, "Requires authentication", result.Error.Message())
}

func TestCreateRepoConcurrentNoError(t *testing.T) {
	// 1) Инициализация
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
	output := make(chan repositories.CreateRepositoriesResult)
	service := repoService{}

	// 2) Вызов функции для теста в корутине
	go service.createRepoConcurrent(request, output)
	result := <-output

	// 3) Валидация
	assert.NotNil(t, result)
	assert.Nil(t, result.Error)
	assert.NotNil(t, result.Response)
	assert.EqualValues(t, 123, result.Response.Id)
	assert.EqualValues(t, "", result.Response.Name)
	assert.EqualValues(t, "", result.Response.Owner)
}

func TestHandleRepoResults(t *testing.T) {
	// 1) инициализация
	input := make(chan repositories.CreateRepositoriesResult)
	output := make(chan repositories.CreateReposResponse)
	defer close(output)
	var wg sync.WaitGroup
	service := repoService{}

	// 2) вызов функции для теста
	go service.handleRepoResults(&wg, input, output)
	wg.Add(1)
	go func() {
		input <- repositories.CreateRepositoriesResult{
			Error: errors.NewBadRequestError("invalid repository name"),
		}
	}()

	wg.Wait()
	close(input)
	result := <-output

	// 3) валидация
	assert.NotNil(t, result)
	assert.EqualValues(t, 0, result.StatusCode)
	assert.EqualValues(t, 1, len(result.Results))

	assert.NotNil(t, result.Results[0].Error)
	assert.EqualValues(t, http.StatusBadRequest, result.Results[0].Error.Status())
	assert.EqualValues(t, "invalid repository name", result.Results[0].Error.Message())
}

func TestCreateReposInvalidRequest(t *testing.T) {
	// 1) инициализация
	requests := []repositories.CreateRepoRequest{
		{},
		{Name: "    "},
	}

	// 2) вызов функции для теста
	result, err := RepositoryService.CreateRepos(requests)

	// 3) валидация
	assert.Nil(t, err)
	assert.NotNil(t, result)
	assert.EqualValues(t, 2, len(result.Results))
	assert.EqualValues(t, http.StatusBadRequest, result.StatusCode)

	assert.EqualValues(t, http.StatusBadRequest, result.Results[0].Error.Status())
	assert.EqualValues(t, "invalid repository name", result.Results[0].Error.Message())

	assert.EqualValues(t, http.StatusBadRequest, result.Results[1].Error.Status())
	assert.EqualValues(t, "invalid repository name", result.Results[1].Error.Message())
}

func TestCreateReposOneSuccessOneFail(t *testing.T) {
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

	requests := []repositories.CreateRepoRequest{
		{},
		{Name: "testing"},
	}

	// 2) вызов функции для теста
	result, err := RepositoryService.CreateRepos(requests)

	// 3) валидация
	assert.Nil(t, err)
	assert.NotNil(t, result)
	assert.EqualValues(t, 2, len(result.Results))
	assert.EqualValues(t, http.StatusPartialContent, result.StatusCode)

	for _, result := range result.Results {
		if result.Error != nil {
			assert.EqualValues(t, http.StatusBadRequest, result.Error.Status())
			assert.EqualValues(t, "invalid repository name", result.Error.Message())
			continue
		}

		assert.EqualValues(t, 123, result.Response.Id)
		assert.EqualValues(t, "", result.Response.Name)
		assert.EqualValues(t, "", result.Response.Owner)
	}
}

func TestCreateReposAllSuccess(t *testing.T) {
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

	requests := []repositories.CreateRepoRequest{
		{Name: "testing"},
		{Name: "testing"},
	}

	// 2) вызов функции для теста
	result, err := RepositoryService.CreateRepos(requests)

	// 3) валидация
	assert.Nil(t, err)
	assert.NotNil(t, result)
	assert.EqualValues(t, 2, len(result.Results))
	assert.EqualValues(t, http.StatusCreated, result.StatusCode)

	for _, result := range result.Results {
		assert.Nil(t, result.Error)
		assert.EqualValues(t, 123, result.Response.Id)
		assert.EqualValues(t, "", result.Response.Name)
		assert.EqualValues(t, "", result.Response.Owner)
	}
}
