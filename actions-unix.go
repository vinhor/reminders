package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/fatih/color"
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
		color.Red("Error getting home directory! Do you have $HOME or %USERPROFILE% set?")
		panic(err)
	}
	filePath := filepath.Join(homeDir, ".local/share/vinhor-reminders.json")
	file, err := os.ReadFile(filePath)
	if err != nil && !os.IsNotExist(err) {
		color.Red("Error opening reminders file!")
		panic(err)
	}
	var reminders []Reminder
	if file != nil {
		err = json.Unmarshal(file, &reminders)
		if err != nil {
			color.Red("Error parsing reminders file!")
			panic(err)
		}
	}
	return reminders, filePath
}
func saveFile(reminders []Reminder, filePath string) {
	newFile, err := json.MarshalIndent(reminders, "", "  ")
	if err != nil {
		color.Red("Error marshaling reminders!")
		panic(err)
	}
	err = os.WriteFile(filePath, newFile, 0644)
	if err != nil {
		color.Red("Error saving to file!")
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
		color.Red("Error reading reminder text")
		os.Exit(1)
	}
	newReminderData.Text = scanner.Text()

	fmt.Println("Please enter the priority (1-3):")
	if !scanner.Scan() {
		color.Red("Error reading priority")
		os.Exit(1)
	}
	newReminderData.Priority, err = strconv.Atoi(scanner.Text())
	if newReminderData.Priority < 1 || newReminderData.Priority > 3 {
		color.Red("Invalid priority")
		os.Exit(1)
	}
	if err != nil {
		color.Red("Error parsing priority")
		panic(err)
	}

	var timeChoice string
	fmt.Println("Do you want to set hour and minute? (y/n)")
	if !scanner.Scan() {
		color.Red("Error reading answer")
		os.Exit(1)
	}
	timeChoice = scanner.Text()
	if timeChoice == "y" {
		newReminderData.WithTime = true
		fmt.Println("Please enter the due date in 24-hour format (2024 Dec 02 18:01):")
		if !scanner.Scan() {
			color.Red("Error reading date")
			os.Exit(1)
		}
		date = scanner.Text()
		newReminderData.Due, err = time.Parse("2006 Jan 02 15:04", date)
		if err != nil {
			color.Red("Invalid date")
			panic(err)
		}
	} else if timeChoice == "n" {
		newReminderData.WithTime = false
		fmt.Println("Please enter the due date in 24-hour format (2024 Dec 02):")
		if !scanner.Scan() {
			color.Red("Error reading date")
			os.Exit(1)
		}
		date = scanner.Text()
		newReminderData.Due, err = time.Parse("2006 Jan 02", date)
		if err != nil {
			color.Red("Invalid date")
			panic(err)
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
		color.Yellow("No reminders found.")
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
	color.Red("Do you really want to remove all reminders? (y/N)")
	var answer string
	fmt.Scanln(&answer)
	if answer != "y" {
		color.Yellow("Aborted.")
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
		color.Red("Error: Invalid reminder ID")
		os.Exit(1)
	}
	reminders = append(reminders[:id], reminders[id+1:]...)
	saveFile(reminders, filePath)
}

func editUnix(id int) {
	id = id - 1
	reminders, filePath := openRemindersUnix()
	if id >= len(reminders) {
		color.Red("Error: Invalid reminder ID")
		os.Exit(1)
	}

	reminder := reminders[id]
	fmt.Println("ID | Priority | Due                              | Text")
	var timeLayout string
	if reminder.WithTime {
		timeLayout = "2006 Jan 02, 15:04, -0700 offset"
	} else {
		timeLayout = "2006 Jan 02                     "
	}
	timeStamp := reminder.Due.Format(timeLayout)
	fmt.Printf("%d  | %d        | %s | %s\n", reminder.Id, reminder.Priority, timeStamp, reminder.Text)

	fmt.Println("What do you want to change? ([p]riority/[d]ue/[t]ext)")
	var choice string
	fmt.Scanln(&choice)
	if choice != "p" && choice != "d" && choice != "t" {
		color.Red("Invalid choice")
		os.Exit(1)
	}

	switch choice {
	case "p":
		var newPriority int
		fmt.Println("Please enter the new priority (1-3):")
		fmt.Scan(&newPriority)
		if newPriority < 1 || newPriority > 3 {
			color.Red("Invalid priority")
			os.Exit(1)
		}
		reminder.Priority = newPriority
	case "d":
		var timeChoice string
		scanner := bufio.NewScanner(os.Stdin)
		fmt.Println("Do you want to set hour and minute? (y/n)")
		if !scanner.Scan() {
			color.Red("Error reading answer")
			os.Exit(1)
		}
		timeChoice = scanner.Text()

		if timeChoice == "y" {
			reminder.WithTime = true
			fmt.Println("Please enter the due date in 24-hour format (2024 Dec 02 18:01):")
			if !scanner.Scan() {
				color.Red("Error reading date")
				os.Exit(1)
			}
			date := scanner.Text()
			var err error
			reminder.Due, err = time.Parse("2006 Jan 02 15:04", date)
			if err != nil {
				color.Red("Invalid date")
				panic(err)
			}
			reminder.Due = reminder.Due.In(time.Local)
		} else if timeChoice == "n" {
			reminder.WithTime = false
			fmt.Println("Please enter the due date in 24-hour format (2024 Dec 02):")
			if !scanner.Scan() {
				color.Red("Error reading date")
				os.Exit(1)
			}
			date := scanner.Text()
			var err error
			fmt.Scanln(&date)
			reminder.Due, err = time.Parse("2006 Jan 02", date)
			if err != nil {
				color.Red("Invalid date")
				panic(err)
			}
			reminder.Due = reminder.Due.In(time.Local)
		} else {
			color.Red("Invalid choice")
			os.Exit(1)
		}
	case "t":
		fmt.Println("Enter the new reminder text:")
		fmt.Scanln(&reminder.Text)
	}

	reminders = append(reminders[:id], reminder)
	saveFile(reminders, filePath)
	fmt.Println("Reminder updated.")
}
