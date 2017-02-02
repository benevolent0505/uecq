package main

import (
	"github.com/PuerkitoBio/goquery"
	"golang.org/x/text/encoding/japanese"
	"regexp"
	"strconv"
	"time"
)

// Lecture stands for lecture of UEC.
type Lecture struct {
	Class   string    `json:"class"`
	Date    time.Time `json:"date"`
	Period  int       `json:"period"`
	Subject string    `json:"subject"`
	Teacher string    `json:"teacher"`
	Remark  string    `json:"remark"`
}

// LecturesSlice is struct for json.
type LecturesSlice struct {
	Lectures []Lecture `json:"lectures"`
}

// URLs are.
const (
	UndergraduateURL = "http://kyoumu.office.uec.ac.jp/kyuukou/kyuukou.html"
	GraduateURL      = "http://kyoumu.office.uec.ac.jp/kyuukou/kyuukou2.html"
)

// ToArray converts Lecture struct to an array.
func (l Lecture) ToArray() []string {
	lecture := []string{
		l.Class,
		l.Date.String(),
		strconv.Itoa(l.Period),
		l.Subject,
		l.Teacher,
		l.Remark,
	}

	return lecture
}

// GetLectures get kyuuko lectures from url.
func GetLectures(url string) []Lecture {
	doc, _ := goquery.NewDocument(url)
	lectures := []Lecture{}
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

		lectures = append(lectures, Lecture{
			Class:   elems[0],
			Date:    date,
			Period:  period,
			Subject: elems[3],
			Teacher: elems[4],
			Remark:  elems[5],
		})
	})

	return lectures
}
