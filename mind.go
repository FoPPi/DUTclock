package main

import (
	api "DUTclock/WorkingWithAPI"
	"encoding/json"
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/widget"
	"image/color"
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
		var err error
		dateNowParsed, err = time.Parse("15:04 02.01.2006", dateNow.Format("15:04 02.01.2006"))
		if err != nil {
			fmt.Println(err)
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
						paraType = "[" + rec.LessonType + "]"
					}

					if !api.SendNotification {
						if StartTime.Equal(dateNowParsed) {
							app.SendNotification(fyne.NewNotification("Пара почалася", PrettyPrint(paraName+" "+paraType)))
						} else if FinishTime.Equal(dateNowParsed) {
							app.SendNotification(fyne.NewNotification("Пара закинчилася", "Ливай нахуй"))
						}
					}

					if dateNowParsed.Before(StartTime) {
						diff = StartTime.Sub(dateNowParsed)

						return true, "До початку: " + PrettyPrint(paraName+" "+paraType), diff
					} else if dateNowParsed.Before(FinishTime) {
						diff = FinishTime.Sub(dateNowParsed)

						return true, "До кінця: " + PrettyPrint(paraName+" "+paraType), diff
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
	url := "https://dut-api.lwjerri.ml/v4/calendar/" + strconv.Itoa(Faculty) + "/" + strconv.Itoa(Course) + "/" + strconv.Itoa(Group) + "/" + Week

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

func TakeRozkald(selectedDate string) (Cards [5]widget.Card) {
	var (
		nowDate       = time.Now()
		nowParsedDate = time.Time{}
		result        WeekJSON
		count         = 0
	)

	if selectedDate == "now" {
		var err error
		nowParsedDate, err = time.Parse("02.01.2006", nowDate.Format("02.01.2006"))
		if err != nil {
			fmt.Println(err)
		}
	} else {
		var err error
		nowParsedDate, err = time.Parse("02.01.2006", selectedDate)
		if err != nil {
			fmt.Println(err)
		}
	}
	isSecondJSON := false

	result = ReadOfflineJSON("files/CURRENT_WeekJSON.json")
	for i := 1; i <= 2; i++ {
		for _, rec := range result.Data {

			paraDate, err := time.Parse("02.01.2006", rec.LessonDate)

			if err != nil {
				fmt.Println(err)
				break
			}

			if nowParsedDate.Equal(paraDate) {

				switch rec.StartAt {
				case "8:00":
					count = 0
					break
				case "09:45":
					count = 1
					break
				case "11:45":
					count = 2
					break
				case "13:30":
					count = 3
					break
				case "15:15":
					count = 4
					break

				}

				Cards[count] = widget.Card{
					Subtitle: "(" + strconv.Itoa(count+1) + ") " + rec.LessonLongName + " [" + rec.LessonType + "]",
					Content:  canvas.NewText(" "+rec.StartAt+" - "+rec.EndAt+" \t"+rec.Cabinet, color.White),
				}

				if !isSecondJSON {
					isSecondJSON = true
				}
			}
		}
		if isSecondJSON {
			return
		} else {
			isSecondJSON = true
			result = ReadOfflineJSON("files/NEXT_WeekJSON.json")
			count = 0
		}
	}
	return
}

func TakeDaysFromJSON() []string {
	isSecondJSON := false
	strArr := ""
	count := 1
	result := ReadOfflineJSON("files/CURRENT_WeekJSON.json")
	for i := 1; i <= 2; i++ {
		for _, rec := range result.Data {
			if len(strArr) == 0 {
				strArr = rec.LessonDate + " "
			}
			if !strings.Contains(strArr, rec.LessonDate) {
				strArr += rec.LessonDate + " "
				count++
			}

		}
		if isSecondJSON {
			return strings.Split(strings.TrimSuffix(strArr, " "), " ")
		} else {
			isSecondJSON = true
			result = ReadOfflineJSON("files/NEXT_WeekJSON.json")
		}
	}
	return strings.Split(strArr, " ")
}
