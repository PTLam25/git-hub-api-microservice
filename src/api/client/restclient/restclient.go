package restclient

// Используем встроенный http client для запроса

import (
	"bytes"
	"encoding/json"
	"net/http"
)

func Post(url string, body interface{}, headers http.Header) (*http.Response, error) {
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
