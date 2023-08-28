package main

import (
	api "DUTclock/WorkingWithAPI"
	"DUTclock/mind"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/cmd/fyne_settings/settings"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"github.com/pkg/browser"
	"time"
)

var (
	InternetExist   bool = true
	FacultySelector      = widget.NewSelect(nil, nil)
	CourseSelector       = widget.NewSelect(nil, nil)
	GroupSelector        = widget.NewSelect(nil, nil)
	DateSelector         = widget.NewSelect(nil, nil)
	arrCards             = [5]widget.Card{}
	ParaNameLabel        = widget.NewLabel("")
	TimerLabel           = widget.NewLabel("")
	OnlineLabel          = widget.NewLabel("")
	LastUpdateLabel      = widget.NewLabel("")
	UpdateButton         = widget.NewButton("Update", nil)
)

func main() {
	a := app.NewWithID("DUTclock")
	api.App = a
	w := a.NewWindow("DUTclock")
	// set icon
	ic, _ := fyne.LoadResourceFromPath("Icon.png")
	w.SetIcon(ic)

	// set size
	w.Resize(fyne.NewSize(492, 492))
	//w.SetFixedSize(true)

	// check internet connection
	status, _ := mind.Ping("google.com")
	if status != 200 {
		InternetExist = false
	} else {
		mind.OwnLTSapi()
	}

	tabs := container.NewAppTabs(
		TimeTab(w),
		CalendarTab(),
		SettingsTab(w),
		AppearanceTab(w),
	)
	sharedPrefs := a.Preferences()

	tabs.Select(tabs.Items[sharedPrefs.Int("LastTabID")])

	// Refresh theme
	tabs.OnSelected = func(t *container.TabItem) {
		t.Content.Refresh()
		if sharedPrefs.Int("GroupID") != 0 {
			switch t.Text {
			case "Time":
				sharedPrefs.SetInt("LastTabID", 0)
				break
			case "Calendar":
				sharedPrefs.SetInt("LastTabID", 1)
				break
			case "Settings":
				sharedPrefs.SetInt("LastTabID", 2)
				break
			case "Appearance":
				sharedPrefs.SetInt("LastTabID", 3)
				break
			}
		}

	}
	DateSelector.OnChanged = func(value string) {
		arrCards = mind.TakeRozkald(value)
		tabs.Items[1].Content.Refresh()
	}

	tabs.Refresh()
	w.SetContent(tabs)

	// show window
	w.ShowAndRun()

}

func TimeTab(w fyne.Window) *container.TabItem {
	// start update timer every minute
	go func() {
		updateTicker := time.NewTicker(time.Minute)
		for range updateTicker.C {
			UpdateTime(ParaNameLabel, TimerLabel)
		}
	}()

	if UpdateButton == nil {
		UpdateButton = widget.NewButton("Update", func() {
			CheckConn(OnlineLabel, LastUpdateLabel, w, true)
			if InternetExist {
				FacultySelector.Options = api.FacultyJSONtoString()
			}
		})
		UpdateButton.Hide()
	}
	sharedPrefs := api.App.Preferences()
	// first call
	if sharedPrefs.Int("GroupID") != 0 {
		UpdateButton.Show()
		CheckConn(OnlineLabel, LastUpdateLabel, w, false)
		UpdateTime(ParaNameLabel, TimerLabel)
		LastUpdateLabel.SetText("Updated: " + sharedPrefs.String("LastUpdate"))
	}

	connTicker := time.NewTicker(time.Hour / 2)
	defer connTicker.Stop()

	// start checking connection every hour
	go func() {
		for range connTicker.C {
			CheckConn(OnlineLabel, LastUpdateLabel, w, false)
			arrCards = mind.TakeRozkald("now")
		}
	}()

	return container.NewTabItem("Time", container.NewCenter(container.NewVBox(
		container.NewVBox(container.NewCenter(ParaNameLabel)),
		container.NewVBox(container.NewCenter(TimerLabel)),
		container.NewVBox(container.NewCenter(LastUpdateLabel)),
		container.NewVBox(container.NewCenter(UpdateButton)),
		container.NewVBox(container.NewCenter(OnlineLabel)),
	)))
}

func CalendarTab() *container.TabItem {

	grid := container.New(layout.NewGridLayout(1), &arrCards[0], &arrCards[1], &arrCards[2], &arrCards[3], &arrCards[4])

	return container.NewTabItem("Calendar", container.NewVBox(DateSelector, grid))
}

