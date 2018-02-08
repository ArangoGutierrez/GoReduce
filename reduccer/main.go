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
	"strings"
)

const (
	dataPath   = "/home/Eduardo/go/src/github.com/ArangoGutierrez/MapReduce/data"
	reducePath = "/home/Eduardo/go/src/github.com/ArangoGutierrez/MapReduce/data/reduce"
	mapsPath   = "/home/Eduardo/go/src/github.com/ArangoGutierrez/MapReduce/data/maps"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func reduce(rm *map[string]int, mapFile os.FileInfo) {
	var dat map[string]int
	// Open the map file to reduce
	rdcFile := fmt.Sprintf("%s/%s", mapsPath, mapFile.Name())
	raw, err := ioutil.ReadFile(rdcFile)
	check(err)
	err = json.Unmarshal(raw, &dat)
	check(err)
	for k, v := range dat {
		(*rm)[k] += v
	}
}

func main() {
	hashtag := flag.String("hashtag", "", "a string")
	flag.Parse()

	reducedMaps := make(map[string]int)
	// Check for the maps
	maps, err := ioutil.ReadDir(mapsPath)
	check(err)
	// Read the maps and reduce them!
	for _, m := range maps {
		if strings.Contains(m.Name(), *hashtag) {
			reduce(&reducedMaps, m)
		} else {
			continue
		}
	}
	//Define and create the reduce file
	reduceFileName := fmt.Sprintf("%s/%s.reduce", reducePath, *hashtag)
	reduceFile, err := os.Create(reduceFileName)
	check(err)
	defer reduceFile.Close()
	w := bufio.NewWriter(reduceFile)
	// Print the reduccer result
	enc := json.NewEncoder(w)
	enc.Encode(reducedMaps)
	w.Flush()
}
