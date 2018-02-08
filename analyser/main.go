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
	"sort"
	"strings"
)

const (
	dataPath   = "/home/Eduardo/go/src/github.com/ArangoGutierrez/MapReduce/data"
	reducePath = "/home/Eduardo/go/src/github.com/ArangoGutierrez/MapReduce/data/reduce"
)

type kv struct {
	Key   string
	Value int
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func loadStopWords() []string {
	laguajes, err := ioutil.ReadDir("stopWords/")
	check(err)
	stopWords := []string{"", "rt"}

	for _, laguaje := range laguajes {
		// Open the stopWords files to map
		l := fmt.Sprintf("stopWords/%s", laguaje.Name())
		f, err := os.Open(l)
		check(err)
		r := bufio.NewReader(f)
		scanner := bufio.NewScanner(r)

		for scanner.Scan() {
			scanner.Text()
			stopWords = append(stopWords, strings.ToLower(scanner.Text()))
		}
	}

	return stopWords
}

func main() {
	hashtag := flag.String("hashtag", "", "a string")
	flag.Parse()

	stopWords := loadStopWords()
	var m map[string]int
	var ss []kv
	// Open the reduced file to analyse
	rdcFile := fmt.Sprintf("%s/%s.reduce", reducePath, *hashtag)
	raw, err := ioutil.ReadFile(rdcFile)
	check(err)
	err = json.Unmarshal(raw, &m)
	check(err)

	for _, d := range stopWords {
		delete(m, d)
	}

	for k, v := range m {
		ss = append(ss, kv{k, v})
	}

	sort.Slice(ss, func(i, j int) bool {
		return ss[i].Value > ss[j].Value
	})

	fmt.Println(ss[:5])
	fmt.Println(ss[6:10])
	fmt.Println(ss[11:15])
	fmt.Println(ss[16:20])
}
