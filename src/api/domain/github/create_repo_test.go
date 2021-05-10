package github

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCreateRepoRequestAsJson(t *testing.T) {
	// Тестируем сериализация и десериализация struct CreateRepoRequest в с JSON с

	// 1) инициализация данных для теста
	request := CreateRepoRequest{
		Name:        "golang introduction",
		Description: "a golang introduction repo",
		Homepage:    "https://github.com",
		Private:     true,
		HasIssues:   true,
		HasProjects: true,
		HasWiki:     true,
	}

	// 2) вызов функции для теста
	// struct в JSON
	bytes, err := json.Marshal(request)
	// JSON в struct
	var target CreateRepoRequest
	err = json.Unmarshal(bytes, &target)

	// 3) валидация
	assert.Nil(t, err)
	assert.NotNil(t, bytes)
	assert.EqualValues(t, `{"name":"golang introduction","description":"a golang introduction repo","homepage":"https://github.com","private":true,"has_issues":true,"has_projects":true,"has_wiki":true}`, string(bytes))

	assert.Nil(t, err)
	assert.EqualValues(t, target.Name, request.Name)
	assert.EqualValues(t, target.Description, request.Description)
	assert.EqualValues(t, target.HasIssues, request.HasIssues)
}
