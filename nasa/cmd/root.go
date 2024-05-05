package cmd

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/anton-uvarenko/headway_test/internal/core"
	"github.com/anton-uvarenko/headway_test/internal/pkg"
	"github.com/anton-uvarenko/headway_test/internal/usecase"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use: "nasa",
	Run: nasaCommnad,
}

func Execute() error {
	return rootCmd.Execute()
}

func nasaCommnad(cmd *cobra.Command, args []string) {
	caller := pkg.NewNasaAPI(http.DefaultClient)
	dataCh := make(chan *core.NasaAPIResponse)
	errCh := make(chan error)

	go usecase.FetchLastSevenDays(caller, dataCh, errCh)
	result, err := pkg.CollectData(dataCh, errCh)
	if err != nil {
		fmt.Println(err)
		return
	}

	data, _ := json.MarshalIndent(result, "", "    ")
	fmt.Println(string(data))
}
