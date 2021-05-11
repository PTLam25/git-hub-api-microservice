package config

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestConstants(t *testing.T) {
	// тестируем константы в файле. Это важно проверять,
	//так как при изменения константы программа может поменять свое поведение
	assert.EqualValues(t, "apiGithubAccessToken", githubAccessToken)
}

func TestGetGithubAccessToken(t *testing.T) {
	// должна быть пустая строка, так как при запуске теста нет никаких env переменных
	assert.EqualValues(t, "", GetGithubAccessToken())
}
