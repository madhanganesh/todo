package main

import (
	"fmt"
	"time"

	bolt "go.etcd.io/bbolt"
)

// TodoRepo struct
type TodoRepo struct {
	db *bolt.DB
}

// Init method
func (r *TodoRepo) Init(db *bolt.DB) {
	r.db = db
}

// Create method
func (r *TodoRepo) CreateTodo(userID string, t Todo) error {
	err := r.db.Update(func(tx *bolt.Tx) error {
		userBucket, err := tx.CreateBucketIfNotExists([]byte(userID))
		if err != nil {
			return err
		}

		// Add todo to user bucket
		err = userBucket.Put(t.id(), t.data())
		if err != nil {
			return err
		}

		// Add to date bucket for easy search on dates
		dateBucket, err := userBucket.CreateBucketIfNotExists(t.due())
		if err != nil {
			return err
		}
		err = dateBucket.Put(t.id(), t.id())
		if err != nil {
			return err
		}

		// Add to pending bucket
		if !t.Done {
			pendingBucket, err := userBucket.CreateBucketIfNotExists(PendingKey)
			if err != nil {
				return err
			}

			err = pendingBucket.Put(t.id(), t.id())
			if err != nil {
				return err
			}
		}

		return nil
	})

	return err
}

// GetPendingTodos method
func (r *TodoRepo) GetPendingTodos(userID string) ([]Todo, error) {
	todos := []Todo{}

	err := r.db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(userID))
		if bucket == nil {
			return nil
		}

		pendingBucket := bucket.Bucket(PendingKey)
		if pendingBucket == nil {
			return nil
		}

		c := pendingBucket.Cursor()
		for id, _ := c.First(); id != nil; id, _ = c.Next() {
			data := bucket.Get(id)
			if data != nil {
				todo, err := makeTodo(data)
				if err != nil {
					return err
				}

				todos = append(todos, todo)
			}
		}

		return nil
	})

	return todos, err
}

// GetTodosByDate method
func (r *TodoRepo) GetTodosByDate(userID string, date time.Time) ([]Todo, error) {
	todos := []Todo{}

	err := r.db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(userID))
		if bucket == nil {
			return nil
		}

		dateBucket := bucket.Bucket([]byte(date.Format("2006-01-02")))
		if dateBucket == nil {
			return nil
		}

		c := dateBucket.Cursor()
		for id, _ := c.First(); id != nil; id, _ = c.Next() {
			data := bucket.Get(id)
			if data != nil {
				todo, err := makeTodo(data)
				if err != nil {
					return err
				}

				todos = append(todos, todo)
			}
		}

		return nil
	})

	return todos, err
}

// GetTodo method
func (r *TodoRepo) GetTodo(userID string, id string) (Todo, error) {
	var todo Todo

	err := r.db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(userID))
		if bucket == nil {
			return nil
		}

		data := bucket.Get([]byte(id))
		if data == nil {
			return fmt.Errorf("Could find tod with ID: %s", id)
		}

		var err error
		todo, err = makeTodo(data)
		if err != nil {
			return err
		}

		return nil
	})

	return todo, err
}

// SetTodoDone method
func (r *TodoRepo) SetTodoDone(userID string, todoID string, status bool) error {
	err := r.db.Update(func(tx *bolt.Tx) error {
		userBucket := tx.Bucket([]byte(userID))
		if userBucket == nil {
			return fmt.Errorf("Unable to find user bucket for %s", userID)
		}

		data := userBucket.Get([]byte(todoID))
		if data == nil {
			return fmt.Errorf("No todo found for ID: %s", todoID)
		}

		todo, err := makeTodo(data)
		if err != nil {
			return err
		}

		todo.Done = status
		todo.Effort = 1.0
		err = userBucket.Put([]byte(todo.ID), todo.data())
		if err != nil {
			return err
		}

		if !todo.Done {
			pendingBucket, err := userBucket.CreateBucketIfNotExists(PendingKey)
			if err != nil {
				return err
			}

			err = pendingBucket.Put(todo.id(), nil)
			if err != nil {
				return err
			}
		}

		if todo.Done {
			pendingBucket := userBucket.Bucket(PendingKey)
			if pendingBucket != nil && pendingBucket.Get(todo.id()) != nil {
				pendingBucket.Delete(todo.id())
			}
		}

		return nil
	})

	return err
}

// SetTodoEffort method
func (r *TodoRepo) SetTodoEffort(userID string, todoID string, effort float32) error {
	err := r.db.Update(func(tx *bolt.Tx) error {
		userBucket := tx.Bucket([]byte(userID))
		if userBucket == nil {
			return fmt.Errorf("Unable to find user bucket for %s", userID)
		}

		data := userBucket.Get([]byte(todoID))
		if data == nil {
			return fmt.Errorf("No todo found for ID: %s", todoID)
		}

		todo, err := makeTodo(data)
		if err != nil {
			return err
		}

		// Update Todo in user bucket
		todo.Effort = effort
		err = userBucket.Put([]byte(todo.ID), todo.data())
		return err
	})

	return err
}

