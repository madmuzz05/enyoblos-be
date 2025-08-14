package helper

import (
	"database/sql/driver"
	"fmt"
	"strings"
	"time"
)

type CustomTime struct {
	time.Time
	format string
}

// MarshalJSON tetap untuk serialisasi ke JSON
func (ct CustomTime) MarshalJSON() ([]byte, error) {
	layout := "2006-01-02 15:04:05" // default layout

	if ct.format == "date" {
		layout = "2006-01-02"
	} else if ct.format == "time" {
		layout = "15:04:05"
	}

	formatted := fmt.Sprintf("\"%s\"", ct.Format(layout))
	return []byte(formatted), nil
}

// UnmarshalJSON tetap untuk parsing JSON ke struct
func (ct *CustomTime) UnmarshalJSON(b []byte) error {
	s := strings.Trim(string(b), "\"")

	var layout string

	switch {
	case strings.Contains(s, " "):
		layout = "2006-01-02 15:04:05"
	case len(s) == len("2006-01-02"):
		layout = "2006-01-02"
		ct.format = "date"
	case len(s) == len("15:04:05"):
		layout = "15:04:05"
		ct.format = "time"
	default:
		return fmt.Errorf("unsupported time format: %s", s)
	}

	parsedTime, err := time.Parse(layout, s)
	if err != nil {
		return err
	}

	ct.Time = parsedTime
	return nil
}

func (ct CustomTime) Value() (driver.Value, error) {
	return ct.Time, nil
}

func (ct *CustomTime) Scan(value interface{}) error {
	if value == nil {
		*ct = CustomTime{Time: time.Time{}}
		return nil
	}

	switch v := value.(type) {
	case time.Time:
		*ct = CustomTime{Time: v}
	case string:
		parsedTime, format, err := parseCustomTimeFormats(v)
		if err != nil {
			return fmt.Errorf("cannot convert %v to CustomTime: %v", value, err)
		}
		*ct = CustomTime{Time: parsedTime, format: format}
	default:
		return fmt.Errorf("cannot convert %v to CustomTime", value)
	}

	if ct.format == "" {
		ct.format = detectFormat(ct.Time)
	}

	return nil
}

func parseCustomTimeFormats(value string) (time.Time, string, error) {
	var parsedTime time.Time
	var format string
	var err error

	switch {
	case strings.Contains(value, " "):
		format = "datetime"
		parsedTime, err = time.Parse("2006-01-02 15:04:05", value)
	case len(value) == len("2006-01-02"):
		format = "date"
		parsedTime, err = time.Parse("2006-01-02", value)
	case len(value) == len("15:04:05"):
		format = "time"
		parsedTime, err = time.Parse("15:04:05", value)
	default:
		err = fmt.Errorf("unsupported time format: %s", value)
	}

	return parsedTime, format, err
}

func detectFormat(t time.Time) string {
	if t.Hour() == 0 && t.Minute() == 0 && t.Second() == 0 {
		return "date"
	}
	return "datetime"
}

// Fungsi Now menggantikan NewCurrentTime untuk mengembalikan waktu saat ini di GMT+7
func Now() CustomTime {
	loc, err := time.LoadLocation("Asia/Jakarta")
	if err != nil {
		loc = time.FixedZone("GMT+7", 7*60*60)
	}
	currentTime := time.Now().In(loc).Format("2006-01-02 15:04:05")
	newCurrentTime, _ := time.ParseInLocation("2006-01-02 15:04:05", currentTime, loc)
	return CustomTime{Time: newCurrentTime}
}

// Fungsi untuk format tanggal "Y-m-d"
func (ct CustomTime) FormatDate() CustomTime {
	time, _ := time.ParseInLocation("2006-01-02", ct.Format("2006-01-02"), ct.Location())
	ct.Time = time
	ct.format = "date"
	return ct
}

