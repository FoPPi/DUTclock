package main

import (
	api "DUTclock/WorkingWithAPI"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	"log"
	"time"
)

// UpdateTime Show time to pare
func UpdateTime(ParaNameLabel *widget.Label, TimerLabel *widget.Label) {
	paraExist, paraName, diff := TakeTime()
	if paraExist {
		ParaNameLabel.SetText(paraName)
		TimerLabel.SetText(diff.String())
	} else {
		ParaNameLabel.SetText(paraName)
		TimerLabel.SetText("")
	}
}

// CheckConn try update WeekJson
func CheckConn(OnlineLabel *widget.Label, LastUpdateLabel *widget.Label, w fyne.Window, SendError bool) {
	updated, err := UpdateOfflineJSON()
	if updated {
		time := time.Now().String()
		api.LastUpdate = time[0:16]
		LastUpdateLabel.SetText("Updated: " + api.LastUpdate)
		OnlineLabel.SetText("Online")
	} else {
		if SendError {
			dialog.ShowError(err, w)
		}

		OnlineLabel.SetText("Offline")
	}
}

func main() {
	a := app.New()
	w := a.NewWindow("DUTclock")

	// set icon
	ic, _ := fyne.LoadResourceFromPath("Icon.png")
	w.SetIcon(ic)

	// set size
	w.Resize(fyne.NewSize(400, 400))
	w.SetFixedSize(true)

	// add app name
	AppLabel := widget.NewLabel("DUTclock")

	// add name of para and timer
	ParaNameLabel := widget.NewLabel("")
	TimerLabel := widget.NewLabel("")

	//update kit
	OnlineLabel := widget.NewLabel("")
	LastUpdateLabel := widget.NewLabel("")
	UpdateButton := widget.NewButton("Update", func() {
		CheckConn(OnlineLabel, LastUpdateLabel, w, true)
	})
	UpdateButton.Hidden = true

	// first call
	api.ReadUserConf()
	if api.GroupID != 0 {
		UpdateButton.Hidden = false
		LastUpdateLabel.SetText("Updated: " + api.LastUpdate)
		CheckConn(OnlineLabel, LastUpdateLabel, w, false)
		UpdateTime(ParaNameLabel, TimerLabel)
	}

	// start update timer every minute
	go func() {
		for range time.Tick(time.Minute) {
			UpdateTime(ParaNameLabel, TimerLabel)
		}
	}()

	// start checking connection every hour
	go func() {
		for range time.Tick(time.Hour) {
			CheckConn(OnlineLabel, LastUpdateLabel, w, false)
		}
	}()

	// add selectors
	GroupLabel := widget.NewLabel("Group")
	GroupSelector := widget.NewSelect([]string{}, func(value string) {
		api.TakeGroupID(value)

		UpdateButton.Hidden = false
		CheckConn(OnlineLabel, LastUpdateLabel, w, false)
		UpdateTime(ParaNameLabel, TimerLabel)
		api.WriteUserConf()
	})
	GroupSelector.Selected = api.GroupName

	CourseLabel := widget.NewLabel("Course")
	CourseSelector := widget.NewSelect([]string{}, func(value string) {
		GroupSelector.Selected = ""
		GroupSelector.Options = api.GroupJSONtoString(value)
	})
	CourseSelector.Selected = api.CourseName

	FacultyLabel := widget.NewLabel("Faculty")
	FacultySelector := widget.NewSelect(api.FacultyJSONtoString(), func(value string) {
		log.Println("Faculty", value)
		GroupSelector.Options = []string{}
		CourseSelector.Selected = ""
		GroupSelector.Selected = ""
		CourseSelector.Options = api.CourseJSONtoString(value)

	})
	FacultySelector.Selected = api.FacultyName

	// add content
	w.SetContent(container.NewVBox(
		container.NewVBox(container.NewCenter(
			AppLabel,
		)),
		container.NewVBox(
			FacultyLabel,
			FacultySelector,
			CourseLabel,
			CourseSelector,
			GroupLabel,
			GroupSelector,
		),
		container.NewVBox(container.NewCenter(
			ParaNameLabel,
		)),
		container.NewVBox(container.NewCenter(
			TimerLabel,
		)),
		container.NewVBox(container.NewCenter(
			LastUpdateLabel,
		)),
		container.NewVBox(container.NewCenter(
			UpdateButton,
		)),
		container.NewVBox(container.NewCenter(
			OnlineLabel,
		)),
	))

	// show window
	w.ShowAndRun()
}
