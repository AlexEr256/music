package utils

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/AlexEr256/musicService/dto"
)

func ValidateEnvParams(user, password, host, port, database string) []string {
	validationErrors := make([]string, 0)

	if user == "" {
		validationErrors = append(validationErrors, "PG_USER value is empty.")
	}
	if password == "" {
		validationErrors = append(validationErrors, "PG_PASSWORD value is empty.")
	}
	if host == "" {
		validationErrors = append(validationErrors, "PG_HOST value is empty.")
	}
	if database == "" {
		validationErrors = append(validationErrors, "PG_DB value is empty.")
	}

	if port == "" {
		validationErrors = append(validationErrors, "PG_PORT value is empty.")
	}
	_, err := strconv.Atoi(port)
	if err != nil {
		validationErrors = append(validationErrors, "PG_PORT must be valid number.")
	}

	return validationErrors
}

func DumbSearchHook(group, song string) (*dto.SongExtraInfoResponse, error) {
	msg := &dto.SongExtraInfoResponse{
		Text:        "Hello,\n my dear friend!\n It is nice\n to see you!",
		Link:        "youtube",
		ReleaseDate: "22.02.2024",
	}

	return msg, nil
}

func SearchHook(group, song string) (*dto.SongExtraInfoResponse, error) {
	host := os.Getenv("INFO_HOOK")
	if host == "" {
		return nil, fmt.Errorf("api hook host is empty. Check input parameters")
	}

	errCh := make(chan error)
	defer close(errCh)

	respChan := make(chan *dto.SongExtraInfoResponse)
	defer close(respChan)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*4)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, host, nil)

	q := req.URL.Query()
	q.Add("group", group)
	q.Add("song", song)
	req.URL.RawQuery = q.Encode()

	if err != nil {
		return nil, fmt.Errorf("failed to create request with timeout - %w", err)
	}

	go func() {
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			errCh <- err
		}
		if resp != nil {
			defer resp.Body.Close()
			var msg dto.SongExtraInfoResponse
			err = json.NewDecoder(resp.Body).Decode(&msg)
			if err != nil {
				errCh <- err
			}
			respChan <- &msg
		}
	}()

	for {
		select {
		case <-ctx.Done():
			return nil, fmt.Errorf("request timeout was exceeded")
		case hookErr := <-errCh:
			return nil, hookErr
		case res := <-respChan:
			return res, nil
		}
	}
}
