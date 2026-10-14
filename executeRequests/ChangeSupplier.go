package executeRequests

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	. "github.com/VladislavTroyan163/POST-request-generation-script/config"
)

type Response struct {
	Tasks []ResponseTask `json:"tasks"`
}

type ResponseTask struct {
	Result ResponseResult `json:"result"`
}

type ResponseResult struct {
	НомерЗаписи int `json:"номерЗаписи"`
}

func ExecuteRequestChangeSupplier(request string, лицевойСчет int) (resultRequest string, err error) {
	config, err := ReadConfig()
	if err != nil {
		return "", fmt.Errorf("ExecuteRequestChangeSupplier: ошибка: \n%w", err)
	}

	body := strings.NewReader(request)
	req, err := http.NewRequest(
		http.MethodPost,
		config.ServerURL,
		body,
	)
	if err != nil {
		return "", fmt.Errorf("ExecuteRequestChangeSupplier: ошибка: \n%w", err)
	}

	req.Header.Set("s-access-token", config.Token)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("ExecuteRequestChangeSupplier: ошибка: \n%w", err)
	}
	defer resp.Body.Close()

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("ExecuteRequestChangeSupplier: ошибка: \n%w", err)
	}

	var response Response

	switch resp.StatusCode {
	case 200:
		err = json.Unmarshal(responseBody, &response)
		if err != nil {
			return "", fmt.Errorf("ExecuteRequestChangeSupplier: ошибка: \n%w", err)
		}
		resultRequest = fmt.Sprintln(
			"\nHTTP-запрос завершился успешно:", resp.Status,
			"\nНомер записи в таблице \"Поставщики\":", response.Tasks[0].Result.НомерЗаписи,
			"\nrow_id лицевого счета:", лицевойСчет)
	case 401:
		resultRequest = fmt.Sprintln(
			"Ошибка авторизации:", resp.Status,
			"\nrow_id лицевого счета:", лицевойСчет)
	default:
		resultRequest = fmt.Sprintln(
			"HTTP-запрос завершился неуспешно:", resp.Status,
			"\nrow_id лицевого счета:", лицевойСчет)
	}
	return resultRequest, nil
}
