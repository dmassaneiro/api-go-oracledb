package services

import (
	"errors"
	"net/http"
)

func ValidateToken(token string) error {
	req, err := http.NewRequest(
		"GET",
		"http://localhost:8081/validate",
		nil,
	)
	if err != nil {
		return err
	}

	req.Header.Set("Authorization", "Bearer "+token)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return errors.New("token inválido")
	}

	return nil
}
