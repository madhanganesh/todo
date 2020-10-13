package main

import (
	"encoding/json"
	"strings"
	"time"
)

// Todo struct
type Todo struct {
	ID     string    `json:"id"`
	Title  string    `json:"title"`
	Done   bool      `json:"done"`
	Due    time.Time `json:"due"`
	Tags   []string  `json:"tags"`
	Effort float32   `json:"duration"`
}

func (t Todo) id() []byte {
	return []byte(t.ID)
}

func (t Todo) data() []byte {
	d, err := json.Marshal(t)
	if err != nil {
		panic(err)
	}
	return d
}

func (t Todo) due() []byte {
	return []byte(t.Due.Format("2006-01-02"))
}

func (t Todo) datestr() string {
	return t.Due.Format("2 Jan 2006")
}

func (t Todo) tagsstr() string {
	return strings.Join(t.Tags, ",")
}

func makeTodo(data []byte) (Todo, error) {
	var todo Todo
	err := json.Unmarshal(data, &todo)
	return todo, err
}
