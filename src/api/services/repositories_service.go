package services

//  Сервис содержит в себе всю бизнес логику для обработки запроса.
import (
	"github.com/PTLam25/git-hub-api-microservice/src/api/config"
	"github.com/PTLam25/git-hub-api-microservice/src/api/domain/github"
	"github.com/PTLam25/git-hub-api-microservice/src/api/domain/repositories"
	"github.com/PTLam25/git-hub-api-microservice/src/api/providers/github_provider"
	"github.com/PTLam25/git-hub-api-microservice/src/api/utils/errors"
	"strings"
)

type repoService struct {
}

type reposServiceInterface interface {
	CreateRepo(request repositories.CreateRepoRequest) (*repositories.CreateResponse, errors.ApiError)
}

var (
	RepositoryService reposServiceInterface
)

func init() {
	RepositoryService = &repoService{}
}

func (rp *repoService) CreateRepo(input repositories.CreateRepoRequest) (*repositories.CreateResponse, errors.ApiError) {
	// 1) валидация входных данных
	input.Name = strings.TrimSpace(input.Name)

	if input.Name == "" {
		return nil, errors.NewBadRequestError(
			"invalid repository name")
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

	result := repositories.CreateResponse{
		Id:    response.Id,
		Name:  response.Name,
		Owner: response.Owner.Login,
	}

	return &result, nil

}
