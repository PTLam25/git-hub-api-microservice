package repositories

// Структура Request и Response, который сервер будет принимать от пользователя и возвращать ему

type CreateRepoRequest struct {
	Name string `json:"name"`
}

type CreateResponse struct {
	Id    int64  `json:"id"`
	Owner string `json:"owner"`
	Name  string `json:"name"`
}