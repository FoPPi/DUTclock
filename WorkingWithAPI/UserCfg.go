package WorkingWithAPI

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
)

func ReadUserConf() {
	jsonFile, err := os.Open("files/UserCfg.json")
	// if we os.Open returns an error then handle it
	if err != nil {
		fmt.Println(err)
		return
	}

	var UserCFG UserJSON

	byteValue, _ := ioutil.ReadAll(jsonFile)

	err = json.Unmarshal(byteValue, &UserCFG)
	if err != nil {
		fmt.Println(err)
		return
	}

	err = jsonFile.Close()
	if err != nil {
		fmt.Println(err)
		return
	}

	FacultyName = UserCFG.Faculty.FacultyName
	FacultyID = UserCFG.Faculty.FacultyID

	CourseName = UserCFG.Course.CourseName
	CourseID = UserCFG.Course.CourseID

	GroupName = UserCFG.Group.GroupName
	GroupID = UserCFG.Group.GroupID

	LessonName = UserCFG.Settings.LessonName
	LessonType = UserCFG.Settings.LessonType
	SendNotification = UserCFG.Settings.SendNotification

	LastUpdate = UserCFG.LastUpdate
}

func WriteUserConf() {
	data := &UserJSON{
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
		Settings: Settings{
			LessonName:       LessonName,
			LessonType:       LessonType,
			SendNotification: SendNotification,
		},
		LastUpdate: LastUpdate,
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		fmt.Printf("could not marshal json: %s\n", err)
		return
	}

	var nestedDir = "files"
	path := filepath.Join(".", nestedDir)
	erro := os.MkdirAll(path, 0777)
	if erro != nil {
		log.Fatal(erro)
	}

	_ = ioutil.WriteFile("files/UserCfg.json", jsonData, 0644)
}

type UserJSON struct {
	Faculty    Faculty  `json:"faculty"`
	Course     Course   `json:"course"`
	Group      Group    `json:"group"`
	Settings   Settings `json:"settings"`
	LastUpdate string   `json:"last_update"`
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

type Settings struct {
	LessonName       bool `json:"lesson_name"`
	LessonType       bool `json:"lesson_type"`
	SendNotification bool `json:"send_notification"`
}
