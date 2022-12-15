package main

import (
	api "DUTclock/WorkingWithAPI"
	"encoding/json"
	"fmt"
	"fyne.io/fyne/v2"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

// TakeTime показывет сколько до начала/конца пары (надо оптимизировать)
func TakeTime(app fyne.App) (paraExist bool, paraName string, diff time.Duration) {

	var (
		dateNow       = time.Now()
		result        WeekJSON
		count         = 0
		paruNaSegodna = 0
		dateNowParsed time.Time
		paraType      string = ""
	)

	if !isWeekend(dateNow) {
		h, min, _ := dateNow.Clock()
		y, mon, d := dateNow.Date()
		dateString := strconv.Itoa(h) + ":" + strconv.Itoa(min) + " " + strconv.Itoa(d) + "." + strconv.Itoa(int(mon)) + "." + strconv.Itoa(y)
		if h < 10 {
			var err error
			dateNowParsed, err = time.Parse("3:4 2.1.2006", dateString)
			if err != nil {
				fmt.Println(err)
			}
		} else {
			var err error
			dateNowParsed, err = time.Parse("15:4 2.1.2006", dateString)
			if err != nil {
				fmt.Println(err)
			}
		}

		isSecondJSON := false

		result = ReadOfflineJSON("files/CURRENT_WeekJSON.json")
		for i := 1; i <= 2; i++ {

			// "1", "1", "1576", "CURRENT"
			for _, rec := range result.Data {

				StartTime, err := time.Parse("15:04 02.01.2006", rec.StartAt+" "+rec.LessonDate)

				FinishTime, err := time.Parse("15:04 02.01.2006", rec.EndAt+" "+rec.LessonDate)

				if err != nil {
					fmt.Println(err)
					break
				}

				if dateNowParsed.Day() == StartTime.Day() &&
					dateNowParsed.Month() == StartTime.Month() &&
					dateNowParsed.Year() == StartTime.Year() {

					paruNaSegodna++

					if api.LessonName {
						paraName = rec.LessonLongName
					} else {
						paraName = rec.LessonShortName
					}

					if api.LessonType {
						paraType = rec.LessonType
					}

					if StartTime.Equal(dateNowParsed) {
						app.SendNotification(fyne.NewNotification("Пара почалася", PrettyPrint(paraName+paraType)))
					} else if FinishTime.Equal(dateNowParsed) {
						app.SendNotification(fyne.NewNotification("Пара закинчилася", "Ливай нахуй"))
					}

					if dateNowParsed.Before(StartTime) {
						diff = StartTime.Sub(dateNowParsed)

						return true, "До початку: " + PrettyPrint(paraName+paraType), diff
					} else if dateNowParsed.Before(FinishTime) {
						diff = FinishTime.Sub(dateNowParsed)

						return true, "До кінця: " + PrettyPrint(paraName+paraType), diff
					} else {
						count++
					}
				}
			}
			if paruNaSegodna == count {
				if isSecondJSON {
					return false, "Пари закінчилися :)", diff
				} else {
					isSecondJSON = true
					result = ReadOfflineJSON("files/NEXT_WeekJSON.json")
					count = 0
					paruNaSegodna = 0
				}
			}
		}
	} else {
		//Вывест что пар нету
		//fmt.Println("Сьогодні немає пар :)")
		return false, "Сьогодні немає пар :)", diff
	}

	return false, "Дикая ошибка", diff
}

// ReadOfflineJSON читает json файл
func ReadOfflineJSON(jsonName string) WeekJSON {
	jsonFile, err := os.Open(jsonName)
	// if we os.Open returns an error then handle it
	if err != nil {
		fmt.Println(err)
	}

	var Week WeekJSON

	byteValue, _ := ioutil.ReadAll(jsonFile)

	err = json.Unmarshal(byteValue, &Week)
	if err != nil {
		fmt.Println(err)
	}

	err = jsonFile.Close()
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
func UpdateOfflineJSON() (Updated bool, error error) {
	CURRENT, err := TakeWeek(api.FacultyID, api.CourseID, api.GroupID, "CURRENT")
	if err != nil {
		return false, &MyError{}
	}
	file1, _ := json.MarshalIndent(CURRENT, "", " ")

	var nestedDir = "files"
	path := filepath.Join(".", nestedDir)
	erro := os.MkdirAll(path, 0777)
	if erro != nil {
		log.Fatal(erro)
	}

	_ = ioutil.WriteFile("files/CURRENT_WeekJSON.json", file1, 0644)

	NEXT, err := TakeWeek(api.FacultyID, api.CourseID, api.GroupID, "NEXT")
	if err != nil {
		return false, &MyError{}
	}
	file2, _ := json.MarshalIndent(NEXT, "", " ")

	_ = ioutil.WriteFile("files/NEXT_WeekJSON.json", file2, 0644)

	return true, nil
}

// TakeWeek читает WeekJSON из url
func TakeWeek(Faculty, Course, Group int, Week string) (*WeekJSON, error) {

	// 1/1/1576/NEXT
	url := "https://dut-api.lwjerri.ml/v3/calendar/" + strconv.Itoa(Faculty) + "/" + strconv.Itoa(Course) + "/" + strconv.Itoa(Group) + "/" + Week

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

// isWeekend проверка на выходной
func isWeekend(t time.Time) bool {
	t = t.UTC()
	switch t.Weekday() {
	case time.Friday:
		h, _, _ := t.Clock()
		if h >= 12+10 {
			return true
		}
	case time.Saturday:
		return true
	case time.Sunday:
		h, m, _ := t.Clock()
		if h < 12+10 {
			return true
		}
		if h == 12+10 && m <= 5 {
			return true
		}
	}
	return false
}

// PrettyPrint вывод текста без ковычек ""
func PrettyPrint(text any) string {
	s, _ := json.Marshal(text)
	res := strings.ReplaceAll(string(s), "\"", "")
	return res
}

// WeekJSON структура json-а
type WeekJSON struct {
	IsCachedResponse bool   `json:"isCachedResponse"`
	IsDataFromDB     bool   `json:"isDataFromDB"`
	DataHash         string `json:"dataHash"`
	Data             []struct {
		LessonId        int    `json:"lessonId"`
		LessonShortName string `json:"lessonShortName"`
		LessonLongName  string `json:"lessonLongName"`
		LessonType      string `json:"lessonType"`
		UpdatedAt       string `json:"updatedAt"`
		AddedAt         string `json:"addedAt"`
		Cabinet         string `json:"cabinet"`
		StartAt         string `json:"startAt"`
		EndAt           string `json:"endAt"`
		LectorShortName string `json:"lectorShortName"`
		LectorFullName  string `json:"lectorFullName"`
		GroupName       string `json:"groupName"`
		LessonDate      string `json:"lessonDate"`
		DayNameShort    string `json:"dayNameShort"`
		DayNameLong     string `json:"dayNameLong"`
	} `json:"data"`
}
