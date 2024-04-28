package helpers

import "time"

func FormatDateToYMD(date *time.Time) string {
	if date != nil {
		return date.Format("2006-01-02")
	}
	return ""
}

func FormatStringToDate(stringDate string) (time.Time, error) {
	dateNow := "2006-01-02"
	birthDateParse, err := time.Parse(dateNow, stringDate)
	return birthDateParse, err

}
