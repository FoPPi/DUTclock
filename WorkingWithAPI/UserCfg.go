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
	}

	var UserCFG UserJSON

	byteValue, _ := ioutil.ReadAll(jsonFile)

	json.Unmarshal(byteValue, &UserCFG)

	jsonFile.Close()

	FacultyName = UserCFG.Faculty.FacultyName
	FacultyID = UserCFG.Faculty.FacultyID

	CourseName = UserCFG.Course.CourseName
	CourseID = UserCFG.Course.CourseID

	GroupName = UserCFG.Group.GroupName
	GroupID = UserCFG.Group.GroupID

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
		LastUpdate: LastUpdate,
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		fmt.Printf("could not marshal json: %s\n", err)
		return
	}

	//_, erro := os.Stat("files")
	//
	//if os.IsNotExist(err) {
	//	errDir := os.Mkdir("./files", os.FileMode(0522))
	//	if errDir != nil {
	//		log.Fatal(erro)
	//	}
	//
	//}

	var nestedDir = "files"
	path := filepath.Join(".", nestedDir)
	erro := os.MkdirAll(path, 0777)
	if erro != nil {
		log.Fatal(erro)
	}

	_ = ioutil.WriteFile("files/UserCfg.json", jsonData, 0644)
}

type UserJSON struct {
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
