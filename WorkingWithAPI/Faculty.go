package WorkingWithAPI

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
)

func TakeFacultyID(value string) {
	faculty, _ := TakeFaculty()

	FacultyName = value
	for _, rec := range faculty {
		if value == rec.Name {
			FacultyID = rec.Id
		}
	}

}

func FacultyJSONtoString() []string {
	faculty, _ := TakeFaculty()

	sArr := make([]string, len(faculty))
	counter := 0
	for _, rec := range faculty {
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
