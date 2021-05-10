package github_provider

import (
	"errors"
	"github.com/PTLam25/git-hub-api-microservice/src/api/client/restclient"
	"github.com/PTLam25/git-hub-api-microservice/src/api/domain/github"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"net/http"
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

func TestConstants(t *testing.T) {
	assert.EqualValues(t, "Authorization", headerAuthorization)
	assert.EqualValues(t, "token %s", headerAuthorizationFormat)
	assert.EqualValues(t, "https://api.github.com/user/repos", urlCreateRepo)
}

func TestGetAuthorizationHeader(t *testing.T) {
	// тестируем функцию getAuthorizationHeader
	header := getAuthorizationHeader("abc123")

	assert.EqualValues(t, "token abc123", header)
}

func TestCreateRepoErrorRestClient(t *testing.T) {
	restclient.FlushMockUps()
	// добавим фейк ответ в список
	restclient.AddMockUp(restclient.Mock{
		Url:        "https://api.github.com/user/repos",
		HttpMethod: http.MethodPost,
		Err:        errors.New("invalid restclient response"),
	})

	response, err := CreateRepo("", github.CreateRepoRequest{})

	assert.Nil(t, response)
	assert.NotNil(t, err)
	assert.EqualValues(t, "invalid restclient response", err.Message)
}

func TestCreateRepoErrorInvalidResponseBody(t *testing.T) {
	restclient.FlushMockUps()

	// рандомное тело ответа
	invalidCloser, _ := os.Open("-asf3")
	// добавим фейк ответ в список
	restclient.AddMockUp(restclient.Mock{
		Url:        "https://api.github.com/user/repos",
		HttpMethod: http.MethodPost,
		Response: &http.Response{
			StatusCode: http.StatusCreated,
			Body:       invalidCloser,
		},
	})

	response, err := CreateRepo("", github.CreateRepoRequest{})

	assert.Nil(t, response)
	assert.NotNil(t, err)
	assert.EqualValues(t, http.StatusInternalServerError, err.StatusCode)
	assert.EqualValues(t, "invalid response body from GitHub", err.Message)

	// стоп режим тестирования
	restclient.StopMockUps()
}

func TestCreateRepoInvalidErrorInterface(t *testing.T) {
	// включаем режим тестирования для restclient
	restclient.StartMockUps()
	restclient.FlushMockUps()

	restclient.AddMockUp(restclient.Mock{
		Url:        "https://api.github.com/user/repos",
		HttpMethod: http.MethodPost,
		Response: &http.Response{
			StatusCode: http.StatusUnauthorized,
			Body:       ioutil.NopCloser(strings.NewReader(`{"message": 1}`)),
		},
	})

	response, err := CreateRepo("", github.CreateRepoRequest{})

	assert.Nil(t, response)
	assert.NotNil(t, err)
	assert.EqualValues(t, http.StatusInternalServerError, err.StatusCode)
	assert.EqualValues(t, "invalid json response body from GitHub", err.Message)
}

func TestCreateRepoUnauthorized(t *testing.T) {
	// включаем режим тестирования для restclient
	restclient.StartMockUps()
	restclient.FlushMockUps()

	restclient.AddMockUp(restclient.Mock{
		Url:        "https://api.github.com/user/repos",
		HttpMethod: http.MethodPost,
		Response: &http.Response{
			StatusCode: http.StatusUnauthorized,
			Body:       ioutil.NopCloser(strings.NewReader(`{"message": "Requires authentication","documentation_url": "https://developer.github.com/v3/repos/#create"}`)),
		},
	})

	response, err := CreateRepo("", github.CreateRepoRequest{})

	assert.Nil(t, response)
	assert.NotNil(t, err)
	assert.EqualValues(t, http.StatusUnauthorized, err.StatusCode)
	assert.EqualValues(t, "Requires authentication", err.Message)
}

func TestCreateRepoInvalidSuccessResponse(t *testing.T) {
	// включаем режим тестирования для restclient
	restclient.StartMockUps()
	restclient.FlushMockUps()

	restclient.AddMockUp(restclient.Mock{
		Url:        "https://api.github.com/user/repos",
		HttpMethod: http.MethodPost,
		Response: &http.Response{
			StatusCode: http.StatusCreated,
			Body:       ioutil.NopCloser(strings.NewReader(`{"id": "123"}`)),
		},
	})

	response, err := CreateRepo("", github.CreateRepoRequest{})

	assert.Nil(t, response)
	assert.NotNil(t, err)
	assert.EqualValues(t, http.StatusInternalServerError, err.StatusCode)
	assert.EqualValues(t, "error unmarshaling GitHub create repo response", err.Message)
}

func TestCreateRepoNoError(t *testing.T) {
	restclient.FlushMockUps()

	restclient.AddMockUp(restclient.Mock{
		Url:        "https://api.github.com/user/repos",
		HttpMethod: http.MethodPost,
		Response: &http.Response{
			StatusCode: http.StatusCreated,
			Body:       ioutil.NopCloser(strings.NewReader(`{"id": 123,"name": "golang-tutorial","full_name": "federicoleon/golang-tutorial"}`)),
		},
	})

	response, err := CreateRepo("", github.CreateRepoRequest{})

	assert.Nil(t, err)
	assert.NotNil(t, response)
	assert.EqualValues(t, 123, response.Id)
	assert.EqualValues(t, "golang-tutorial", response.Name)
	assert.EqualValues(t, "federicoleon/golang-tutorial", response.FullName)
}
