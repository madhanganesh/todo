package main

import (
	"fmt"
	"log"
	"os"
	"runtime"
	"time"
)

func startOfDay(t time.Time) time.Time {
	year, month, day := t.Date()
	return time.Date(year, month, day, 0, 0, 0, 0, t.Location())
}

func today() time.Time {
	return startOfDay(time.Now())
}

func tomorrow() time.Time {
	return startOfDay(time.Now().Add(24 * time.Hour))
}

func toDate(str string) time.Time {
	switch str {
	case "today":
		return today()
	case "yesterday":
		return startOfDay(time.Now().Add(-1 * 24 * time.Hour))
	case "tomorrow":
		return startOfDay(time.Now().Add(24 * time.Hour))
	default:
		zone, _ := time.Now().Zone()
		parsed, err := time.Parse("2006-01-02 MST", str+" "+zone)
		if err != nil {
			panic("unknown value for toDate " + str)
		}
		return startOfDay(parsed)
	}
}

func printTodo(todo Todo) {
	fmt.Printf(`
Task   : %s
Due    : %s
Done   : %v
Effort : %.1f hours
Tags   : %s

`, todo.Title, todo.datestr(), todo.Done, todo.Effort, todo.tagsstr())
}

func userHomeDir() string {
	if runtime.GOOS == "windows" {
		home := os.Getenv("HOMEDRIVE") + os.Getenv("HOMEPATH")
		if home == "" {
			home = os.Getenv("USERPROFILE")
		}
		return home
	} else if runtime.GOOS == "linux" {
		home := os.Getenv("XDG_CONFIG_HOME")
		if home != "" {
			return home
		}
	}
	return os.Getenv("HOME")
}

func getDbPath() string {
	homeDir := userHomeDir()
	dbDir := homeDir + string(os.PathSeparator) + "todo"
	if _, err := os.Stat(dbDir); os.IsNotExist(err) {
		err = os.Mkdir(dbDir, 0700)
		if err != nil {
			log.Fatal(err)
		}
	}

	dbFile := dbDir + string(os.PathSeparator) + "todo.db"
	return dbFile
}
