package createRequests

import (
	"encoding/json"
	"fmt"
)

type GraphQLRequest struct {
	Query     string         `json:"query"`
	Variables LoginVariables `json:"variables"`
}

type LoginVariables struct {
	Login       string `json:"login"`
	Password    string `json:"password"`
	Fingerprint string `json:"fingerprint"`
	ForceLogin  bool   `json:"forceLogin"`
}

func CreateRequestAuthorization() (string, error) {
	var login string
	var password string
	query := `mutation ($login: String!, $password: String!, $fingerprint: String!, $forceLogin: Boolean!) {
		authentication {
			security {
				login(
					username: $login
					password: $password
					fingerprint: $fingerprint
					forceLogin: $forceLogin
				) {
					accessToken
				}
			}
		}
	}`

	fmt.Print("Введите логин: ")
	_, err := fmt.Scanln(&login)
	fmt.Print("Введите пароль: ")
	_, err = fmt.Scanln(&password)
	if err != nil {
		return "", fmt.Errorf("CreateRequestAuthorization: ошибка ввода логина или пароля: \n%w", err)
	}

	variables := LoginVariables{
		Login:      login,
		Password:   password,
		ForceLogin: true,
	}

	request := GraphQLRequest{
		Query:     query,
		Variables: variables,
	}

	jsonRequest, err := json.Marshal(request)
	if err != nil {
		return "", fmt.Errorf("CreateRequestAuthorization: ошибка сериализации: \n%w", err)
	}

	return string(jsonRequest), nil
}
