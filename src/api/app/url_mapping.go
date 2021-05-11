package app

import (
	"github.com/PTLam25/git-hub-api-microservice/src/api/controllers/polo"
	"github.com/PTLam25/git-hub-api-microservice/src/api/controllers/repositories"
)

func mapUrls() {
	router.GET("/marco", polo.Marco)

	router.POST("/repository", repositories.CreateRepo)
	router.POST("/repositories", repositories.CreateRepos)
}
