package WorkingWithAPI

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strconv"
)

func TakeCourseID(value string) {
	CourseJson, _ := TakeCourse(App.Preferences().Int("FacultyID"))

	sharedPrefs := App.Preferences()
	sharedPrefs.SetString("CourseName", value)
	for _, rec := range CourseJson.Data {
		if value == rec.Name {
			sharedPrefs.SetInt("CourseID", rec.Id)
			break
		}
	}
}

func CourseJSONtoString(value string) []string {

	TakeFacultyID(value)

	CourseJson, _ := TakeCourse(App.Preferences().Int("FacultyID"))

	sArr := make([]string, len(CourseJson.Data))
	counter := 0
	for _, rec := range CourseJson.Data {
		sArr[counter] = rec.Name
		counter++
	}

	return sArr

}

// TakeCourse читает CourseJSON из url
func TakeCourse(FacultyID int) (*CourseJSON, error) {
	url := ApiURL + "/v" + strconv.Itoa(App.Preferences().Int("LastApiVersion")) + "/course/" + strconv.Itoa(FacultyID)

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

	var result *CourseJSON
	if err := json.Unmarshal(body, &result); err != nil { // Parse []byte to the go struct pointer
		fmt.Println("Can not unmarshal JSON")
		return nil, err
	}

	return result, err
}

// CourseJSON структура json-а
type CourseJSON struct {
	IsCachedResponse bool   `json:"isCachedResponse"`
	IsDataFromDB     bool   `json:"isDataFromDB"`
	DataHash         string `json:"dataHash"`
	Data             []struct {
		Name string `json:"name"`
		Id   int    `json:"id"`
	} `json:"data"`
}
