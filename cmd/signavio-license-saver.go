package main

import (
	"flag"
	"log/slog"
	"os"
	"slices"
	"strings"
	"time"

	"github.com/torbendury/signavio-license-saver/pkg/signavio"
)

var signavioTenant string
var signavioUser string
var signavioPassword string
var signavioURL string
var allowlist string

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	flag.StringVar(&signavioTenant, "tenant", "", "Signavio tenant ID")
	flag.StringVar(&signavioUser, "user", "", "Signavio user")
	flag.StringVar(&signavioPassword, "password", "", "Signavio password")
	flag.StringVar(&signavioURL, "url", "", "Signavio URL (e.g. https://editor.signavio.com, no trailing slash)")
	flag.StringVar(&allowlist, "allowlist", "", "Optional: Comma separated list of emails to keep")
	flag.Parse()

	if signavioTenant == "" || signavioUser == "" || signavioPassword == "" || signavioURL == "" {
		flag.PrintDefaults()
		os.Exit(1)
	}

	if allowlist == "" {
		logger.Info("No allowlist provided, at least keeping API user")
		allowlist = signavioUser
	}

	allowedUsers := strings.Split(allowlist, ",")
	if !slices.Contains(allowedUsers, signavioUser) {
		logger.Info("API user is not in allowlist, adding for safety")
		allowedUsers = append(allowedUsers, signavioUser)
	}

	client := signavio.New(signavioTenant, signavioURL, signavioUser, signavioPassword, logger)
	err := client.Login()
	if err != nil {
		logger.Error("Login failed", "error", err)
		os.Exit(1)
	}
	logger.Info("Logged in")

	users, err := client.GetUsers()
	if err != nil {
		logger.Error("GetUsers failed", "error", err)
		os.Exit(1)
	}
	logger.Info("Retrieved users", "users", len(*users))

	for _, user := range *users {
		if user.Rep.Email == "" {
			continue
		}
		if !slices.Contains(allowedUsers, user.Rep.Email) {
			logger.Info("Deleting user", "email", user.Rep.Email)
			job, err := client.DeleteUser(user)
			if err != nil {
				logger.Error("DeleteUser failed", "error", err)
				continue
			}
			job, err = client.GetJobStatus(job)
			if err != nil {
				logger.Error("GetJobStatus failed", "error", err)
				os.Exit(1)
			}
			for job.Status == string(signavio.Scheduled) || job.Status == string(signavio.Running) {
				job, err = client.GetJobStatus(job)
				if err != nil {
					logger.Error("GetJobStatus failed", "error", err)
				}
				time.Sleep(500 * time.Millisecond)
			}
			logger.Info("Job finished", "job", *job)

		} else {
			logger.Info("Keeping user", "email", user.Rep.Email)
		}
	}
}
