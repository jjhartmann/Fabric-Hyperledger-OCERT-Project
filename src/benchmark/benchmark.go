/*
 * Used to benchmark ocert_scheme
 */

package main

import (
	"fmt"
	"ocert"
)

type DB struct {
	DB map[string][]byte
}

func (db *DB) GetState(key string) ([]byte, error) {
	val, exist := db.DB[key]
	if exist {
		return val, nil
	} else {
		return nil, fmt.Errorf("Failed to get state")
	}
}

func (db *DB) PutState(key string, value []byte) error {
	db.DB[key] = value
	return nil
}

func main() {
	db := new(DB)
	db.DB = make(map[string][]byte)

	// TODO delete
	var argPut []string
	argPut = append(argPut, "key1")
	argPut = append(argPut, "value1")

	ocert.Put(db, argPut)

	var argGet []string
	argGet = append(argGet, "key1")
	val, err := ocert.Get(db, argGet)

	if err == nil {
		fmt.Println(string(val))
	}

	// Benchmark starts here
}