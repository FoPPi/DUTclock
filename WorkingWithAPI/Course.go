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
	CourseJson, _ := TakeCourse(FacultyID)

	CourseName = value
	for _, rec := range CourseJson {
		if value == rec.Name {
			CourseID = rec.Id
		}
	}
}

func CourseJSONtoString(value string) []string {

	TakeFacultyID(value)

	CourseJson, _ := TakeCourse(FacultyID)

	sArr := make([]string, len(CourseJson))
	counter := 0
	for _, rec := range CourseJson {
		sArr[counter] = rec.Name
		counter++
	}

	return sArr

}

// TakeCourse читает CourseJSON из url
func TakeCourse(FacultyID int) ([]CourseJSON, error) {
	url := "https://dutcalendar-tracker.lwjerri.ml/v1/course/" + strconv.Itoa(FacultyID)

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

	var result []CourseJSON
	if err := json.Unmarshal(body, &result); err != nil { // Parse []byte to the go struct pointer
		fmt.Println("Can not unmarshal JSON")
		return nil, err
	}

	return result, err
}

// CourseJSON структура json-а
type CourseJSON struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}
