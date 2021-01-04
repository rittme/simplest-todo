package model

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
)

// Defines the type for the entries list
type Entries struct {
	Entries []Entry `json:"entries"`
}

// Defines the type for an entry
type Entry struct {
	Label  string `json:"label"`
	Done   bool `json:"done"`
}

// Defines the type for an entry with ID (line number)
type IDEntry struct {
	ID    int
	Entry Entry
}

// Defines the path for the list file
const FilePath = "list.json"

func ReadFile() (entries Entries, err error) {
	jsonFile, err := os.Open(FilePath)
	if err != nil {
		return  
	}
	defer jsonFile.Close()

	byteValue, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		return
	}
	err = json.Unmarshal(byteValue, &entries)
	if err != nil {
		return
	}
	return
}

func WriteToFile(entries Entries) (err error) {
	jsonFile, err := os.OpenFile(FilePath, os.O_RDWR, os.ModePerm)
	if err != nil {
		return
	}
	defer jsonFile.Close()

	err = jsonFile.Truncate(0)
	if err != nil {
		return
	}
	_, err = jsonFile.Seek(0, 0)
	if err != nil {
		return
	}

	err = json.NewEncoder(jsonFile).Encode(entries)
	return
}

func GetRandom(includeDone string) (entry IDEntry, err error) {
	entries, err := ReadFile()	
	if err != nil {
		fmt.Println(err)
		return
	}

	filteredEntries := []IDEntry{}
	
	for i := range entries.Entries {
			if includeDone == "true" || !entries.Entries[i].Done {
					filteredEntries = append(filteredEntries, IDEntry{i, entries.Entries[i]})
			}
	}
	if len(filteredEntries) == 0 {
		return
	}
	if len (filteredEntries) == 1 {
		entry = filteredEntries[0]
	}

	random := rand.Intn(len(filteredEntries)-1)
	entry = filteredEntries[random]
	return
}

func CreateNew(label string) (entries Entries, err error) {
	entries, err = ReadFile()
	if err != nil {
		fmt.Println(err)
		return
	}

	entries.Entries = append(entries.Entries, Entry{string(label), false})

	jsonFile, err := os.OpenFile(FilePath, os.O_RDWR, os.ModePerm)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer jsonFile.Close()
	
	err = WriteToFile(entries)

	return
}

func ToggleDone(id int) (entries Entries, err error) {
	entries, err = ReadFile()
	if err != nil {
		return
	}

	if id < 0 || id >= len(entries.Entries) {
		err = fmt.Errorf("Index out of bounds: %d", id)
		return
	}

	entries.Entries[id].Done = !entries.Entries[id].Done

	err = WriteToFile(entries)

	return
}

func DeleteEntry(id int) (entries Entries, err error) {
	entries, err = ReadFile()
	if err != nil {
		return
	}

	if id < 0 || id >= len(entries.Entries) {
		err = fmt.Errorf("Index out of bounds: %d", id)
		return
	}

  entries.Entries = append(entries.Entries[:id], entries.Entries[id+1:]...)
	
	err = WriteToFile(entries)

	return
}