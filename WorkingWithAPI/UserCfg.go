package WorkingWithAPI

import (
	"encoding/json"
	"fmt"
	"fyne.io/fyne/v2"
)

func ReadUserConf(sharedPrefs fyne.Preferences) {

	value := sharedPrefs.String("SettingsJSON")

	var SettingsJSON Settings

	err := json.Unmarshal([]byte(value), &SettingsJSON)
	if err != nil {
		fmt.Println(err)
		return
	}

	FacultyName = SettingsJSON.DutInfo.Faculty.FacultyName
	FacultyID = SettingsJSON.DutInfo.Faculty.FacultyID

	CourseName = SettingsJSON.DutInfo.Course.CourseName
	CourseID = SettingsJSON.DutInfo.Course.CourseID

	GroupName = SettingsJSON.DutInfo.Group.GroupName
	GroupID = SettingsJSON.DutInfo.Group.GroupID

	LessonName = SettingsJSON.LessonName
	LessonType = SettingsJSON.LessonType
	SendNotification = SettingsJSON.SendNotification

	LastUpdate = SettingsJSON.DutInfo.LastUpdate

	LastTabID = SettingsJSON.LastTabID

	LastApiVersion = SettingsJSON.LastApiVersion
}

func WriteUserConf(sharedPrefs fyne.Preferences) {
	data := &Settings{
		DutInfo: DutInfo{
			Faculty: Faculty{
				FacultyName: FacultyName,
				FacultyID:   FacultyID,
			},
			Course: Course{
				CourseName: CourseName,
				CourseID:   CourseID,
			},
			Group: Group{
				GroupName: GroupName,
				GroupID:   GroupID,
			},
			LastUpdate: LastUpdate,
		},
		LessonName:       LessonName,
		LessonType:       LessonType,
		SendNotification: SendNotification,
		LastTabID:        LastTabID,
		LastApiVersion:   LastApiVersion,
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		fmt.Printf("could not marshal json: %s\n", err)
		return
	}
	sharedPrefs.SetString("SettingsJSON", string(jsonData))

}

type Settings struct {
	DutInfo          DutInfo `json:"dut_info"`
	LessonName       bool    `json:"lesson_name"`
	LessonType       bool    `json:"lesson_type"`
	SendNotification bool    `json:"send_notification"`
	LastTabID        int     `json:"last_tab_id"`
	LastApiVersion   int     `json:"last_api_version"`
}

type DutInfo struct {
	Faculty    Faculty `json:"faculty"`
	Course     Course  `json:"course"`
	Group      Group   `json:"group"`
	LastUpdate string  `json:"last_update"`
}

type Faculty struct {
	FacultyName string `json:"faculty_name"`
	FacultyID   int    `json:"faculty_id"`
}

type Course struct {
	CourseName string `json:"course_name"`
	CourseID   int    `json:"course_id"`
}

type Group struct {
	GroupName string `json:"group_name"`
	GroupID   int    `json:"group_id"`
}
