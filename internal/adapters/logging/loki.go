package logging

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/golang-auth/internal/core/ports"
)

type lokiAdapter struct {
	url    string
	labels map[string]string
}

func NewLokiLogger(url string, labels map[string]string) ports.Logger {
	return &lokiAdapter{
		url:    url,
		labels: labels,
	}
}

func (l *lokiAdapter) send(level, msg string, args ...any) {
	// 1. Start with the basic line
	line := fmt.Sprintf("level=%s msg=%q", level, msg)

	// 2. Loop through args to create individual rows/keys
	// We assume args come in pairs: ["key", value, "key2", value2]
	for i := 0; i < len(args); i += 2 {
		if i+1 < len(args) {
			key := fmt.Sprintf("%v", args[i])
			val := fmt.Sprintf("%v", args[i+1])
			// Append each as its own key=value pair
			line = fmt.Sprintf("%s %s=%q", line, key, val)
		} else {
			// If there's an odd number of args, just append the last one
			line = fmt.Sprintf("%s extra=%q", line, args[i])
		}
	}

	ts := fmt.Sprintf("%d", time.Now().UnixNano())
	lokiMsg := map[string]any{
		"streams": []map[string]any{
			{
				"stream": l.labels,
				"values": [][]string{{ts, line}},
			},
		},
	}

	body, _ := json.Marshal(lokiMsg)
	// Suggestion: Use a persistent client with a timeout
	http.Post(l.url, "application/json", bytes.NewBuffer(body))
}

func (l *lokiAdapter) Debug(msg string, args ...any) {
	l.send("Debug", msg, args...)
}

func (l *lokiAdapter) Info(msg string, args ...any) {
	l.send("Info", msg, args...)
}

func (l *lokiAdapter) Warn(msg string, args ...any) {
	l.send("Warn", msg, args...)
}

func (l *lokiAdapter) Error(msg string, args ...any) {
	l.send("Error", msg, args...)
}

func (l *lokiAdapter) Fatal(msg string, args ...any) {
	l.send("Fatal", msg, args...)
	os.Exit(1)
}
