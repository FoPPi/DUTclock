package WorkingWithAPI

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
)

func TakeFacultyID(value string) {
	json, _ := TakeFaculty()

	FacultyName = value
	for _, rec := range json {
		if value == rec.Name {
			FacultyID = rec.Id
		}
	}

}

func FacultyJSONtoString() []string {
	json, _ := TakeFaculty()

	sArr := make([]string, len(json))
	counter := 0
	for _, rec := range json {
		sArr[counter] = rec.Name
		counter++
	}

	return sArr

}

// TakeFaculty читает FacultyJSON из url
func TakeFaculty() ([]FacultyJSON, error) {
	url := "https://dutcalendar-tracker.lwjerri.ml/v1/faculty"

	// Get request
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println("No response from request")
		return nil, err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			fmt.Println("Error", err)
		}
	}(resp.Body)
	body, err := ioutil.ReadAll(resp.Body) // response body is []byte

	var result []FacultyJSON
	if err := json.Unmarshal(body, &result); err != nil { // Parse []byte to the go struct pointer
		fmt.Println("Can not unmarshal JSON")
		return nil, err
	}

	return result, err
}

// FacultyJSON структура json-а
type FacultyJSON struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}
