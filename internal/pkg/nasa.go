package pkg

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"sync"
	"time"

	"github.com/anton-uvarenko/headway_test/internal/core"
)

var (
	ErrPerformingRquest = errors.New("error performing request to NASA")
	ErrNasaInternal     = errors.New("error Nasa's internal error")
	ErrBadRequest       = errors.New("error ")
)

func FetchFromNasaApi(date time.Time) (*core.NasaAPIResponse, error) {
	q := url.Values{}
	q.Add("start_date", date.String())
	q.Add("end_date", date.String())
	q.Add("api_key", os.Getenv("API_KEY"))

	resp, err := http.DefaultClient.Get("https://api.nasa.gov/neo/rest/v1/feed?" + q.Encode())
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrPerformingRquest, err)
	}

	if resp.StatusCode >= 500 {
		return nil, fmt.Errorf("%w: %d", ErrNasaInternal, resp.StatusCode)
	}

	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("%w: %d", ErrBadRequest, resp.StatusCode)
	}

	var result core.NasaAPIResponse
	json.NewDecoder(resp.Body).Decode(&result)

	return &result, nil
}

func FetchLastSevenDays() error {
	wg := &sync.WaitGroup{}
	mx := sync.Mutex{}
	wg.Add(7)
	for i := range 7 {
		go func(wg *sync.WaitGroup) {
			defer wg.Done()
			date := time.Now().Add(-time.Hour * 24 * time.Duration(i+1))
			data, err := FetchFromNasaApi(date)
			if err != nil {
				fmt.Println(err)
				return
			}

			mx.Lock()

			cliResp, err := NasaRespToCliResp(data)
			if err != nil {
				fmt.Println(err)
				return
			}
			result, err := json.MarshalIndent(cliResp, "", "    ")
			if err != nil {
				fmt.Println(err)
				return
			}
			fmt.Println(string(result))

			mx.Unlock()
		}(wg)
	}
	wg.Wait()

	return nil
}
