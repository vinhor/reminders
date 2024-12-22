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
	Id       int       `json:"id"`
	WithTime bool      `json:"withTime"`
}

// TODO: Change the error messages to be more readable.

func openRemindersUnix() ([]Reminder, string) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}
	filePath := filepath.Join(homeDir, ".local/share/vinhor-reminders.json")
	file, err := os.ReadFile(filePath)
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
	return reminders, filePath
}
func saveFile(reminders []Reminder, filePath string) {
	newFile, err := json.MarshalIndent(reminders, "", "  ")
	if err != nil {
		panic(err)
	}
	err = os.WriteFile(filePath, newFile, 0644)
	if err != nil {
		panic(err)
	}
}

func addUnix() {
	reminders, filePath := openRemindersUnix()
	var newReminderData Reminder
	var date string
	var err error

	scanner := bufio.NewScanner(os.Stdin)

	fmt.Println("Please enter the reminder text:")
	if !scanner.Scan() {
		panic("Error reading reminder")
	}
	newReminderData.Text = scanner.Text()

	fmt.Println("Please enter the priority (1-3):")
	if !scanner.Scan() {
		panic("Error reading priority")
	}
	newReminderData.Priority, err = strconv.Atoi(scanner.Text())
	if err != nil || newReminderData.Priority < 1 || newReminderData.Priority > 3 {
		panic("Invalid priority")
	}

	var timeChoice string
	fmt.Println("Do you want to set hour and minute? y/n")
	if !scanner.Scan() {
		panic("Error reading answer")
	}
	timeChoice = scanner.Text()
	if timeChoice == "y" {
		newReminderData.WithTime = true
		fmt.Println("Please enter the due date in 24-hour format (2024 Dec 02 18:01):")
		if !scanner.Scan() {
			panic("Error reading date")
		}
		date = scanner.Text()
		newReminderData.Due, err = time.Parse("2006 Jan 02 15:04", date)
		if err != nil {
			panic("Invalid date")
		}
	} else if timeChoice == "n" {
		newReminderData.WithTime = false
		fmt.Println("Please enter the due date in 24-hour format (2024 Dec 02):")
		if !scanner.Scan() {
			panic("Error reading date")
		}
		date = scanner.Text()
		newReminderData.Due, err = time.Parse("2006 Jan 02", date)
		if err != nil {
			panic("Invalid date")
		}
	}

	newReminderData.Id = len(reminders) + 1

	newReminderData.Due = newReminderData.Due.In(time.Local)
	reminders = append(reminders, newReminderData)
	saveFile(reminders, filePath)
}

func listUnix() {
	reminders, _ := openRemindersUnix()
	if len(reminders) == 0 {
		fmt.Println("No reminders found.")
		return
	}
	fmt.Println("ID | Priority | Due                              | Text")
	for _, reminder := range reminders {
		var timeLayout string
		if reminder.WithTime {
			timeLayout = "2006 Jan 02, 15:04, -0700 offset"
		} else {
			timeLayout = "2006 Jan 02                     "
		}
		timeStamp := reminder.Due.Format(timeLayout)
		fmt.Printf("%d  | %d        | %s | %s\n", reminder.Id, reminder.Priority, timeStamp, reminder.Text)
	}
}

func rmAllUnix() {
	fmt.Println("Do you really want to remove all reminders? (y/N)")
	var answer string
	fmt.Scanln(&answer)
	if answer != "y" {
		fmt.Println("Aborted.")
		return
	}
	homeDir, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}
	filePath := filepath.Join(homeDir, ".local/share/vinhor-reminders.json")
	err = os.Remove(filePath)
	if err != nil {
		panic(err)
	}
	fmt.Println("All reminders removed.")
}

func removeUnix(id int) {
	id = id - 1
	reminders, filePath := openRemindersUnix()
	if id >= len(reminders) {
		fmt.Println("Error: Invalid reminder ID")
		return
	}
	reminders = append(reminders[:id], reminders[id+1:]...)
	saveFile(reminders, filePath)
}
