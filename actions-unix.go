package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

type Reminder struct {
	Text     string    `json:"text"`
	Priority int       `json:"priority"`
	Due      time.Time `json:"due"`
}

func addUnix() {
	homeDir, _ := os.UserHomeDir()
	filesPath := filepath.Join(homeDir, ".local/share/vinhor-reminders.json")
	file, err := os.ReadFile(filesPath)
	if err != nil && !os.IsNotExist(err) {
		panic(err)
	}
	var reminders []Reminder
	if file != nil {
		err = json.Unmarshal(file, &reminders)
		if err != nil {
			panic(err)
		}
	}
	var newReminderData Reminder
	var date string

	scanner := bufio.NewScanner(os.Stdin)

	fmt.Println("Please enter the reminder text:")
	if scanner.Scan() {
		newReminderData.Text = scanner.Text()
	} else {
		panic("Error reading reminder")
	}

	fmt.Println("Please enter the priority (1-3):")
	if scanner.Scan() {
		newReminderData.Priority, err = strconv.Atoi(scanner.Text())
		if err != nil || newReminderData.Priority < 1 || newReminderData.Priority > 3 {
			panic("Invalid priority")
		}
	} else {
		panic("Error reading priority")
	}

	var timeChoice string
	fmt.Println("Do you want to set hour and minute? y/n")
	if scanner.Scan() {
		timeChoice = scanner.Text()
	} else {
		panic("Error reading answer")
	}
	if timeChoice == "y" {
		fmt.Println("Please enter the due date in 24-hour format (2024 Dec 02 18:01):")
		if scanner.Scan() {
			date = scanner.Text()
			newReminderData.Due, err = time.Parse("2006 Jan 02 15:04", date)
			if err != nil {
				panic("Invalid date")
			}
		} else {
			panic("Error reading date")
		}
	} else if timeChoice == "n" {
		fmt.Println("Please enter the due date in 24-hour format (2024 Dec 02):")
		if scanner.Scan() {
			date = scanner.Text()
			newReminderData.Due, err = time.Parse("2006 Jan 02", date)
			if err != nil {
				panic("Invalid date")
			}
		} else {
			panic("Error reading date")
		}
	}

	newReminderData.Due = newReminderData.Due.In(time.Local)
	reminders = append(reminders, newReminderData)
	newFile, err := json.MarshalIndent(reminders, "", "  ")
	if err != nil {
		panic(err)
	}
	err = os.WriteFile(filesPath, newFile, 0644)
	if err != nil {
		panic(err)
	}
}
