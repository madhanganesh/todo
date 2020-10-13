package main

import (
	"log"
	"os"

	bolt "go.etcd.io/bbolt"
)

func main() {
	db, err := bolt.Open(getDbPath(), 0600, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	repo := &TodoRepo{}
	repo.Init(db)

	opts, err := getOpts(os.Args)
	if err != nil {
		log.Fatal(err)
	}

	operation := OperationMap[opts.option]
	if operation == nil {
		log.Fatal("Unknown option: ", opts.option)
	}

	err = operation(opts, repo)
	if err != nil {
		log.Fatal(err)
	}
}
