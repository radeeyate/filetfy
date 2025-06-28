package main

import (
	"fmt"
	"log"
	"log/slog"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/go-co-op/gocron/v2"
	"github.com/joho/godotenv"
	"resty.dev/v3"
)

type Directory struct {
	Path  string
	Files []os.DirEntry
}

var (
	oldDirectories []Directory
)

func main() {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	client := resty.New()
	defer client.Close()

	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}

	ntfyTopic := os.Getenv("NTFY_TOPIC")
	ntfyServer := os.Getenv("NTFY_SERVER")
	cronString := os.Getenv("CRON_STRING")
	directories := os.Getenv("DIRECTORIES")
	createIndicator := os.Getenv("CREATE_INDICATOR")
	deleteIndicator := os.Getenv("DELETE_INDICATOR")

	directoriesList := strings.Split(directories, ",")

	if len(directoriesList) == 0 {
		panic("DIRECTORIES is not set")
	}

	for _, directory := range directoriesList {
		if exists, err := exists(directory); err != nil {
			log.Fatal(err)
		} else if !exists {
			log.Fatal("Directory does not exist: " + directory)
		}

		files, err := os.ReadDir(directory)
		if err != nil {
			log.Fatal(err)
		}

		oldDirectories = append(oldDirectories, Directory{
			Path:  directory,
			Files: files,
		})

		slog.Info("Directory registered", "dir", directory)
	}


	if ntfyServer == "" {
		slog.Info("NTFY_SERVER is not set, defaulting to ntfy.sh")
		ntfyServer = "https://ntfy.sh/"
	}

	if ntfyTopic == "" {
		panic("NTFY_TOPIC is not set")
	}

	if cronString == "" {
		panic("CRON_STRING is not set")
	}

	if createIndicator == "" {
		createIndicator = "+"
		slog.Info("CREATE_INDICATOR is not set, defaulting to +")
	}

	if deleteIndicator == "" {
		deleteIndicator = "-"
		slog.Info("DELETE_INDICATOR is not set, defaulting to -")
	}

	s, err := gocron.NewScheduler()
	if err != nil {
		log.Fatal(err)
	}

	j, err := s.NewJob(
		gocron.CronJob(cronString, true),
		gocron.NewTask(
			func() {
				totalNew := 0
				totalDeleted := 0
				finalMessage := ""

				for index, directory := range oldDirectories {
					message := "- " + directory.Path
					files, err := os.ReadDir(directory.Path)
					if err != nil {
						slog.Error("Error reading directory", "err", err)
					}

					new := difference(files, directory.Files)
					deleted := difference(directory.Files, files)

					totalNew += len(new)
					totalDeleted += len(deleted)

					if len(new) == 0 && len(deleted) == 0 {
						continue
					}

					for _, f := range new {
						message += fmt.Sprintf("\n\t%s %s", createIndicator, f.Name())
					}

					for _, f := range deleted {
						message += fmt.Sprintf("\n\t%s %s", deleteIndicator, f.Name())
					}

					finalMessage += message + "\n\n"

					oldDirectories[index] = Directory{
						Path:  directory.Path,
						Files: files,
					}
				}

				if totalNew == 0 && totalDeleted == 0 {
					return
				}

				diff := fmt.Sprintf("+%d, -%d", totalNew, totalDeleted)

				_, err = client.R().
					SetBody(finalMessage).
					SetHeader("Title", "filetfy: "+diff).
					Post(ntfyServer + ntfyTopic)

				if err != nil {
					slog.Error("Erroring sending notification", "err", err)
				} else {
					slog.Info("Notification sent", "diff", diff)
				}
			},
		),
	)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(j.ID())

	s.Start()

	go func() {
		sig := <-sigChan
		fmt.Printf("\nReceived signal %v, shutting down...\n", sig)

		if err := s.Shutdown(); err != nil {
			fmt.Println("Error during scheduler shutdown:", err)
		}

		os.Exit(0)
	}()

	select {}
}

func difference(a, b []os.DirEntry) []os.DirEntry {
	bNames := make(map[string]struct{}, len(b))
	for _, entry := range b {
		bNames[entry.Name()] = struct{}{}
	}

	var diff []os.DirEntry
	for _, entryA := range a {
		if _, found := bNames[entryA.Name()]; !found {
			diff = append(diff, entryA)
		}
	}
	return diff
}

func exists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}
