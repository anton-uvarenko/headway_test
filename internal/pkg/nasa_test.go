package pkg

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"testing"
	"time"

	"github.com/anton-uvarenko/headway_test/internal/core"
	"github.com/magiconair/properties/assert"
)

type HTTPClientMock struct {
	StatusCodeToReturn int
	ErrorOnRequest     bool
	ErrorOnDecode      bool
}

func (c HTTPClientMock) Get(url string) (*http.Response, error) {
	if c.ErrorOnRequest {
		return nil, errors.New("failing on request")
	}

	response := &http.Response{
		StatusCode: c.StatusCodeToReturn,
	}

	if c.ErrorOnDecode {
		response.Body = io.NopCloser(bytes.NewBuffer([]byte("asdfa")))
		return response, nil
	}

	if c.StatusCodeToReturn == http.StatusOK {
		result := core.NasaAPIResponse{
			Total: 1,
		}
		data, _ := json.Marshal(result)
		response.Body = io.NopCloser(bytes.NewBuffer(data))
		response.ContentLength = int64(len(data))
	}

	return response, nil
}

func TestFetchFromNasaApi(t *testing.T) {
	testTable := []struct {
		Name               string
		ExpectedStatusCode int
		ExpectedOutput     *core.NasaAPIResponse
		ExpectedError      error
	}{
		{
			Name:               "OK",
			ExpectedStatusCode: http.StatusOK,
			ExpectedOutput: &core.NasaAPIResponse{
				Total: 1,
			},
			ExpectedError: nil,
		},
		{
			Name:               "Server error",
			ExpectedStatusCode: http.StatusInternalServerError,
			ExpectedOutput:     nil,
			ExpectedError:      ErrNasaInternal,
		},
		{
			Name:               "Forbiden",
			ExpectedStatusCode: http.StatusForbidden,
			ExpectedOutput:     nil,
			ExpectedError:      ErrInvalidAPIKey,
		},
		{
			Name:               "Unknown error",
			ExpectedStatusCode: http.StatusTeapot,
			ExpectedOutput:     nil,
			ExpectedError:      ErrUnknown,
		},
		{
			Name:               "Request error",
			ExpectedStatusCode: 0,
			ExpectedOutput:     nil,
			ExpectedError:      ErrPerformingRquest,
		},
		{
			Name:               "Response decode error",
			ExpectedStatusCode: http.StatusOK,
			ExpectedOutput:     nil,
			ExpectedError:      ErrDecodingResponse,
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.Name, func(t *testing.T) {
			client := &HTTPClientMock{}
			nasaApi := NewNasaAPI(client)

			client.StatusCodeToReturn = testCase.ExpectedStatusCode

			if errors.Is(testCase.ExpectedError, ErrPerformingRquest) {
				client.ErrorOnRequest = true
			}
			if errors.Is(testCase.ExpectedError, ErrDecodingResponse) {
				client.ErrorOnDecode = true
			}

			result, err := nasaApi.FetchFromNasaApi(time.Now())

			assert.Equal(t, result, testCase.ExpectedOutput)
			assert.Equal(t, errors.Is(err, testCase.ExpectedError), true)

			client.ErrorOnRequest = false
			client.ErrorOnDecode = false
		})
	}
}
