package repositories

import (
	"github.com/PTLam25/git-hub-api-microservice/src/api/domain/repositories"
	"github.com/PTLam25/git-hub-api-microservice/src/api/services"
	"github.com/PTLam25/git-hub-api-microservice/src/api/utils/errors"
	"github.com/gin-gonic/gin"
	"net/http"
)

func CreateRepo(c *gin.Context) {
	var request repositories.CreateRepoRequest

	// смотрим на тело запроса и сходится ли с структурой CreateRepoRequest, если нет, то ошибка
	if err := c.ShouldBindJSON(&request); err != nil {
		apiError := errors.NewBadRequestError("invalid json body")
		c.JSON(apiError.Status(), apiError)
		return
	}

	result, err := services.RepositoryService.CreateRepo(request)
	if err != nil {
		c.JSON(err.Status(), err)
		return
	}

	c.JSON(http.StatusCreated, result)
}

func CreateRepos(c *gin.Context) {
	var request []repositories.CreateRepoRequest

	// смотрим на тело запроса и сходится ли с структурой CreateRepoRequest, если нет, то ошибка
	if err := c.ShouldBindJSON(&request); err != nil {
		apiError := errors.NewBadRequestError("invalid json body")
		c.JSON(apiError.Status(), apiError)
		return
	}

	result, err := services.RepositoryService.CreateRepos(request)
	if err != nil {
		c.JSON(err.Status(), err)
		return
	}

	c.JSON(result.StatusCode, result)
}
