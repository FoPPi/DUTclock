package mind

import (
	api "DUTclock/WorkingWithAPI"
	"encoding/json"
	"fmt"
	"fyne.io/fyne/v2"
	"io"
	"io/ioutil"
	"net/http"
	"strconv"
)

// ReadOfflineJSON читает json файл
func ReadOfflineJSON(jsonName string, sharedPrefs fyne.Preferences) WeekJSON {

	value := sharedPrefs.String(jsonName)

	var Week WeekJSON

	err := json.Unmarshal([]byte(value), &Week)
	if err != nil {
		fmt.Println(err)
	}

	return Week
}

// MyError custom error
type MyError struct{}

func (m *MyError) Error() string {
	return "No internet connection"
}

// UpdateOfflineJSON читает json из url и записывет его в файл
func UpdateOfflineJSON(sharedPrefs fyne.Preferences) (Updated bool, error error) {
	CURRENT, err := TakeWeek(api.FacultyID, api.CourseID, api.GroupID, "CURRENT")
	if err != nil {
		return false, &MyError{}
	}
	file1, _ := json.MarshalIndent(CURRENT, "", " ")

	sharedPrefs.SetString("CURRENT_WeekJSON", string(file1))

	NEXT, err := TakeWeek(api.FacultyID, api.CourseID, api.GroupID, "NEXT")
	if err != nil {
		return false, &MyError{}
	}
	file2, _ := json.MarshalIndent(NEXT, "", " ")

	sharedPrefs.SetString("NEXT_WeekJSON", string(file2))

	return true, nil
}

// TakeWeek читает WeekJSON из url
func TakeWeek(Faculty, Course, Group int, Week string) (*WeekJSON, error) {

	// 1/1/1576/NEXT
	url := "https://dut-api.lwjerri.ml/v" + strconv.Itoa(api.LastApiVersion) + "/student-calendar/" + strconv.Itoa(Faculty) + "/" + strconv.Itoa(Course) + "/" + strconv.Itoa(Group) + "/" + Week

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

	var result *WeekJSON
	if err := json.Unmarshal(body, &result); err != nil { // Parse []byte to the go struct pointer
		fmt.Println("Can not unmarshal JSON")
		return nil, err
	}

	return result, err
}

// WeekJSON структура json-а
type WeekJSON struct {
	IsCachedResponse bool   `json:"isCachedResponse"`
	IsDataFromDB     bool   `json:"isDataFromDB"`
	DataHash         string `json:"dataHash"`
	Data             []struct {
		AddedAt         string `json:"addedAt"`
		Cabinet         string `json:"cabinet"`
		DayNameLong     string `json:"dayNameLong"`
		DayNameShort    string `json:"dayNameShort"`
		EndAt           string `json:"endAt"`
		GroupName       string `json:"groupName"`
		LectorFullName  string `json:"lectorFullName"`
		LectorShortName string `json:"lectorShortName"`
		LessonDate      string `json:"lessonDate"`
		LessonLongName  string `json:"lessonLongName"`
		LessonNumber    int    `json:"lessonNumber"`
		LessonShortName string `json:"lessonShortName"`
		LessonType      string `json:"lessonType"`
		StartAt         string `json:"startAt"`
	} `json:"data"`
}
