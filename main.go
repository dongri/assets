package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"time"
)

var httpClient = &http.Client{Timeout: 10 * time.Second}

func getJSON(url string, target interface{}) error {
	r, err := httpClient.Get(url)
	if err != nil {
		return err
	}
	defer r.Body.Close()
	return json.NewDecoder(r.Body).Decode(target)
}

type res struct {
	Rate string `json:"rate"`
}

// CryptoCurrency ...
type CryptoCurrency []struct {
	Ticker string  `json:"ticker"`
	Hold   float64 `json:"hold"`
}

func main() {
	api := "https://coincheck.com/api/rate/%s_jpy"
	raw, err := ioutil.ReadFile("./assets.json")
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	var cc CryptoCurrency
	json.Unmarshal(raw, &cc)
	assets := float64(0)
	for _, c := range cc {
		r := new(res)
		getJSON(fmt.Sprintf(api, c.Ticker), r)
		rate, _ := strconv.ParseFloat(r.Rate, 64)
		hold := c.Hold
		fmt.Printf("%s: %f * %f\n", c.Ticker, hold, rate)
		assets += hold * rate
	}
	fmt.Printf("Sum: %f\n", assets)
}