// Fungsi untuk format waktu "H:i:s"
func (ct CustomTime) FormatTime() CustomTime {
	time, _ := time.ParseInLocation("15:04:05", ct.Format("15:04:05"), ct.Location())
	ct.Time = time
	ct.format = "time"
	return ct
}

// Fungsi untuk format lengkap "Y-m-d H:i:s"
func (ct CustomTime) FormatDateTime() CustomTime {
	time, _ := time.ParseInLocation("2006-01-02 15:04:05", ct.Format("2006-01-02 15:04:05"), ct.Location())
	ct.Time = time
	ct.format = "datetime"
	return ct
}

// Fungsi ToString untuk mengembalikan string sesuai format
func (ct CustomTime) ToString() string {
	var layout string

	switch ct.format {
	case "date":
		layout = "2006-01-02"
	case "time":
		layout = "15:04:05"
	default:
		layout = "2006-01-02 15:04:05"
	}

	formatted := ct.Format(layout)
	return formatted
}

func (ct CustomTime) FormatID() string {
	bulan := []string{
		"Januari", "Februari", "Maret", "April", "Mei", "Juni",
		"Juli", "Agustus", "September", "Oktober", "November", "Desember",
	}
	hari := []string{
		"Minggu", "Senin", "Selasa", "Rabu", "Kamis", "Jumat", "Sabtu",
	}
	return fmt.Sprintf("%s, %02d %s %d",
		hari[ct.Weekday()], ct.Day(), bulan[ct.Month()-1], ct.Year())
}

func (ct CustomTime) FormatDateTimeID() string {
	bulan := []string{
		"Januari", "Februari", "Maret", "April", "Mei", "Juni",
		"Juli", "Agustus", "September", "Oktober", "November", "Desember",
	}
	hari := []string{
		"Minggu", "Senin", "Selasa", "Rabu", "Kamis", "Jumat", "Sabtu",
	}
	return fmt.Sprintf("%s, %02d %s %d %s",
		hari[ct.Weekday()], ct.Day(), bulan[ct.Month()-1], ct.Year(), ct.Format("15:04:05"))
}

// ParseStringToCustomTime converts a string to CustomTime based on its format.
func ParseStringToCustomTime(input string) (CustomTime, error) {
	input = strings.TrimSpace(input)

	// Define supported formats with their layouts and types
	formats := []struct {
		layout string
		format string
	}{
		// DateTime formats
		{"2006-01-02 15:04:05", "datetime"},
		{"2006/01/02 15:04:05", "datetime"},
		{"2006.01.02 15:04:05", "datetime"},
		{"2006-01-02T15:04:05", "datetime"},
		{"2006-01-02T15:04:05Z", "datetime"},
		{"2006-01-02T15:04:05-07:00", "datetime"},

		// Date formats
		{"2006-01-02", "date"},
		{"2006/01/02", "date"},
		{"2006.01.02", "date"},
		{"02-01-2006", "date"},
		{"02/01/2006", "date"},
		{"02.01.2006", "date"},
		{"01-02-2006", "date"},
		{"01/02/2006", "date"},
		{"02-01-06", "date"},
		{"02/01/06", "date"},
		{"06-01-02", "date"},
		{"06/01/02", "date"},
		{"20060102", "date"},

		// Time formats
		{"15:04:05", "time"},
		{"15:04", "time"},
		{"3:04:05 PM", "time"},
		{"3:04 PM", "time"},
	}

	// Try each format until one works
	for _, f := range formats {
		if parsedTime, err := time.Parse(f.layout, input); err == nil {
			return CustomTime{Time: parsedTime, format: f.format}, nil
		}
	}

	return CustomTime{}, fmt.Errorf("unsupported time format: %s", input)
}

func (ct CustomTime) FormatMonthYearID() string {
	bulan := []string{
		"Januari", "Februari", "Maret", "April", "Mei", "Juni",
		"Juli", "Agustus", "September", "Oktober", "November", "Desember",
	}
	return fmt.Sprintf("%s-%d", bulan[ct.Month()-1], ct.Year())
}
