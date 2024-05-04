package pkg

import (
	"errors"
	"fmt"
	"time"

	"github.com/anton-uvarenko/headway_test/internal/core"
)

var ErrParsingTime = errors.New("pkg.NasaRespNoCliResp error parsing time")

func NasaRespToCliResp(resp *core.NASAResponse) (*core.CliResponse, error) {
	result := &core.CliResponse{
		Total: resp.Total,
	}

	for date, objects := range resp.NearEarthObjects {
		for _, object := range objects {
			d, err := time.Parse(time.DateOnly, date)
			if err != nil {
				return nil, fmt.Errorf("%w: [%w]", ErrParsingTime, err)
			}
			result.NearEarthObjects = append(result.NearEarthObjects, core.CliNEObject{
				Date:                           d,
				Id:                             object.Id,
				Name:                           object.Name,
				IsPotentiallyHazardousAsteroid: object.IsPotentiallyHazardousAsteroid,
			})
		}
	}

	return result, nil
}
