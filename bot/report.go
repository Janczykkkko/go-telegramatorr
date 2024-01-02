package bot

import (
	"fmt"
	"log"
	"strings"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func botGenerateReports(chatID int64, bot *tgbotapi.BotAPI, dblocation string) {
	//periodically check if should report & report
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()
	for range ticker.C {
		if TimeToReport() {
			report, err := GenerateReport(dblocation)
			if err != nil {
				log.Println("Error generating report", err)
				// Sending an error message instead
				errorMsg := tgbotapi.NewMessage(chatID, "Failed to generate a report, please check the logs")
				errorMsg.DisableNotification = true
				if _, err := bot.Send(errorMsg); err != nil {
					log.Printf("Error sending error message: %s", err)
				}
				continue // Skips sending message
			}
			// Sending the message if no error occurred during report generation
			msg := tgbotapi.NewMessage(chatID, report)
			msg.DisableNotification = true
			if _, err := bot.Send(msg); err != nil {
				log.Printf("Error sending report: %s", err)
				continue
			}
			log.Println("Report succesfully generated and sent\n", report)
		}
	}
}

func TimeToReport() bool {
	now := time.Now()
	if now.Hour() == 10 && now.Minute() == 0 && now.Second() == 0 {
		return true
	}
	return false
}

func GenerateReport(dblocation string) (string, error) {
	reportTime := 24 //default daily
	showDays := false
	cleandb := false
	var report strings.Builder
	if time.Now().Weekday() == time.Sunday {
		reportTime = 168 //weekly
		showDays = true
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
		cleandb = true
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
	if cleandb {
		CleanDB(dblocation) //clean db after a weekly report
		CreateDb(dblocation)
	}
	return report.String(), nil
}
