package main

import (
	"container/list"
)

type HashTable struct {
	count    int16
	capacity int16
	values   []list.List // Store tuples of key and value
}

type KeyValuePair struct {
	key   string
	value string
}

// Source: https://stackoverflow.com/questions/7666509/hash-function-for-string
func get_hash(hash_table *HashTable, input string) int16 {
	hash := int64(5281)

	for char := range input {
		hash = ((hash << 5) + hash) + int64(char)
	}

	return int16(hash % int64(hash_table.capacity))
}

func create_table() *HashTable {
	return create_table_capacity(10) // Default size
}

func create_table_capacity(capacity int16) *HashTable {
	hash_table := HashTable{}
	hash_table.values = make([]list.List, capacity)
	hash_table.capacity = capacity
	return &hash_table
}

func insert_value_table(hash_table *HashTable, key string, value string) {
	hash := get_hash(hash_table, key)

	// Create new linked list
	if hash_table.values[hash].Len() == 0 {
		hash_table.values[hash] = *list.New()
	}

	// Insert element
	hash_table.values[hash].PushFront(KeyValuePair{key, value})
	hash_table.count += 1

	// Resize hash table (double the size)
	if float32(hash_table.count) > 0.7*float32(hash_table.capacity) {
		old_values := get_all_table(hash_table)
		*hash_table = *create_table_capacity(hash_table.capacity * 2)

		// Insert all old values
		for _, key_value := range old_values {
			insert_value_table(hash_table, key_value.key, key_value.value)
		}
	}
}

func delete_value_table(hash_table *HashTable, key string) {
	hash := get_hash(hash_table, key)
	list := &hash_table.values[hash]

	// Find value in list
	for element := list.Front(); element != nil; element = element.Next() {
		if element.Value.(KeyValuePair).key == key {
			list.Remove(element)
			hash_table.count -= 1
			return
		}
	}
}

func get_value_table(hash_table *HashTable, key string) (string, bool) {
	hash := get_hash(hash_table, key)
	list := hash_table.values[hash]

	// Find value in list
	for element := list.Front(); element != nil; element = element.Next() {
		if element.Value.(KeyValuePair).key == key {
			return element.Value.(KeyValuePair).value, true
		}
	}

	return "", false
}

func get_all_table(hash_table *HashTable) []KeyValuePair {
	var all []KeyValuePair

	for _, values := range hash_table.values {
		for element := values.Front(); element != nil; element = element.Next() {
			if element.Value != nil {
				key_value := element.Value.(KeyValuePair)
				all = append(all, key_value)
			}
		}
	}

	return all
}