func SettingsTab(w fyne.Window) *container.TabItem {
	sharedPrefs := api.App.Preferences()
	// add selectors
	GroupLabel := widget.NewLabel("Group")
	GroupSelector = widget.NewSelect([]string{}, func(value string) {
		if value != sharedPrefs.String("GroupName") {
			api.TakeGroupID(value)

			UpdateButton.Hidden = false

			_, err := mind.UpdateOfflineJSON()
			if err != nil {
				dialog.ShowError(err, w)
			}

			arrCards = mind.TakeRozkald("now")
			DateSelector.Options = mind.TakeDaysFromJSON()
			DateSelector.Selected = time.Now().Format("02.01.2006")

		}
	})

	CourseLabel := widget.NewLabel("Course")
	CourseSelector = widget.NewSelect([]string{}, func(value string) {
		GroupSelector.Selected = ""
		GroupSelector.Options = api.GroupJSONtoString(value)
	})

	FacultyLabel := widget.NewLabel("Faculty")
	FacultySelector = widget.NewSelect([]string{}, func(value string) {
		if value != sharedPrefs.String("FacultyName") {
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
	if sharedPrefs.Int("GroupID") != 0 {
		GroupSelector.Selected = sharedPrefs.String("GroupName")
		CourseSelector.Selected = sharedPrefs.String("CourseName")
		FacultySelector.Selected = sharedPrefs.String("FacultyName")
		arrCards = mind.TakeRozkald("now")
		DateSelector.Options = mind.TakeDaysFromJSON()
		DateSelector.Selected = time.Now().Format("02.01.2006")

		if InternetExist && len(CourseSelector.Options) == 0 {
			CourseSelector.Options = api.CourseJSONtoString(sharedPrefs.String("FacultyName"))
			GroupSelector.Options = api.GroupJSONtoString(sharedPrefs.String("CourseName"))
		}
	}

	LessonNameLabel := widget.NewLabel("Lesson Name")
	LessonNameRadio := widget.NewRadioGroup([]string{"Long", "Short"}, func(value string) {
		if value == "Long" {
			sharedPrefs.SetBool("LessonName", true)
		} else if value == "Short" {
			sharedPrefs.SetBool("LessonName", false)
		}
		UpdateTime(ParaNameLabel, TimerLabel)
	})

	LessonTypeLabel := widget.NewLabel("Lesson Type")
	LessonTypeRadio := widget.NewRadioGroup([]string{"Show", "Hide"}, func(value string) {
		if value == "Show" {
			sharedPrefs.SetBool("LessonType", true)
		} else if value == "Hide" {
			sharedPrefs.SetBool("LessonType", false)
		}
		UpdateTime(ParaNameLabel, TimerLabel)
	})

	SendNotificationLabel := widget.NewLabel("Send Notif")
	SendNotificationRadio := widget.NewRadioGroup([]string{"Yes", "No"}, func(value string) {
		if value == "Yes" {
			sharedPrefs.SetBool("SendNotification", true)
		} else if value == "No" {
			sharedPrefs.SetBool("SendNotification", false)
		}
	})

	if sharedPrefs.Bool("LessonName") == true {
		LessonNameRadio.Selected = "Long"
	} else {
		LessonNameRadio.Selected = "Short"
	}

	if sharedPrefs.Bool("LessonType") == true {
		LessonTypeRadio.Selected = "Show"
	} else {
		LessonTypeRadio.Selected = "Hide"
	}

	if sharedPrefs.Bool("SendNotification") == true {
		SendNotificationRadio.Selected = "Yes"
	} else {
		SendNotificationRadio.Selected = "No"
	}

	return container.NewTabItem("Settings", container.NewVBox(
		container.NewVBox(
			FacultyLabel,
			FacultySelector,
			CourseLabel,
			CourseSelector,
			GroupLabel,
			GroupSelector,
		),
		container.NewCenter(
			container.NewHBox(
				container.NewVBox(
					LessonNameLabel,
					LessonNameRadio,
				),
				container.NewVBox(
					LessonTypeLabel,
					LessonTypeRadio,
				),
				container.NewVBox(
					SendNotificationLabel,
					SendNotificationRadio,
				),
			),
		)))
}

func AppearanceTab(w fyne.Window) *container.TabItem {

	s := settings.NewSettings()
	appearance := s.LoadAppearanceScreen(w)

	return &container.TabItem{Text: "Appearance", Content: appearance}
}

// UpdateTime Show time to pare
func UpdateTime(ParaNameLabel *widget.Label, TimerLabel *widget.Label) {
	paraExist, paraName, diff := mind.TakeTime()
	ParaNameLabel.SetText(paraName)
	if paraExist {
		TimerLabel.SetText(diff.String())
	} else {
		TimerLabel.SetText("")
	}
}

// CheckConn try update WeekJson
func CheckConn(OnlineLabel *widget.Label, LastUpdateLabel *widget.Label, w fyne.Window, SendError bool) {
	sharedPrefs := api.App.Preferences()
	updated, err := mind.UpdateOfflineJSON()
	if updated {
		t := time.Now().String()
		sharedPrefs.SetString("LastUpdate", t[0:16])
		LastUpdateLabel.SetText("Updated: " + sharedPrefs.String("LastUpdate"))
		OnlineLabel.SetText("Online")

		if !InternetExist {
			FacultySelector.Options = api.FacultyJSONtoString()
			CourseSelector.Options = api.CourseJSONtoString(sharedPrefs.String("FacultyName"))
			GroupSelector.Options = api.GroupJSONtoString(sharedPrefs.String("CourseName"))
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

// CheckUpdate check update on GitHub
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
