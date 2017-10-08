package main

import (
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"path/filepath"
	"time"
)

func generateID() string {
	t := time.Now().UTC()
	entropy := rand.New(rand.NewSource(t.UnixNano()))

	id := fmt.Sprintf("%X%X", t.Nanosecond(), entropy.Int())

	// Check that this ID doesn't exist yet
	filepath, err := filepath.Abs(uploadDir + id)
	if err != nil {
		// Failed to construct a filename, generate a new ID
		return generateID()
	}

	if _, err := os.Stat(filepath); os.IsNotExist(err) {
		return id
	}

	// Filename exists - generate a new ID
	return generateID()
}

func getPaste(id string) (string, error) {
	filepath, err := filepath.Abs(uploadDir + id)
	if err != nil {
		panic(err.Error())
	}

	if _, err := os.Stat(filepath); os.IsNotExist(err) {
		return "", err
	}

	bs, err := ioutil.ReadFile(filepath)
	if err != nil {
		return "", err
	}
	return string(bs), nil
}

func savePasteInFile(paste string) (string, error) {
	// Generate ID
	id := generateID()

	filepath, err := filepath.Abs(uploadDir + id)
	if err != nil {
		return "", err
	}

	file, err := os.Create(filepath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	file.WriteString(paste)
	return id, nil
}
