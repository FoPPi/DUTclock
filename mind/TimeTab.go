package mind

import (
	api "DUTclock/WorkingWithAPI"
	"encoding/json"
	"fmt"
	"fyne.io/fyne/v2"
	"strings"
	"time"
)

// TakeTime показывет сколько до начала/конца пары (надо оптимизировать)
func TakeTime(app fyne.App) (paraExist bool, paraName string, diff time.Duration) {
	dateNow := time.Now()

	if isWeekend(dateNow) {
		return false, "Сьогодні немає пар :)", diff
	}

	dateNowParsed, err := time.Parse("15:04 02.01.2006", dateNow.Format("15:04 02.01.2006"))
	if err != nil {
		fmt.Println(err)
		return false, "Дикая ошибка", diff
	}

	for _, jsonName := range []string{"CURRENT_WeekJSON", "NEXT_WeekJSON"} {
		result := ReadOfflineJSON(jsonName, app.Preferences())
		for _, rec := range result.Data {
			StartTime, err := time.Parse("15:04 02.01.2006", rec.StartAt+" "+rec.LessonDate)
			if err != nil {
				fmt.Println(err)
				continue
			}
			if dateNowParsed.Day() != StartTime.Day() || dateNowParsed.Month() != StartTime.Month() || dateNowParsed.Year() != StartTime.Year() {
				continue
			}
			FinishTime, err := time.Parse("15:04 02.01.2006", rec.EndAt+" "+rec.LessonDate)
			if err != nil {
				fmt.Println(err)
				continue
			}
			paraName, paraType := rec.LessonShortName, ""
			if api.LessonName {
				paraName = rec.LessonLongName
			}
			if api.LessonType {
				paraType = "[" + rec.LessonType + "]"
			}
			if api.SendNotification {
				if StartTime.Equal(dateNowParsed) {
					app.SendNotification(fyne.NewNotification("Пара почалася", PrettyPrint(paraName+" "+paraType)))
				} else if FinishTime.Equal(dateNowParsed) {
					app.SendNotification(fyne.NewNotification("Пара закінчилася", "Ливай"))
				}
			}
			if dateNowParsed.Before(StartTime) {
				diff = StartTime.Sub(dateNowParsed)
				return true, "До початку: " + PrettyPrint(paraName+" "+paraType), diff
			} else if dateNowParsed.Before(FinishTime) {
				diff = FinishTime.Sub(dateNowParsed)
				return true, "До кінця: " + PrettyPrint(paraName+" "+paraType), diff
			}
		}
		result = ReadOfflineJSON("NEXT_WeekJSON", app.Preferences())
	}
	return false, "Пари закінчилися :)", diff
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
