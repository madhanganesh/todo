package main

// UserKey key
var UserKey = "SELF"

// PendingKey key
var PendingKey = []byte("pending")

// ListingKey key
var ListingKey = []byte("listing")

// OpType type
type OpType string

const (
	// ShowHelp option
	ShowHelp OpType = "help"
	// AddTodo option
	AddTodo = "add"
	// ListTodos option
	ListTodos = "list"
	// ShowTodoDetail option
	ShowTodoDetail = "detail"
	// SetDone option
	SetDone = "setDone"
	// SetDue option
	SetDue = "setDue"
	// SetTags option
	SetTags = "setTags"
	// SetEffort option
	SetEffort = "setEffort"
	// DeleteTodo option
	DeleteTodo = "delete"
	// UpdateTodo option
	UpdateTodo = "update"
)

// Operation type
type Operation func(Opts, *TodoRepo) error

// OperationMap mapping
var OperationMap = map[OpType]Operation{
	AddTodo:        addTodo,
	ListTodos:      listTodos,
	ShowTodoDetail: showTodo,
	SetDone:        setDone,
	SetDue:         setDue,
	SetTags:        setTags,
	SetEffort:      setEffort,
	DeleteTodo:     deleteTodo,
	UpdateTodo:     updateTodo,
}
