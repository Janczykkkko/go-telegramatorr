package bot

import (
	"log"
	"time"
)

func TimeToReport() bool {
	now := time.Now()
	if now.Hour() == 23 && now.Minute() == 0 && now.Second() == 0 {
		return true
	}
	return false
}

func GetReport(dblocation string) (string, error) {
	if time.Now().Weekday() == time.Sunday {
		report, err := weeklyReport(dblocation)
		if err != nil {
			log.Println("Error generating report", err)
			return "", err
		}
		return report, nil
	}
	report, err := dailyReport(dblocation)
	if err != nil {
		log.Println("Error generating report", err)
		return "", err
	}
	return report, nil
}

func dailyReport(dblocation string) (string, error) {

}

func weeklyReport(dblocation string) (string, error) {

}
