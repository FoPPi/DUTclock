package WorkingWithAPI

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strconv"
)

func TakeGroupID(value string) {
	GroupJson, _ := TakeGroup(FacultyID, CourseID)

	GroupName = value
	for _, rec := range GroupJson.Data {
		if value == rec.Name {
			GroupID = rec.Id
		}
	}
}

func GroupJSONtoString(value string) []string {

	TakeCourseID(value)

	GroupJson, _ := TakeGroup(FacultyID, CourseID)

	sArr := make([]string, len(GroupJson.Data))
	counter := 0
	for _, rec := range GroupJson.Data {
		sArr[counter] = rec.Name
		counter++
	}

	return sArr

}

// TakeGroup читает GroupJSON из url
func TakeGroup(FacultyID, CourseID int) (*GroupJSON, error) {
	url := "https://dut-api.lwjerri.ml/v4/group/" + strconv.Itoa(FacultyID) + "/" + strconv.Itoa(CourseID)

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

	var result *GroupJSON
	if err := json.Unmarshal(body, &result); err != nil { // Parse []byte to the go struct pointer
		fmt.Println("Can not unmarshal JSON")
		return nil, err
	}

	return result, err
}

// GroupJSON структура json-а
type GroupJSON struct {
	IsCachedResponse bool   `json:"isCachedResponse"`
	IsDataFromDB     bool   `json:"isDataFromDB"`
	DataHash         string `json:"dataHash"`
	Data             []struct {
		Name string `json:"name"`
		Id   int    `json:"id"`
	} `json:"data"`
}
