package main

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/rs/xid"
)

func addTodo(opts Opts, repo *TodoRepo) error {
	todo := Todo{
		ID:     xid.New().String(),
		Title:  opts.params["title"].(string),
		Due:    opts.params["due"].(time.Time),
		Tags:   opts.params["tags"].([]string),
		Done:   opts.params["done"].(bool),
		Effort: opts.params["effort"].(float32),
	}
	return repo.CreateTodo(UserKey, todo)
}

func listTodos(opts Opts, repo *TodoRepo) error {
	todos, err := getTodosByFilter(opts, repo)
	if err != nil {
		return err
	}

	listMapping := printTodos(opts, todos)
	repo.SetListMapping(listMapping)

	return nil
}

func showTodo(opts Opts, repo *TodoRepo) error {
	id := idFromOpts(opts, repo)
	todo, err := repo.GetTodo(UserKey, id)
	if err != nil {
		return err
	}
	printTodo(todo)
	return nil
}

func setDone(opts Opts, repo *TodoRepo) error {
	id := idFromOpts(opts, repo)
	status := opts.params["done"].(bool)
	return repo.SetTodoDone(UserKey, id, status)
}

func setDue(opts Opts, repo *TodoRepo) error {
	id := idFromOpts(opts, repo)
	due := opts.params["due"].(time.Time)
	return repo.SetTodoDue(UserKey, id, due)
}

func setTags(opts Opts, repo *TodoRepo) error {
	id := idFromOpts(opts, repo)
	tags := opts.params["tags"].([]string)
	return repo.SetTodoTags(UserKey, id, tags)
}

func setEffort(opts Opts, repo *TodoRepo) error {
	id := idFromOpts(opts, repo)
	effort := opts.params["effort"].(float32)
	return repo.SetTodoEffort(UserKey, id, effort)
}

func deleteTodo(opts Opts, repo *TodoRepo) error {
	id := idFromOpts(opts, repo)
	return repo.DeleteTodo(UserKey, id)
}

func updateTodo(opts Opts, repo *TodoRepo) error {
	id := idFromOpts(opts, repo)
	todo, err := repo.GetTodo(UserKey, id)
	if err != nil {
		return err
	}

	if donei, present := opts.params["done"]; present {
		todo.Done = donei.(bool)
	}

	if duei, present := opts.params["due"]; present {
		todo.Due = duei.(time.Time)
	}

	if efforti, present := opts.params["effort"]; present {
		todo.Effort = efforti.(float32)
	}

	if tagsi, present := opts.params["tags"]; present {
		todo.Tags = tagsi.([]string)
	}

	return repo.UpdateTodo(UserKey, id, todo)
}

func getTodosByFilter(opts Opts, repo *TodoRepo) ([]Todo, error) {
	filter := opts.params["type"].(string)
	switch filter {
	case "pending":
		return repo.GetPendingTodos(UserKey)
	case "bydate":
		return repo.GetTodosByDate(UserKey, opts.params["date"].(time.Time))
	default:
		return nil, fmt.Errorf("Unknow option for type %s", filter)
	}
}

func printTodos(opts Opts, todos []Todo) map[string]string {
	heading := getHeadingForPrint(opts)

	fmt.Printf("\n%s\n", strings.Title(heading))
	for range heading {
		fmt.Print("-")
	}
	fmt.Println()
	mapping := map[string]string{}
	totalEffort := float32(0.0)
	completedTodos := 0

	for i, todo := range todos {
		s := strconv.Itoa(i + 1)
		mapping[s] = todo.ID
		if heading == "Pending" {
			fmt.Printf("%s. %s - %s\n", s, todo.Due.Format("02 Jan"), todo.Title)
			continue
		}

		if todo.Done {
			fmt.Printf("%s. [X] (%0.1f) %s\n", s, todo.Effort, todo.Title)
		} else {
			fmt.Printf("%s. [ ] (%0.1f) %s\n", s, todo.Effort, todo.Title)
		}

		totalEffort += todo.Effort
		if todo.Done {
			completedTodos++
		}
	}
	fmt.Println()

	if heading != "Pending" {
		fmt.Printf("%d / %d Todos pending\n", len(todos)-completedTodos, len(todos))
		fmt.Printf("%.1f hours of total effort\n\n", totalEffort)
	}

	return mapping
}

func getHeadingForPrint(opts Opts) string {
	filter := opts.params["type"].(string)
	switch filter {
	case "pending":
		return "Pending"
	case "bydate":
		date := opts.params["date"].(time.Time)
		return date.Format("2006-01-02")
	default:
		return "Unknown"
	}
}

func idFromOpts(opts Opts, repo *TodoRepo) string {
	key := opts.params["id"].(string)
	listMapping := repo.GetListMapping()
	id := listMapping[key]

	return id
}
