package mind

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/widget"
	"image/color"
	"time"
)

func TakeRozkald(selectedDate string, sharedPrefs fyne.Preferences) (Cards [5]widget.Card) {
	var (
		nowParsedDate time.Time
		result        WeekJSON
		count         int
		isSecondJSON  bool
		WeekTypes     = []string{"CURRENT_WeekJSON", "NEXT_WeekJSON"}
	)

	if selectedDate == "now" {
		nowParsedDate = time.Now().Truncate(24 * time.Hour)
	} else {
		nowParsedDate, _ = time.Parse("02.01.2006", selectedDate)
	}

	for _, jsonName := range WeekTypes {
		result = ReadOfflineJSON(jsonName, sharedPrefs)

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
				case "09:45":
					count = 1
				case "11:45":
					count = 2
				case "13:30":
					count = 3
				case "15:15":
					count = 4
				}

				Cards[count] = widget.Card{
					Subtitle: fmt.Sprintf("(%d) %s [%s]", count+1, rec.LessonLongName, rec.LessonType),
					Content:  canvas.NewText(fmt.Sprintf(" %s - %s \t%s", rec.StartAt, rec.EndAt, rec.Cabinet), color.White),
				}
				isSecondJSON = true
			}
		}

		if isSecondJSON {
			return
		}
	}

	return [5]widget.Card{}
}

func TakeDaysFromJSON(sharedPrefs fyne.Preferences) []string {
	dates := make(map[string]bool)
	var days []string

	for _, jsonName := range []string{"CURRENT_WeekJSON", "NEXT_WeekJSON"} {
		result := ReadOfflineJSON(jsonName, sharedPrefs)
		for _, rec := range result.Data {
			if _, ok := dates[rec.LessonDate]; !ok {
				dates[rec.LessonDate] = true
				days = append(days, rec.LessonDate)
			}
		}
	}

	return days
}
