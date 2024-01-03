package bot

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
)

func botGenerateReports(dblocation string) {
	http.HandleFunc("/", reportHandler)
	fmt.Println("Server listening on port 8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Println("Error starting server:", err)
	}
}

func reportHandler(w http.ResponseWriter, r *http.Request) {
	// Get the integer value from the query parameter or set a default value
	time := r.URL.Query().Get("int")
	if time == "" {
		time = "24" // Default value if no parameter is provided
	}

	timeInt, err := strconv.Atoi(time)
	if err != nil {
		http.Error(w, "Invalid integer", http.StatusBadRequest)
		return
	}

	report, err := GenerateReport(dblocation, timeInt)
	if err != nil {
		report = err.Error()
	}
	htmlReport := "<p>" + strings.ReplaceAll(report, "\n", "<br>") + "</p>"
	// Generate HTML response dynamically with a form field and default background color
	htmlResponse := fmt.Sprintf(`
		<!DOCTYPE html>
		<html>
		<head>
			<title>Player sessions</title>
		</head>
		<body style="background-color: black; color: white;">
			<h1>Player sessions</h1>
			<form action="/" method="get">
				<label for="intInput">Enter an Integer:</label>
				<input type="number" id="intInput" name="int" value="%d">
				<input type="submit" value="Submit">
			</form>
			<p>%s</p>
		</body>
		</html>
	`, timeInt, htmlReport)

	// Write the HTML response to the client
	w.Header().Set("Content-Type", "text/html")
	w.Write([]byte(htmlResponse))
}

func GenerateReport(dblocation string, time int) (string, error) {
	showDays := false
	var report strings.Builder
	if time > 24 {
		showDays = true
	}
	sessiondata, err := GetSessionsByUserFromDB(dblocation, time)
	if err != nil {
		log.Println("Error getting data from db", err)
		return "", err
	}
	report.WriteString(fmt.Sprintf("Here's %d hour report from media players:\n", time))
	report.WriteString("-------------\n")
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
