package main

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// Opts struct
type Opts struct {
	option OpType
	params map[string]interface{}
}

func getOpts(args []string) (Opts, error) {
	var opts Opts
	opts.params = map[string]interface{}{}

	if len(args) == 1 {
		opts.option = ListTodos
		opts.params["type"] = "pending"
		return opts, nil
	}

	if len(args) == 2 {
		match, _ := regexp.MatchString("^[-/]?h(elp)?", args[1])
		if match {
			opts.option = ShowHelp
			return opts, nil
		}

		match, _ = regexp.MatchString(`^(today|tomorrow|yesterday|\d\d\d\d-\d\d-\d\d)$`, args[1])
		if match {
			opts.option = ListTodos
			opts.params["type"] = "bydate"
			opts.params["date"] = toDate(args[1])
			return opts, nil
		}

		if args[1] == "thisweek" || args[1] == "lastweek" || args[1] == "nextweek" {
			panic("not implemented")
		}

		opts.option = ShowTodoDetail
		opts.params["id"] = args[1]
	}

	if len(args) == 3 {
		opts.params["id"] = args[1]

		if args[2] == "delete" {
			opts.option = DeleteTodo
			return opts, nil
		}

		if args[2] == "done" {
			opts.option = SetDone
			opts.params["done"] = true
			return opts, nil
		}

		if args[2] == "pending" {
			opts.option = SetDone
			opts.params["done"] = false
			return opts, nil
		}

		match, _ := regexp.MatchString(`^(today|tomorrow|yesterday|\d\d\d\d-\d\d-\d\d)$`, args[2])
		if match {
			opts.option = SetDue
			opts.params["due"] = toDate(args[2])
			return opts, nil
		}

		if strings.HasPrefix(args[2], "#") {
			tags := []string{}
			tagsh := strings.Split(args[2], ",")
			for _, tag := range tagsh {
				if string(tag[0]) == "#" {
					tag = tag[1:]
				}
				tags = append(tags, tag)
			}

			opts.option = SetTags
			opts.params["tags"] = tags
			return opts, nil
		}

		effort, err := strconv.ParseFloat(args[2], 32)
		if err == nil {
			opts.option = SetEffort
			opts.params["effort"] = float32(effort)
			return opts, nil
		}

		return Opts{}, fmt.Errorf("Unexpected Arguments. To Add Todo minimum 3 words are required")
	}

	if len(args) > 3 {
		if isIndex, index := isIndex(args[1]); isIndex {
			params := args[1:]
			opts.option = UpdateTodo
			opts.params["id"] = index
			fillInParams(params, &opts)
			return opts, nil
		}
	}

	if len(args) > 3 {
		opts.option = AddTodo
		temp := []string{}
		params := []string{}

		for i, t := range args[1:] {
			if t == "-p" {
				params = args[i+2:]
				break
			}

			if i == 0 {
				t = strings.Title(t)
			}
			temp = append(temp, t)
		}

		opts.params["title"] = strings.Join(temp, " ")
		opts.params["due"] = today()
		opts.params["done"] = false
		opts.params["effort"] = float32(0.0)
		opts.params["tags"] = []string{}
		fillInParams(params, &opts)
		return opts, nil
	}

	return opts, nil
}

func fillInParams(params []string, opts *Opts) {
	for _, param := range params {
		fillInParam(param, opts)
	}
}

func fillInParam(param string, opts *Opts) {
	if isDone, done := isDone(param); isDone {
		opts.params["done"] = done
		return

	}

	if isDate, date := isDate(param); isDate {
		opts.params["due"] = date
		return
	}

	if isEffort, effort := isEffort(param); isEffort {
		opts.params["effort"] = float32(effort)
		return
	}

	if isTags, tags := isTags(param); isTags {
		opts.params["tags"] = tags
		return
	}
}

func isIndex(param string) (bool, string) {
	_, err := strconv.Atoi(param)
	return err == nil, param
}

func isDone(param string) (bool, bool) {
	if param == "done" {
		return true, true
	}

	if param == "pending" {
		return true, false
	}

	return false, false
}

func isDate(param string) (bool, time.Time) {
	match, _ := regexp.MatchString(`^(today|tomorrow|yesterday|\d\d\d\d-\d\d-\d\d)$`, param)
	if match {
		return true, toDate(param)
	}
	return false, time.Now()
}

func isEffort(param string) (bool, float32) {
	effort, err := strconv.ParseFloat(param, 32)
	if err == nil {
		return true, float32(effort)
	}

	return false, 0.0
}

func isTags(param string) (bool, []string) {
	if strings.HasPrefix(param, "#") {
		tags := []string{}
		tagsh := strings.Split(param, ",")
		for _, tag := range tagsh {
			if string(tag[0]) == "#" {
				tag = tag[1:]
			}
			tags = append(tags, tag)
		}

		return true, tags
	}

	return false, nil
}