// SetTodoDue method
func (r *TodoRepo) SetTodoDue(userID string, todoID string, due time.Time) error {
	err := r.db.Update(func(tx *bolt.Tx) error {
		userBucket := tx.Bucket([]byte(userID))
		if userBucket == nil {
			return fmt.Errorf("Unable to find user bucket for %s", userID)
		}

		data := userBucket.Get([]byte(todoID))
		if data == nil {
			return fmt.Errorf("No todo found for ID: %s", todoID)
		}

		todo, err := makeTodo(data)
		if err != nil {
			return err
		}

		// Remove from the curret date bucket
		currentDateBucket := userBucket.Bucket(todo.due())
		if currentDateBucket != nil {
			currentDateBucket.Delete(todo.id())
		}

		// Update Todo in user bucket
		todo.Due = due
		err = userBucket.Put([]byte(todo.ID), todo.data())
		if err != nil {
			return err
		}

		// Add to new date bucket
		newDateBucket, err := userBucket.CreateBucketIfNotExists(todo.due())
		if err != nil {
			return err
		}
		err = newDateBucket.Put(todo.id(), todo.id())

		return err
	})

	return err
}

// SetTodoTags method
func (r *TodoRepo) SetTodoTags(userID string, todoID string, tags []string) error {
	err := r.db.Update(func(tx *bolt.Tx) error {
		userBucket := tx.Bucket([]byte(userID))
		if userBucket == nil {
			return fmt.Errorf("Unable to find user bucket for %s", userID)
		}

		data := userBucket.Get([]byte(todoID))
		if data == nil {
			return fmt.Errorf("No todo found for ID: %s", todoID)
		}

		todo, err := makeTodo(data)
		if err != nil {
			return err
		}

		todo.Tags = tags
		err = userBucket.Put(todo.id(), todo.data())
		if err != nil {
			return err
		}
		return err
	})

	return err
}

// DeleteTodo method
func (r *TodoRepo) DeleteTodo(userID string, todoID string) error {
	err := r.db.Update(func(tx *bolt.Tx) error {
		userBucket := tx.Bucket([]byte(userID))
		if userBucket == nil {
			return fmt.Errorf("Unable to find user bucket for %s", userID)
		}

		// First retrieve the Todo
		data := userBucket.Get([]byte(todoID))
		if data == nil {
			return fmt.Errorf("No todo found for ID: %s", todoID)
		}
		todo, err := makeTodo(data)
		if err != nil {
			return err
		}

		// Delete from the date bucket
		dateBucket := userBucket.Bucket(todo.due())
		if dateBucket == nil {
			return fmt.Errorf("No date bucket found for date %v", todo.datestr())
		}
		err = dateBucket.Delete(todo.id())
		if err != nil {
			return err
		}

		// Delete from teh pending bucket, if present
		pendingBucket := userBucket.Bucket(PendingKey)
		if pendingBucket != nil {
			err = pendingBucket.Delete([]byte(todoID))
			if err != nil {
				return err
			}
		}

		// Finally delete from user bucket
		err = userBucket.Delete([]byte(todoID))
		return err
	})

	return err
}

// UpdateTodo method
func (r *TodoRepo) UpdateTodo(userID string, todoID string, todo Todo) error {
	err := r.DeleteTodo(userID, todoID)
	if err != nil {
		return err
	}

	return r.CreateTodo(userID, todo)
}

// SetListMapping is a method to persist the todo list index that
// is dieplayed to teh user to ID. This mapping is persisted after each
// display of the Todos. Users will provide just the index like 1, 2, ..
// for subsequent operations and the mapping will be retrieved to fetch the ID
func (r *TodoRepo) SetListMapping(mapping map[string]string) {
	r.db.Update(func(t *bolt.Tx) error {
		bucket, err := t.CreateBucketIfNotExists(ListingKey)
		if err != nil {
			return err
		}

		for k, v := range mapping {
			bucket.Put([]byte(k), []byte(v))
		}

		return nil
	})
}

// GetListMapping is a method to retrive back the mapping from display
// index number to Todo ID
func (r *TodoRepo) GetListMapping() map[string]string {
	listMapping := map[string]string{}

	r.db.View(func(t *bolt.Tx) error {
		bucket := t.Bucket(ListingKey)
		if bucket != nil {
			c := bucket.Cursor()
			for k, v := c.First(); k != nil; k, v = c.Next() {
				listMapping[string(k)] = string(v)
			}
		}

		return nil
	})

	return listMapping
}
