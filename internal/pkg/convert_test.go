package pkg

import (
	"errors"
	"testing"
	"time"

	"github.com/anton-uvarenko/headway_test/internal/core"
	"github.com/magiconair/properties/assert"
)

func TestNasaRespToCliResp(t *testing.T) {
	successDate, _ := time.Parse(time.DateOnly, "2024-04-01")
	testTable := []struct {
		Name     string
		NasaResp *core.NasaAPIResponse
		CliResp  *core.CliResponse
		Error    error
	}{
		{
			Name: "OK",
			NasaResp: &core.NasaAPIResponse{
				Total: 1,
				NearEarthObjects: map[string][]core.NasaNEObject{
					"2024-04-01": {
						{
							Id:                             "123123",
							Name:                           "some_name",
							IsPotentiallyHazardousAsteroid: false,
						},
					},
				},
			},
			CliResp: &core.CliResponse{
				Total: 1,
				NearEarthObjects: []core.CliNEObject{
					{
						Date:                           successDate,
						Id:                             "123123",
						Name:                           "some_name",
						IsPotentiallyHazardousAsteroid: false,
					},
				},
			},
			Error: nil,
		},
		{
			Name: "Date parse error",
			NasaResp: &core.NasaAPIResponse{
				Total: 1,
				NearEarthObjects: map[string][]core.NasaNEObject{
					"01-04-2024": {
						{
							Id:                             "123123",
							Name:                           "some_name",
							IsPotentiallyHazardousAsteroid: false,
						},
					},
				},
			},
			CliResp: nil,
			Error:   ErrParsingTime,
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.Name, func(t *testing.T) {
			result, err := NasaRespToCliResp(testCase.NasaResp)

			assert.Equal(t, errors.Is(err, testCase.Error), true)
			assert.Equal(t, result, testCase.CliResp)
		})
	}
}
