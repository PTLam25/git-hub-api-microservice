package github_provider

import (
	"encoding/json"
	"fmt"
	"github.com/PTLam25/git-hub-api-microservice/src/api/client/restclient"
	"github.com/PTLam25/git-hub-api-microservice/src/api/domain/github"
	"io/ioutil"
	"log"
	"net/http"
)

const (
	headerAuthorization       = "Authorization"
	headerAuthorizationFormat = "token %s"

	urlCreateRepo = "https://api.github.com/user/repos"
)

func getAuthorizationHeader(accessToken string) string {
	return fmt.Sprintf(headerAuthorizationFormat, accessToken)
}

func CreateRepo(assessToken string, request github.CreateRepoRequest) (*github.CreateRepoResponse, *github.GitHubErrorResponse) {
	// 1) создаем Header с accessToken
	headers := http.Header{}
	headers.Set(headerAuthorization, getAuthorizationHeader(assessToken))

	// 2) делаем запрос на GitHub
	response, err := restclient.Post(urlCreateRepo, request, headers)

	// 3) обработка ответа
	// когда не удалось получить ответа от GitHub
	if err != nil {
		// log error
		log.Printf("error when trying to create  new repo in github: %s", err.Error())
		return nil, &github.GitHubErrorResponse{
			StatusCode: http.StatusInternalServerError, Message: err.Error(),
		}
	}

	// читаем биты с тело ответа
	bytes, err := ioutil.ReadAll(response.Body)
	// ошибка при чтения битов
	if err != nil {
		// возвращаем 500, так как мы не знаем как обработать ответ
		return nil, &github.GitHubErrorResponse{
			StatusCode: http.StatusInternalServerError, Message: "invalid response body from GitHub",
		}
	}
	// закрываем поток чтения из тела ответа когда функция возвращает, чтобы избежать утечки памяти
	defer response.Body.Close()

	// обрабатываем ответ со статусом > 299
	if response.StatusCode > 299 {
		var errResponse github.GitHubErrorResponse
		if err := json.Unmarshal(bytes, &errResponse); err != nil {
			return nil, &github.GitHubErrorResponse{
				StatusCode: http.StatusInternalServerError, Message: "invalid json response body from GitHub",
			}
		}
		errResponse.StatusCode = response.StatusCode
		return nil, &errResponse
	}

	// десериализация данных в ответе
	var result github.CreateRepoResponse
	if err := json.Unmarshal(bytes, &result); err != nil {
		// log error
		log.Printf("error when trying to unmarshal create repo successful response from github: %s", err.Error())
		return nil, &github.GitHubErrorResponse{
			StatusCode: http.StatusInternalServerError, Message: "error unmarshaling GitHub create repo response",
		}
	}

	// отправляем результат
	return &result, nil
}
