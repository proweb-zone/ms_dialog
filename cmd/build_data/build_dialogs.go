package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

type BuildDataDialog struct {
	url   string
	token string
}

func main() {
	newToken := "dCagXAppjfkLggxhgMfvYGVvbOYLvPiT"

	words := []string{"привет", "как втои дела", "У меня все норм", "что сегодня делал"}
	buildDataDialog := NewBuildDataDialog(newToken)

	for _, msg := range words {
		fmt.Println(msg)
		go buildDataDialog.sendRequest(msg)
		time.Sleep(1 * time.Second)
	}

}

func NewBuildDataDialog(newToken string) *BuildDataDialog {
	newUrl := "http://localhost:3002/v2/dialog/2/send"
	return &BuildDataDialog{url: newUrl, token: newToken}
}

func (b *BuildDataDialog) sendRequest(msg string) error {
	// Подготовка тела запроса
	requestBody := map[string]string{"Msg": msg}
	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		return fmt.Errorf("error marshaling JSON: %v", err)
	}

	// Создание запроса
	req, err := http.NewRequest("POST", b.url, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("error creating request: %v", err)
	}

	// Установка заголовков
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+b.token)

	// Выполнение запроса
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("error sending request: %v", err)
	}
	defer resp.Body.Close()

	// Проверка статуса
	if resp.StatusCode >= 400 {
		return fmt.Errorf("server returned error status: %d", resp.StatusCode)
	}

	// Чтение ответа
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("error reading response: %v", err)
	}

	fmt.Printf("Success! Status: %d\n", resp.StatusCode)
	fmt.Printf("Response: %s\n", string(body))

	return nil
}
