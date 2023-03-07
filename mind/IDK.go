package mind

import (
	api "DUTclock/WorkingWithAPI"
	"io"
	"log"
	"net/http"
	"strconv"
)

// Ping site and return status code
func Ping(domain string) (int, error) {
	var client = http.Client{}

	url := "http://" + domain
	req, err := http.NewRequest("HEAD", url, nil)
	if err != nil {
		return 0, err
	}
	resp, err := client.Do(req)
	if err != nil {
		return 0, err
	}

	err = resp.Body.Close()
	if err != nil {
		return 0, err
	}
	return resp.StatusCode, nil
}

func OwnLTSapi() int {
	count := api.LastApiVersion
	for {
		url := "https://dut-api.lwjerri.ml/v" + strconv.Itoa(count) + "/faculty"
		resp, err := http.Get(url)
		if err != nil {
			log.Fatal(err)
		}
		defer func(Body io.ReadCloser) {
			err := Body.Close()
			if err != nil {
				log.Println(err.Error())
			}
		}(resp.Body)
		if resp.StatusCode == http.StatusOK {
			break
		}
		count++
	}
	return count
}
