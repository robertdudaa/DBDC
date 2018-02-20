package main

import (
	"encoding/xml"
	s "strings"
	"os"
	"io/ioutil"
	"fmt"
	"time"
)

/*
The DeutscheBahn Developer API only respond in application/xml
 */

type Timetable struct {
	Station string `xml:"station,attr"`
	Stops []struct {
		Id string `xml:"id,attr"`
		Triplabel struct {
			Flag string `xml:"f,attr"`
			Train string `xml:"c,attr"`
			Number string `xml:"n,attr"`
		} `xml:"tl"`
		Arrival struct {
			PlannedTime string `xml:"pt,attr"`
			PlannedPlatform string `xml:"pp,attr"`
			PlannedPath string `xml:"ppth,attr"`
			ChangedTime string `xml:"ct,attr"`
			ChangedPlatform string `xml:"cp,attr"`
			ChangedPath string `xml:"cpth,attr"`
		} `xml:"ar"`
		Departure struct {
			PlannedTime string `xml:"pt,attr"`
			PlannedPlatform string `xml:"pp,attr"`
			PlannedPath string `xml:"ppth,attr"`
			ChangedTime string `xml:"ct,attr"`
			ChangedPlatform string `xml:"cp,attr"`
			ChangedPath string `xml:"cpth,attr"`
		} `xml:"dp"`
	} `xml:"s"`
}



func main() {
	/*
	planned.xml & changes.xml are some examples
	 */
	planned, err := ioutil.ReadFile("xml/planned.xml")
	changes, err := ioutil.ReadFile("xml/changes.xml")
	if err != nil {
		os.Exit(1)
	}
	var timetable Timetable
	var timetchange Timetable

	xml.Unmarshal(planned, &timetable)
	xml.Unmarshal(changes, &timetchange)

	calcDelays(timetable, timetchange)
}

func calcDelays(timetable Timetable, timetchange Timetable) {
	v := make(map[string]time.Duration)

	for _, changes := range timetchange.Stops {
		for _, pstop := range timetable.Stops {
			if s.Contains(pstop.Id, changes.Id) {
				if pstop.Departure.PlannedTime != "" && changes.Departure.ChangedTime != "" {
					ptime := convertTime(pstop.Departure.PlannedTime)
					rtime := convertTime(changes.Departure.ChangedTime)
					if (rtime.After(ptime)) {
						diff := rtime.Sub(ptime)
						v[pstop.Id] = diff
					}
				}
			}
		}
	}

	if len(v) != 0 {
		vges := time.Time{}.Add(0)
		fmt.Println("DELAYS:")
		for id, dur := range v {
			out := time.Time{}.Add(dur)
			vges = vges.Add(dur)

			fmt.Println("ID: ", id)
			fmt.Println("DELAY: ", out.Format("15:04"))
		}
		fmt.Println("\nTOTAL DELAY: ", vges.Format("15:04"))
	} else {
		fmt.Println("THERE WERE NO DELAYS")
	}
}

func convertTime(timestring string) time.Time {
	time, err := time.Parse("0601021504", timestring)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	return time
}
