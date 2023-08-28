package WorkingWithAPI

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strconv"
)

func TakeFacultyID(value string) {
	faculty, _ := TakeFaculty()

	sharedPrefs := App.Preferences()
	sharedPrefs.SetString("FacultyName", value)
	for _, rec := range faculty.Data {
		if value == rec.Name {
			sharedPrefs.SetInt("FacultyID", rec.Id)
			break
		}
	}

}

func FacultyJSONtoString() []string {
	faculty, _ := TakeFaculty()

	sArr := make([]string, len(faculty.Data))
	counter := 0
	for _, rec := range faculty.Data {
		sArr[counter] = rec.Name
		counter++
	}

	return sArr

}

// TakeFaculty читает FacultyJSON из url
func TakeFaculty() (*FacultyJSON, error) {
	url := ApiURL + "/v" + strconv.Itoa(App.Preferences().Int("LastApiVersion")) + "/faculty"

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

	var result *FacultyJSON
	if err := json.Unmarshal(body, &result); err != nil { // Parse []byte to the go struct pointer
		fmt.Println("Can not unmarshal JSON")
		return nil, err
	}

	return result, err
}

// FacultyJSON структура json-а
type FacultyJSON struct {
	IsCachedResponse bool   `json:"isCachedResponse"`
	IsDataFromDB     bool   `json:"isDataFromDB"`
	DataHash         string `json:"dataHash"`
	Data             []struct {
		Name string `json:"name"`
		Id   int    `json:"id"`
	} `json:"data"`
}
