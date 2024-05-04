package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/anton-uvarenko/headway_test/internal/core"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()

	q := url.Values{}
	q.Add("start_date", time.Now().Add(-time.Hour*24*7).String())
	q.Add("end_date", time.Now().String())
	q.Add("api_key", os.Getenv("API_KEY"))

	fmt.Println("API_KEY: ", os.Getenv("API_KEY"))

	resp, err := http.DefaultClient.Get("https://api.nasa.gov/neo/rest/v1/feed?" + q.Encode())
	if err != nil {
		fmt.Println("error performing request to NASA: %w", err)
		return
	}

	var result core.NASAResponse
	json.NewDecoder(resp.Body).Decode(&result)
	fmt.Println(result)

	// out, _ := os.Create("result.json")
	// result, _ := io.ReadAll(resp.Body)
	// out.Write(result)
	// fmt.Println(string(result))
}
