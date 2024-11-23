package test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"sync"
)

func RequestOnPath(path string, method string, data map[string]interface{}) string {
	var jsonData []byte
	var err error

	if data != nil {
		jsonData, err = json.Marshal(data)
		if err != nil {
			return "fail: " + err.Error()
		}
	}

	req, err := http.NewRequest(method, path, bytes.NewBuffer(jsonData))
	if err != nil {
		return "fail: " + err.Error()
	}
	req.Header.Set("Content-Type", "application/json") // Устанавливаем заголовок для JSON

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "fail: " + err.Error()
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "fail: " + err.Error()
	}

	cleanedBody := strings.ReplaceAll(string(body), "\"", "")
	return strings.TrimSpace(cleanedBody)
}

func StartTest() {
	createWalletPath := "http://localhost:8080/api/v1/wallets/"
	updateWalletPath := "http://localhost:8080/api/v1/wallet"
	getWalletPath := "http://localhost:8080/api/v1/wallets/"

	balanceMap := map[string]interface{}{
		"balance": 10000,
	}

	wg := sync.WaitGroup{}
	idList := make([]string, 0, 100)
	resultAfterUpdate := make([]string, 0, 100)

	// Мьютекс для защиты доступа к idList и resultAfterUpdate
	var mu sync.Mutex

	// Создание кошельков
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			id := RequestOnPath(createWalletPath, "POST", balanceMap)
			mu.Lock()
			idList = append(idList, id)
			mu.Unlock()
		}()
	}
	wg.Wait()

	// Обновление кошельков
	for i := 0; i < 2; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			mu.Lock()
			id := idList[1]
			mu.Unlock()

			updateBalanceMap := map[string]interface{}{
				"id":            id,
				"typeOperation": "WITHDRAW",
				"amount":        100,
			}
			RequestOnPath(updateWalletPath, "PUT", updateBalanceMap)
		}()
	}
	wg.Wait()

	// Получение обновленных данных
	for i := 0; i < 1; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			mu.Lock()
			if i >= len(idList) {
				fmt.Printf("Error: index out of range for ID list: %d\n", i)
				mu.Unlock()
				return
			}
			id := idList[1]
			mu.Unlock()

			newPath := getWalletPath + id
			response := RequestOnPath(newPath, "GET", nil)
			mu.Lock()
			resultAfterUpdate = append(resultAfterUpdate, response)
			mu.Unlock()
		}(i)
	}
	wg.Wait()

	fmt.Println("Results after update:", resultAfterUpdate)
}
