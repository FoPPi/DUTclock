package main

import (
	api "DUTclock/WorkingWithAPI"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/cmd/fyne_settings/settings"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	"github.com/pkg/browser"
	"time"
)

var (
	InternetExist   bool = true
	FacultySelector      = widget.NewSelect([]string{}, func(value string) {})
	CourseSelector       = widget.NewSelect([]string{}, func(value string) {})
	GroupSelector        = widget.NewSelect([]string{}, func(value string) {})
)

// UpdateTime Show time to pare
func UpdateTime(ParaNameLabel *widget.Label, TimerLabel *widget.Label, app fyne.App) {
	paraExist, paraName, diff := TakeTime(app)
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
		t := time.Now().String()
		api.LastUpdate = t[0:16]
		LastUpdateLabel.SetText("Updated: " + api.LastUpdate)
		OnlineLabel.SetText("Online")

		if !InternetExist {
			FacultySelector.Options = api.FacultyJSONtoString()
			CourseSelector.Options = api.CourseJSONtoString(api.FacultyName)
			GroupSelector.Options = api.GroupJSONtoString(api.CourseName)
			InternetExist = true

			CheckUpdate(w)

		}
	} else {
		if SendError {
			dialog.ShowError(err, w)
		}

		OnlineLabel.SetText("Offline")
	}
}

func CheckUpdate(w fyne.Window) {
	newVersion, release := api.CheckUpdateOnGitHub()
	if newVersion {
		dialog.ShowConfirm("Update", "New version available", func(b bool) {
			if b {
				err := browser.OpenURL(release)
				if err != nil {
					println(err)
				}
			}
		}, w)
	}
}

