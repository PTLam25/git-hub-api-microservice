package restclient

// Используем встроенный http client для запроса

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

var (
	enabledMocks = false
	mocks        = make(map[string]*Mock)
)

type Mock struct {
	Url        string
	HttpMethod string
	Response   *http.Response
	Err        error
}

func GetMockId(httpMethod string, url string) string {
	// получить из словаря нужный объект ответа по ключу
	return fmt.Sprintf("%s_%s", httpMethod, url)
}

func StartMockUps() {
	enabledMocks = true
}

func FlushMockUps() {
	// очистить словарь
	mocks = make(map[string]*Mock)
}

func StopMockUps() {
	enabledMocks = false
}

func AddMockUp(mock Mock) {
	mocks[GetMockId(mock.HttpMethod, mock.Url)] = &mock
}

func Post(url string, body interface{}, headers http.Header) (*http.Response, error) {
	if enabledMocks {
		// данные для теста без запроса в GitHub
		mock := mocks[GetMockId(http.MethodPost, url)]
		if mock == nil {
			return nil, errors.New("no mockup found for give request")
		}
		return mock.Response, mock.Err
	}

	// 1) сериализуем body в JSON
	jsonBytes, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}

	// 2) создаем новый объект request
	request, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(jsonBytes))
	request.Header = headers

	// 3) делаем запрос
	client := http.Client{}
	return client.Do(request)
}
