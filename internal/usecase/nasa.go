package usecase

import (
	"time"

	"github.com/anton-uvarenko/headway_test/internal/core"
	"github.com/anton-uvarenko/headway_test/internal/pkg"
	"golang.org/x/sync/errgroup"
)

const Day = time.Hour * 24

func FetchLastSevenDays(caller *pkg.NasaAPI, dataCh chan<- *core.NasaAPIResponse, errCh chan<- error) {
	defer close(dataCh)
	defer close(errCh)

	eg := &errgroup.Group{}

	for i := range 7 {
		eg.Go(func() error {
			date := time.Now().Add(-Day * time.Duration(i+1))

			data, err := caller.FetchFromNasaApi(date)
			if err != nil {
				return err
			}

			dataCh <- data
			return nil
		})
	}

	err := eg.Wait()
	if err != nil {
		errCh <- err
	}
}
