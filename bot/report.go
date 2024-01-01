package bot

import (
	"fmt"
	"log"
	"strings"
	"time"
)

func TimeToReport() bool {
	now := time.Now()
	if now.Hour() == 23 && now.Minute() == 0 && now.Second() == 0 {
		return true
	}
	return false
}

func GenerateReport(dblocation string) (string, error) {
	reportTime := 24 //default daily
	showDays := false
	var report strings.Builder
	if time.Now().Weekday() == time.Sunday {
		reportTime = 168 //weekly
	}
	sessiondata, err := GetSessionsByUserFromDB(dblocation, reportTime)
	if err != nil {
		log.Println("Error getting data from db", err)
		return "", err
	}
	if reportTime == 24 {
		report.WriteString("Here's a daily report from media players:\n")
	} else {
		report.WriteString("Here's a weekly report from media players:\n")
		showDays = true
	}

	for _, userSessions := range sessiondata {
		report.WriteString(fmt.Sprintf("User: %s\n", userSessions.UserName))
		for _, s := range userSessions.Sessions {
			report.WriteString(fmt.Sprintf("%s - %s on %s(%s) for %s minutes\nmethod: %s\nbitrate: %s Mbps\nsubs: %s\n",
				FormatTimeToNiceString(s.StartTime, showDays),
				s.Name,
				s.Service,
				s.DeviceName,
				s.Duration,
				s.PlayMethod,
				s.Bitrate,
				s.SubStream))
		}
		report.WriteString("-------------\n")
	}
	return report.String(), nil
}
