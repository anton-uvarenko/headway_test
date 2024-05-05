package pkg

import (
	"errors"
	"testing"
	"time"

	"github.com/anton-uvarenko/headway_test/internal/core"
	"github.com/magiconair/properties/assert"
)

func TestCollectData(t *testing.T) {
	today, _ := time.Parse(time.DateOnly, time.Now().Format(time.DateOnly))
	yesterday, _ := time.Parse(time.DateOnly, time.Now().Add(-time.Hour*24).Format(time.DateOnly))
	testTable := []struct {
		Name           string
		Input          []*core.NasaAPIResponse
		ExpectedOutput *core.CliResponse
		ExpectedError  error
	}{
		{
			Name: "Ok",
			Input: []*core.NasaAPIResponse{
				{
					Total: 1,
					NearEarthObjects: map[string][]core.NasaNEObject{
						time.Now().Format(time.DateOnly): {
							{
								Id:                             "123",
								Name:                           "123",
								IsPotentiallyHazardousAsteroid: false,
							},
						},
					},
				},
				{
					Total: 1,
					NearEarthObjects: map[string][]core.NasaNEObject{
						time.Now().Add(-time.Hour * 24).Format(time.DateOnly): {
							{
								Id:                             "122",
								Name:                           "122",
								IsPotentiallyHazardousAsteroid: true,
							},
						},
					},
				},
			},
			ExpectedOutput: &core.CliResponse{
				Total: 2,
				NearEarthObjects: []core.CliNEObject{
					{
						Date:                           today,
						Id:                             "123",
						Name:                           "123",
						IsPotentiallyHazardousAsteroid: false,
					},
					{
						Date:                           yesterday,
						Id:                             "122",
						Name:                           "122",
						IsPotentiallyHazardousAsteroid: true,
					},
				},
			},
			ExpectedError: nil,
		},
		{
			Name:           "Error",
			Input:          []*core.NasaAPIResponse{},
			ExpectedOutput: nil,
			ExpectedError:  ErrUnknown,
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.Name, func(t *testing.T) {
			dataCh := make(chan *core.NasaAPIResponse)
			errCh := make(chan error)

			go func() {
				defer close(dataCh)
				defer close(errCh)

				for _, v := range testCase.Input {
					dataCh <- v
				}

				if testCase.ExpectedError != nil {
					errCh <- testCase.ExpectedError
				}
			}()

			result, err := CollectData(dataCh, errCh)

			assert.Equal(t, result, testCase.ExpectedOutput)
			assert.Equal(t, errors.Is(err, testCase.ExpectedError), true)
		})
	}
}
