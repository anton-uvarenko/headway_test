package usecase

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"testing"

	"github.com/anton-uvarenko/headway_test/internal/core"
	"github.com/anton-uvarenko/headway_test/internal/pkg"
	"github.com/magiconair/properties/assert"
)

type HTTPClientMock struct {
	ToFail bool
}

func (c HTTPClientMock) Get(url string) (*http.Response, error) {
	if c.ToFail {
		return nil, errors.New("fail")
	}
	response := &http.Response{
		StatusCode: http.StatusOK,
	}

	result := core.NasaAPIResponse{
		Total: 1,
	}
	data, _ := json.Marshal(result)
	response.Body = io.NopCloser(bytes.NewBuffer(data))
	response.ContentLength = int64(len(data))

	return response, nil
}

func TestFetchLastSevenDays(t *testing.T) {
	testTable := []struct {
		Name                 string
		ExpectedOutputLength int
		ExpectedError        error
	}{
		{
			Name:                 "OK",
			ExpectedOutputLength: 7,
			ExpectedError:        nil,
		},
		{
			Name:                 "Fail",
			ExpectedOutputLength: 0,
			ExpectedError:        pkg.ErrPerformingRquest,
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.Name, func(t *testing.T) {
			client := &HTTPClientMock{}
			nasaApi := pkg.NewNasaAPI(client)
			var err error

			dataCh := make(chan *core.NasaAPIResponse)
			errCh := make(chan error)

			if testCase.ExpectedError != nil {
				client.ToFail = true
			}
			go FetchLastSevenDays(nasaApi, dataCh, errCh)

			total := 0
		outer:
			for {
				select {
				case _, ok := <-dataCh:
					if !ok {
						break outer
					}

					total++
				case e := <-errCh:
					err = e
					break outer
				}
			}

			assert.Equal(t, total, testCase.ExpectedOutputLength)
			assert.Equal(t, errors.Is(err, testCase.ExpectedError), true)
		})
	}
}
