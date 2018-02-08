/*
 * Copyright 2017 The GoReduce Authors
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"strings"
)

const (
	dataPath    = "/home/Eduardo/go/src/github.com/ArangoGutierrez/MapReduce/data"
	chuncksPath = "/home/Eduardo/go/src/github.com/ArangoGutierrez/MapReduce/data/dataChunck"
	mapsPath    = "/home/Eduardo/go/src/github.com/ArangoGutierrez/MapReduce/data/maps"
	bigFile     = "/home/Eduardo/go/src/github.com/ArangoGutierrez/MapReduce/data/love.txt"
	bigFileName = "love.txt"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func mapper(file os.FileInfo, messages chan string) {
	m := make(map[string]int)
	var txt string
	// Open the chunk file to map
	chunck := fmt.Sprintf("%s/dataChunck/%s", dataPath, file.Name())
	f, err := os.Open(chunck)
	check(err)
	r := bufio.NewReader(f)
	scanner := bufio.NewScanner(r)
	// Scanner run to avoid '\n' being counted on as a string
	for scanner.Scan() {
		scanner.Text()
		txt = strings.ToLower(scanner.Text())
		// Make a Regex to say we only want
		reg, _ := regexp.Compile("[^a-zA-Z @#Ññ]+")
		processedString := reg.ReplaceAllString(txt, "")
		for _, v := range strings.Split(processedString, " ") {
			m[v]++
		}
	}
	//Define and create the map file
	mapFileName := fmt.Sprintf("%s/%s.map", mapsPath, file.Name())
	mapFile, err := os.Create(mapFileName)
	check(err)
	defer mapFile.Close()
	w := bufio.NewWriter(mapFile)
	// Print the mapper result
	enc := json.NewEncoder(w)
	enc.Encode(m)
	w.Flush()
	messages <- fmt.Sprintf("Map:%s", file.Name())
}

func main() {
	hashtag := flag.String("hashtag", "", "a string")
	flag.Parse()

	messages := make(chan string)
	chuncks, err := ioutil.ReadDir(chuncksPath)
	check(err)
	c := 0
	for _, chunck := range chuncks {
		if strings.Contains(chunck.Name(), *hashtag) {
			go mapper(chunck, messages)
			c++
		} else {
			continue
		}
	}

	for i := 0; i < c; i++ {
		msg := <-messages
		fmt.Println(msg)
	}

}
