package main

import (
	"github.com/PuerkitoBio/goquery"
	"golang.org/x/text/encoding/japanese"
	"regexp"
	"strconv"
	"time"
)

// Lesson stands for lesson of UEC.
type Lesson struct {
	Class   string    `json:"class"`
	Date    time.Time `json:"date"`
	Period  int       `json:"period"`
	Subject string    `json:"subject"`
	Teacher string    `json:"teacher"`
	Remark  string    `json:"remark"`
}

// LessonsSlice is struct for json.
type LessonsSlice struct {
	Lessons []Lesson `json:"lessons"`
}

// URLs are.
const (
	UndergraduateURL = "http://kyoumu.office.uec.ac.jp/kyuukou/kyuukou.html"
	GraduateURL      = "http://kyoumu.office.uec.ac.jp/kyuukou/kyuukou2.html"
)

// ToArray converts Lesson struct to an array.
func (l Lesson) ToArray() []string {
	lesson := []string{
		l.Class,
		l.Date.String(),
		strconv.Itoa(l.Period),
		l.Subject,
		l.Teacher,
		l.Remark,
	}

	return lesson
}

// GetLessons get kyuuko lessons from url.
func GetLessons(url string) []Lesson {
	doc, _ := goquery.NewDocument(url)
	lessons := []Lesson{}
	doc.Find("table > tbody > tr").Next().Each(func(_ int, row *goquery.Selection) {
		elems := make([]string, 6)
		row.Find("td").Each(func(i int, s *goquery.Selection) {
			raw, _ := s.Html()
			if raw != string(0xA0) { // nbsp避け
				text, _ := japanese.ShiftJIS.NewDecoder().String(raw)
				elems[i] = text
			}
		})

		period, _ := strconv.Atoi(elems[2])

		monthAndDay := regexp.MustCompile("(\\d+)月(\\d+)日").FindStringSubmatch(elems[1])
		month, _ := strconv.Atoi(monthAndDay[1])
		day, _ := strconv.Atoi(monthAndDay[2])
		loc, _ := time.LoadLocation("Asia/Tokyo")
		date := time.Date(time.Now().Year(), time.Month(month), day, 0, 0, 0, 0, loc)

		lessons = append(lessons, Lesson{
			Class:   elems[0],
			Date:    date,
			Period:  period,
			Subject: elems[3],
			Teacher: elems[4],
			Remark:  elems[5],
		})
	})

	return lessons
}
