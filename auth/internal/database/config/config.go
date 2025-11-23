package config

import (
	"context"
	"log"
	"sync"
)

var mu sync.Mutex



func InitDB(ctx context.Context) error {
	mu.Lock()
	defer mu.Unlock()

	if db != nil {
		return nil
	}

	pool, err := NewPool(ctx)

	if err != nil {
		return err
	}

	db = pool
	return nil

}

func Get() *Pool {
	return db
}

func Close() {
	if db != nil && db.Pool != nil {
		db.Close()
		db = nil
		log.Println("Database connection closed!")
	}
}

//for testing purpose
func Reset(){
	db = nil
}