package pkg

import "github.com/anton-uvarenko/headway_test/internal/core"

func CollectData(dataCh <-chan *core.NasaAPIResponse, errCh <-chan error) (*core.CliResponse, error) {
	result := &core.CliResponse{}
	for {
		select {
		case resp, ok := <-dataCh:
			if !ok {
				return result, nil
			}
			cliResp, err := NasaRespToCliResp(resp)
			if err != nil {
				return nil, err
			}

			result.Total += cliResp.Total
			result.NearEarthObjects = append(result.NearEarthObjects, cliResp.NearEarthObjects...)

		case err := <-errCh:
			return nil, err
		}
	}
}