func main() {
	a := app.NewWithID("DUTclock")
	w := a.NewWindow("DUTclock")

	// set icon
	ic, _ := fyne.LoadResourceFromPath("Icon.png")
	w.SetIcon(ic)

	// set size
	w.Resize(fyne.NewSize(492, 492))
	w.SetFixedSize(true)

	// Time tab
	//---------------------------------------------------------------------

	// add app name
	//AppLabel := widget.NewLabel("DUTclock")

	// add name of para and timer
	ParaNameLabel := widget.NewLabel("")
	TimerLabel := widget.NewLabel("")

	// check internet connection
	status, _ := Ping("google.com")
	if status != 200 {
		InternetExist = false
	}

	//update kit
	OnlineLabel := widget.NewLabel("")
	LastUpdateLabel := widget.NewLabel("")
	UpdateButton := widget.NewButton("Update", func() {
		CheckConn(OnlineLabel, LastUpdateLabel, w, true)
		FacultySelector.Options = api.FacultyJSONtoString()

	})
	UpdateButton.Hidden = true

	// first call
	api.ReadUserConf()
	if api.GroupID != 0 {
		UpdateButton.Hidden = false
		CheckConn(OnlineLabel, LastUpdateLabel, w, false)
		UpdateTime(ParaNameLabel, TimerLabel, a)
		LastUpdateLabel.SetText("Updated: " + api.LastUpdate)
	}

	// start update timer every minute
	go func() {
		for range time.Tick(time.Minute) {
			UpdateTime(ParaNameLabel, TimerLabel, a)
		}
	}()

	// start checking connection every hour
	go func() {
		for range time.Tick(time.Hour / 2) {
			CheckConn(OnlineLabel, LastUpdateLabel, w, false)
		}
	}()

	// add selectors
	GroupLabel := widget.NewLabel("Group")
	GroupSelector = widget.NewSelect([]string{}, func(value string) {
		if value != api.GroupName {
			api.TakeGroupID(value)

			UpdateButton.Hidden = false
			CheckConn(OnlineLabel, LastUpdateLabel, w, false)
			UpdateTime(ParaNameLabel, TimerLabel, a)
			api.WriteUserConf()
		}
	})

	CourseLabel := widget.NewLabel("Course")
	CourseSelector = widget.NewSelect([]string{}, func(value string) {
		GroupSelector.Selected = ""
		GroupSelector.Options = api.GroupJSONtoString(value)
	})

	FacultyLabel := widget.NewLabel("Faculty")
	FacultySelector = widget.NewSelect([]string{}, func(value string) {
		if value != api.FacultyName {
			GroupSelector.Options = []string{}
			CourseSelector.Selected = ""
			GroupSelector.Selected = ""
			CourseSelector.Options = api.CourseJSONtoString(value)
		}
	})

	if InternetExist {
		FacultySelector.Options = api.FacultyJSONtoString()
		CheckUpdate(w)
	}

	// add selectors names if is not first start
	if api.GroupID != 0 {
		GroupSelector.Selected = api.GroupName
		CourseSelector.Selected = api.CourseName
		FacultySelector.Selected = api.FacultyName
		if InternetExist {
			CourseSelector.Options = api.CourseJSONtoString(api.FacultyName)
			GroupSelector.Options = api.GroupJSONtoString(api.CourseName)
		}
	}

	//Settings tab
	//---------------------------------------------------------------------

	LessonNameLabel := widget.NewLabel("Lesson Name")
	LessonNameRadio := widget.NewRadioGroup([]string{"Long", "Short"}, func(value string) {
		if value == "Long" {
			api.LessonName = true
		} else if value == "Short" {
			api.LessonName = false
		}
		api.WriteUserConf()
		UpdateTime(ParaNameLabel, TimerLabel, a)
	})

	LessonTypeLabel := widget.NewLabel("Lesson Type")
	LessonTypeRadio := widget.NewRadioGroup([]string{"Show", "Hide"}, func(value string) {
		if value == "Show" {
			api.LessonType = true
		} else if value == "Hide" {
			api.LessonType = false
		}
		api.WriteUserConf()
		UpdateTime(ParaNameLabel, TimerLabel, a)
	})

	//ColorPickerButton := widget.NewButton("Color", func() {
	//	s := settings.NewSettings()
	//	w := a.NewWindow("DUTclock Settings")
	//	w.SetIcon(ic)
	//
	//	appearance := s.LoadAppearanceScreen(w)
	//	tabs := container.NewAppTabs(
	//		&container.TabItem{Text: "Appearance", Icon: s.AppearanceIcon(), Content: appearance})
	//	tabs.SetTabLocation(container.TabLocationLeading)
	//	w.SetContent(tabs)
	//
	//	w.Resize(fyne.NewSize(480, 480))
	//	w.Show()
	//})

	s := settings.NewSettings()
	appearance := s.LoadAppearanceScreen(w)

	if api.LessonName == true {
		LessonNameRadio.Selected = "Long"
	} else if api.LessonName == false {
		LessonNameRadio.Selected = "Short"
	}

	if api.LessonType == true {
		LessonTypeRadio.Selected = "Show"
	} else if api.LessonType == false {
		LessonTypeRadio.Selected = "Hide"
	}

	// Add content
	//---------------------------------------------------------------------

	tabs := container.NewAppTabs(
		container.NewTabItem("Time", container.NewVBox(
			//container.NewVBox(container.NewCenter(
			//	AppLabel,
			//)),
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
		)),
		container.NewTabItem("Settings", container.NewVBox(container.NewCenter(container.NewHBox(

			container.NewVBox(
				LessonNameLabel,
				LessonNameRadio,
			),
			container.NewVBox(
				LessonTypeLabel,
				LessonTypeRadio,
			),
		)),
		//ColorPickerButton,
		),
		),
		&container.TabItem{Text: "Appearance", Content: appearance},
	)

	// Refresh theme
	tabs.OnSelected = func(t *container.TabItem) {
		t.Content.Refresh()
	}

	w.SetContent(tabs)

	// show window
	w.ShowAndRun()

}
