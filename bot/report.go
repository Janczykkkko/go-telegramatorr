package bot

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
)

func botGenerateReports(dblocation string) {
	r := mux.NewRouter()

	r.HandleFunc("/{int}", intHandler)

	log.Println("Reports server listening on port 8080")
	err := http.ListenAndServe(":8080", r)
	if err != nil {
		log.Println("Error starting server:", err)
	}
}

func intHandler(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	intParam := vars["int"]
	response := ""

	w.Header().Set("Content-Type", "text/plain")

	time, err := strconv.Atoi(intParam)
	if err != nil {
		http.Error(w, "Invalid integer", http.StatusBadRequest)
		return
	}
	sessions, err := GenerateReport(dblocation, time)
	if err != nil {
		response = err.Error()
		fmt.Fprintf(w, "%s", response)
	}
	fmt.Fprintf(w, "%s", sessions)

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
