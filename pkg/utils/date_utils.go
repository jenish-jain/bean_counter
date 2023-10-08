package utils

import (
	"fmt"
	"regexp"
	"strconv"
	"time"
	_ "time/tzdata"
)

func ParseMDYYYYToDate(dateString string) time.Time {
	re := regexp.MustCompile("([0-9]+)/([0-9]+)/([0-9]+)")
	result := re.FindAllStringSubmatch(dateString, -1)
	month, _ := strconv.Atoi(result[0][1])
	day, _ := strconv.Atoi(result[0][2])
	year, _ := strconv.Atoi(result[0][3])
	return time.Date(year, time.Month(month), day, 0, 0, 0, 0, asiaKolkataTimeZone())
}

func GetPreviousMonth() time.Month {
	return time.Now().In(asiaKolkataTimeZone()).Month() - 1
}

func GetCurrentYear() int {
	return time.Now().In(asiaKolkataTimeZone()).Year()
}

func asiaKolkataTimeZone() *time.Location {
	location, err := time.LoadLocation("Asia/Kolkata")
	if err != nil {
		fmt.Printf("Error getting asia/kolkata location %+v \n", err)
		panic("unable to load asia kolkata location")
	}

	return location
}
