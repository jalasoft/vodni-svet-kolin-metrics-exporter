package vodnisvetkolin

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"regexp"
	"strconv"
)

var (
	url                    string
	occupancyRegexp        *regexp.Regexp
	visitedTodayRegexp     *regexp.Regexp
	visitedThisWeekRegexp  *regexp.Regexp
	visitedThisMonthRegexp *regexp.Regexp
	visitedThisYearRegexp  *regexp.Regexp
)

func init() {
	flag.StringVar(&url, "url", "https://www.vodnisvetkolin.cz/", "URL")

	occupancyRegexp = regexp.MustCompile(`<div class="bubble">\s*<div class="value">\s+(\d+)\s+</div>`)
	visitedTodayRegexp = regexp.MustCompile(`<div class="desc">\s*dnes v aquaparku:\s*</div>\s*<div class="value">\s*(\d+)\s*</div>`)
	visitedThisWeekRegexp = regexp.MustCompile(`<div class="desc">\s*tento týden v aquaparku:\s*</div>\s*<div class="value">\s*(\d+)\s*</div>`)
	visitedThisMonthRegexp = regexp.MustCompile(`<div class="desc">\s*tento měsíc v aquaparku:\s*</div>\s*<div class="value">\s*(\d+)\s*</div>`)
	visitedThisYearRegexp = regexp.MustCompile(`<div class="desc">\s*tento rok v aquaparku:\s*</div>\s*<div class="value">\s*(\d+)\s*</div>`)
}

type VodniSvetKolin struct {
	Occupancy        uint16
	VisitedToday     uint32
	VisitedThisWeek  uint32
	VisitedThisMonth uint32
	VisitedThisYear  uint32
}

func (v VodniSvetKolin) String() string {
	return fmt.Sprintf("VodniSvetKolin[now:%d, today:%d, this week:%d, this month:%d, this year:%d]", v.Occupancy, v.VisitedToday, v.VisitedThisWeek, v.VisitedThisMonth, v.VisitedThisYear)
}

func ReadStatistics() (VodniSvetKolin, error) {
	page, err := loadPage()

	if err != nil {
		log.Printf("Error: %s", err.Error())
		return VodniSvetKolin{}, err
	}

	occupancy, err := findUint32(page, occupancyRegexp)

	if err != nil {
		log.Printf("Error: %s", err.Error())
		return VodniSvetKolin{}, err
	}

	visitedToday, err := findUint32(page, visitedTodayRegexp)

	if err != nil {
		log.Printf("Error: %s", err.Error())
		return VodniSvetKolin{}, err
	}

	visitedThisWeek, err := findUint32(page, visitedThisWeekRegexp)

	if err != nil {
		log.Printf("Error: %s", err.Error())
		return VodniSvetKolin{}, err
	}

	visitedThisMonth, err := findUint32(page, visitedThisMonthRegexp)

	if err != nil {
		log.Printf("Error: %s", err.Error())
		return VodniSvetKolin{}, err
	}

	visitedThisYear, err := findUint32(page, visitedThisYearRegexp)

	if err != nil {
		log.Printf("Error: %s", err.Error())
		return VodniSvetKolin{}, err
	}

	result := VodniSvetKolin{
		Occupancy:        uint16(occupancy),
		VisitedToday:     uint32(visitedToday),
		VisitedThisWeek:  uint32(visitedThisWeek),
		VisitedThisMonth: uint32(visitedThisMonth),
		VisitedThisYear:  uint32(visitedThisYear),
	}

	log.Printf("%s", result)

	return result, nil
}

func loadPage() (string, error) {

	respo, err := http.Get(url)

	if err != nil {
		return "", err
	}

	if respo.StatusCode != 200 {
		return "", fmt.Errorf("Cannot load page, having response status %d\n", respo.StatusCode)
	}

	defer respo.Body.Close()

	var buff *bytes.Buffer = new(bytes.Buffer)
	_, err = io.Copy(buff, respo.Body)

	if err != nil {
		return "", fmt.Errorf("Cannot copy response body: %v\n", err)
	}

	return buff.String(), nil
}

func findUint32(page string, pattern *regexp.Regexp) (uint32, error) {
	matches := pattern.FindAllStringSubmatch(page, -1)

	if len(matches) == 0 {
		return 0, fmt.Errorf("Cannot find occupancy pattern. Reading from '%s'", url)
	}

	if len(matches) > 1 {
		return 0, fmt.Errorf("More than 1 occurrence of occupancy pattern. Reading from '%s'", url)
	}

	var groups []string = matches[0]

	number, err := strconv.Atoi(groups[1])

	if err != nil {
		return 0, fmt.Errorf("Cannot convert string '%s' to an integer.", groups[1])
	}

	return uint32(number), nil
}
