package pkg

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/anton-uvarenko/headway_test/internal/core"
)

var (
	ErrPerformingRquest = errors.New("error performing request to NASA")
	ErrNasaInternal     = errors.New("error Nasa's internal error")
	ErrInvalidAPIKey    = errors.New("error invalid api key")
	ErrUnknown          = errors.New("error unknown")
	ErrDecodingResponse = errors.New("error decoding response")
)

type NasaAPI struct {
	client HTTPClient
}

func NewNasaAPI(client HTTPClient) *NasaAPI {
	return &NasaAPI{
		client: client,
	}
}

type HTTPClient interface {
	Get(string) (*http.Response, error)
}

func (c NasaAPI) FetchFromNasaApi(date time.Time) (*core.NasaAPIResponse, error) {
	q := url.Values{}
	q.Add("start_date", date.String())
	q.Add("end_date", date.String())
	q.Add("api_key", os.Getenv("API_KEY"))

	resp, err := c.client.Get("https://api.nasa.gov/neo/rest/v1/feed?" + q.Encode())
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrPerformingRquest, err)
	}

	switch {
	case resp.StatusCode == http.StatusOK:
		break
	case resp.StatusCode == http.StatusForbidden:
		return nil, fmt.Errorf("%w: %d", ErrInvalidAPIKey, resp.StatusCode)
	case resp.StatusCode >= 500:
		return nil, fmt.Errorf("%w: %d", ErrNasaInternal, resp.StatusCode)
	default:
		return nil, fmt.Errorf("%w: %d", ErrUnknown, resp.StatusCode)
	}

	var result core.NasaAPIResponse
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return nil, ErrDecodingResponse
	}

	return &result, nil
}
