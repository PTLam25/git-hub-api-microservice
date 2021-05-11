package services

//  Сервис содержит в себе всю бизнес логику для обработки запроса.
import (
	"github.com/PTLam25/git-hub-api-microservice/src/api/config"
	"github.com/PTLam25/git-hub-api-microservice/src/api/domain/github"
	"github.com/PTLam25/git-hub-api-microservice/src/api/domain/repositories"
	"github.com/PTLam25/git-hub-api-microservice/src/api/providers/github_provider"
	"github.com/PTLam25/git-hub-api-microservice/src/api/utils/errors"
	"net/http"
	"sync"
)

type repoService struct {
}

type reposServiceInterface interface {
	CreateRepo(request repositories.CreateRepoRequest) (*repositories.CreateRepoResponse, errors.ApiError)
	CreateRepos(request []repositories.CreateRepoRequest) (repositories.CreateReposResponse, errors.ApiError)
}

var (
	RepositoryService reposServiceInterface
)

func init() {
	RepositoryService = &repoService{}
}

func (rp *repoService) CreateRepo(input repositories.CreateRepoRequest) (*repositories.CreateRepoResponse, errors.ApiError) {
	// 1) валидация входных данных
	if err := input.Validate(); err != nil {
		return nil, err
	}

	// 3) подготовка данных для дальнейшего запроса
	request := github.CreateRepoRequest{
		Name:        input.Name,
		Description: input.Description,
		Private:     false,
	}

	// 2) обращения к внешнему АПИ за данными
	response, err := github_provider.CreateRepo(config.GetGithubAccessToken(), request)

	//  3) валидации ответа от АПИ и возвращения ответа клиенту
	if err != nil {
		return nil, errors.NewApiError(err.StatusCode, err.Message)
	}

	result := repositories.CreateRepoResponse{
		Id:    response.Id,
		Name:  response.Name,
		Owner: response.Owner.Login,
	}

	return &result, nil
}

func (rp *repoService) CreateRepos(requests []repositories.CreateRepoRequest) (repositories.CreateReposResponse, errors.ApiError) {
	// create channel for go routines
	input := make(chan repositories.CreateRepositoriesResult)
	output := make(chan repositories.CreateReposResponse)
	// закрываем канал output когда функция возвращает
	defer close(output)

	var wg sync.WaitGroup

	go rp.handleRepoResults(&wg, input, output)

	// run go routine for creating repo for each request in concurrency
	for _, currentRequest := range requests {
		// добавляем счетчик корутин в WaitGroup, счетчик уйдет когда корутина закончит работу
		wg.Add(1)
		go rp.createRepoConcurrent(currentRequest, input)
	}

	// просим не совершать функцию, а подождать пока все корутины закончат работу
	//(когда счетчик WaitGroup дойдет до 0, тогда продолжить выполнения кода
	wg.Wait()
	// закрываем канал input, так как там все корутины закончили работу
	close(input)

	// ждем результат с канала output
	result := <-output

	// проверяем сколько запроса было успешно выполнено
	successCreations := 0
	for _, current := range result.Results {
		if current.Response != nil {
			successCreations++
		}
	}
	if successCreations == 0 {
		// все запросы провалились берем ошибку, который вернул GitHub с 1-ого элемента
		result.StatusCode = result.Results[0].Error.Status()
	} else if successCreations == len(requests) {
		result.StatusCode = http.StatusCreated
	} else {
		result.StatusCode = http.StatusPartialContent
	}

	return result, nil
}

func (rp *repoService) handleRepoResults(wg *sync.WaitGroup, input chan repositories.CreateRepositoriesResult, output chan repositories.CreateReposResponse) {
	var results repositories.CreateReposResponse

	// начинаем слушать и получать все событии из канала input пока он открыт
	for incomingEvent := range input {
		repoResult := repositories.CreateRepositoriesResult{
			Response: incomingEvent.Response,
			Error:    incomingEvent.Error,
		}
		results.Results = append(results.Results, repoResult)
		// Уменьшаем счетчик корутин в WaitGroup на 1, так как корутина закончила выполнения
		wg.Done()
	}

	// когда канал input будет закрыт, то выйдем из цикла выше, отправляем результат в канал output
	output <- results
}

func (rp *repoService) createRepoConcurrent(input repositories.CreateRepoRequest, output chan repositories.CreateRepositoriesResult) {
	// function for creating single repo in go routine
	// валидация данных
	if err := input.Validate(); err != nil {
		output <- repositories.CreateRepositoriesResult{Error: err}
		return
	}

	// делаем запрос ПОСТ на создания репо
	result, err := rp.CreateRepo(input)

	// валидация результата запроса
	if err != nil {
		output <- repositories.CreateRepositoriesResult{Error: err}
		return
	}

	output <- repositories.CreateRepositoriesResult{Response: result}
}
