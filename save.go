package main

import (
	"encoding/json"
	"fmt"
	"os"
)

// This file allows the application to load/save from json files

func (m model) SaveToFile() {
	var f *os.File
	var err error

	f, err = os.Create("tasks.json")
	if err != nil {
		fmt.Println(err)
	}
	defer f.Close()

	dat, _ := json.Marshal(m.state.tasks.tasks)
	f.Write(dat)
}

func ReadFromFile() ([]task, error) {
	dat, err := os.ReadFile("./tasks.json")
	if err != nil {
		return nil, fmt.Errorf("ile does not exist")
	}


	var tasks = []task{}

	err = json.Unmarshal(dat, &tasks)
	if err != nil {
		panic(err)
	}

	fmt.Printf("%#v", tasks)

	return tasks, nil
}
