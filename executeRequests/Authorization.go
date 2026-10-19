package executeRequests

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	. "github.com/VladislavTroyan163/POST-request-generation-script/config"
)

type AuthorizationResponse struct {
	Data   AuthorizationData `json:"data"`
	Errors []GraphQLError    `json:"errors"`
}

type AuthorizationData struct {
	Authentication Authentication `json:"authentication"`
}

type Authentication struct {
	Security Security `json:"security"`
}

type Security struct {
	Login LoginResult `json:"login"`
}

type LoginResult struct {
	AccessToken string `json:"accessToken"`
}

type GraphQLError struct {
	Message    string `json:"message"`
	Extensions struct {
		MessageID      string `json:"messageId"`
		Classification string `json:"classification"`
	} `json:"extensions"`
}

func ExecuteRequestAuthorization(request string, config Config) (resultRequest string, err error) {
	body := strings.NewReader(request)
	req, err := http.NewRequest(
		http.MethodPost,
		config.ServerURLAuthorization,
		body,
	)
	if err != nil {
		return "", fmt.Errorf("ExecuteRequestAuthorization: ошибка: \n%w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("ExecuteRequestAuthorization: ошибка: \n%w", err)
	}
	defer resp.Body.Close()

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("ExecuteRequestAuthorization: ошибка: \n%w", err)
	}
	var response AuthorizationResponse

	err = json.Unmarshal(responseBody, &response)
	if err != nil {
		return "", fmt.Errorf("ExecuteRequestAuthorization: ошибка: \n%w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf(
			"ExecuteRequestAuthorization: сервер вернул HTTP-статус %s",
			resp.Status,
		)
	}

	if len(response.Errors) > 0 {
		return "", fmt.Errorf(
			"ExecuteRequestAuthorization: ошибка авторизации: %s",
			response.Errors[0].Message,
		)
	}

	token := response.Data.Authentication.Security.Login.AccessToken
	if token == "" {
		return "", fmt.Errorf(
			"ExecuteRequestAuthorization: сервер не вернул токен авторизации",
		)
	}

	return token, nil
}
