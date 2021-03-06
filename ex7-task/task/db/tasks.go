package db

import (
	"encoding/binary"
	"time"

	"github.com/boltdb/bolt"
)

// package level variables.
// not super recommended for web applications etc.
var taskBucket = []byte("tasks")
var db *bolt.DB

// Task describes a task in the to-do list
// it contains a Key and a Value
type Task struct {
	Key   int
	Value string
}

// Init function that gets called by the application
func Init(dbPath string) error {
	var err error
	// Open the my.db data file in current directory.
	// Will be created if it doesn't exist.
	db, err = bolt.Open(dbPath, 0600, &bolt.Options{Timeout: 1 * time.Second})

	if err != nil {
		return err
	}

	return db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists(taskBucket)
		return err
	})
}

// CreateTask creates a task and puts it in the bucket
func CreateTask(task string) (int, error) {
	var id int
	err := db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(taskBucket)
		id64, _ := b.NextSequence()
		id = int(id64)
		key := itob(id)
		return b.Put(key, []byte(task))
	})
	if err != nil {
		return -1, err
	}
	return id, nil
}

// AllTasks returns all the tasks from the db
func AllTasks() ([]Task, error) {
	var tasks []Task
	err := db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(taskBucket)
		c := b.Cursor()
		for k, v := c.First(); k != nil; k, v = c.Next() {
			tasks = append(tasks, Task{
				Key:   btoi(k),
				Value: string(v),
			})
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return tasks, nil
}

// DeleteTask deletes a task from the bucket
func DeleteTask(key int) error {
	return db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(taskBucket)
		return b.Delete(itob(key))
	})
}

// Convert integers to byte slices
func itob(v int) []byte {
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, uint64(v))
	return b
}

// Convert byte slices to integers
func btoi(b []byte) int {
	return int(binary.BigEndian.Uint64(b))
}
