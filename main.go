package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/schollz/progressbar/v3"
)

func parseCmd(args []string) (string, int64, error) {

	cmd := args[1:]
	l := len(cmd)
	var n int64
	var err error

	if l > 2 {
		return "", 0, errors.New(fmt.Sprintf("Too many args provided %v, usage: jspand <filepath> <desired number of bytes>(optional, default is 500mb)", args))
	}

	if l == 1 {
		n = 500000000
	} else {
		n, err = strconv.ParseInt(cmd[1], 10, 64)
		if err != nil {
			return "", 0, errors.New(fmt.Sprintf("Error parsing provided integer: %v", err))
		}
	}

	path := strings.ToLower(cmd[0])

	return path, n, nil

}

func parseJSON(jsonPath string) map[string]any {
	file, err := os.Open(jsonPath)
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}

	defer file.Close()

	var hashmap map[string]any
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&hashmap); err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}

	return hashmap

}

func main() {

	// jspand takes a small json file and blows it up so that it can be used for testing stuff
	// this works by simply duplicating existing data and changing keys slightly
	// benefit of this is that it will retain your schema but scale it up

	// parse args
	//      - ex. jspand [filepath] [optional integer]
	//      - integer arg is optional and if not provided defaults to 50000?

	// testCMD := []string{"jspand", "../fastTravel/prod/fastTravel.JSON"}
	filepath, size, err := parseCmd(os.Args)
	if err != nil {
		fmt.Println(err)
	}

	// validate
	//      - verify it's a json file
	//      - verify it's under predetermined size

	_, ext, found := strings.Cut(filepath, ".json")
	if !found || ext != "" {
		fmt.Printf("'%v' is not a json file or has '.json' in it's name - ", filepath)
		os.Exit(1)
	}

	f, err := os.Stat(filepath)
	if err != nil {
		fmt.Printf("unable to obtain file info for determining size of file '%v'", filepath)
	}

	if fs := f.Size(); fs > size {
		fmt.Printf("file size (%v bytes) is larger than max byte size (%v bytes)", fs, size)
		os.Exit(1)
	}

	// take keys and duplicate key: value pairs
	//      - ex. key -> key1, anotherKey -> anotherKey1
	//      - repeat until size limit reached

	data := parseJSON(filepath)
	newFilepath, _, _ := strings.Cut(filepath, ".json")
	newFilepath = newFilepath + "_BFF.json"
	newFile, err := os.Create(newFilepath)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer newFile.Close()

	var buffer bytes.Buffer
	buffer.Write([]byte("{"))
	newFile.Write(buffer.Bytes())
	buffer.Reset()

	firstKey := true
	iterator := 1

	fmt.Println("jspand is generating a bigger and better file for you...")
	bar := progressbar.Default(size)

	for fileSize := int64(0); fileSize < size; iterator++ {
		for k, v := range data {

			var key string
			if iterator != 1 {
				key = k + strconv.Itoa(iterator)
			} else {
				key = k
			}
			pair, err := json.Marshal(map[string]interface{}{key: v})
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}

			if !firstKey {
				buffer.Write([]byte(","))
			} else {
				firstKey = false
			}

			buffer.Write(pair[1 : len(pair)-1])

		}

		newFile.Write(buffer.Bytes())
		buffer.Reset()

		info, err := newFile.Stat()
		if err != nil {
			fmt.Printf("unable to obtain file info for determining size of file '%v'", filepath)
		}
		fileSize = info.Size()

		bar.Set64(fileSize)
		// fmt.Printf("i: %v ", iterator)
		// fmt.Printf("size: %v ", fileSize)
		// fmt.Printf("limit: %v ", size)
		// fmt.Println(fileSize < size)
	}

	buffer.Write([]byte("}"))
	newFile.Write(buffer.Bytes())
	buffer.Reset()

	fmt.Printf("'%v' created successfully!", newFilepath)
}
