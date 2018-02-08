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
	"context"
	"flag"
	"fmt"
	"log"
	"net/url"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/ChimeraCoder/anaconda"
)

var (
	consumerKey       = getenv("TWITTER_CONSUMER_KEY")
	consumerSecret    = getenv("TWITTER_CONSUMER_SECRET")
	accessToken       = getenv("TWITTER_ACCESS_TOKEN")
	accessTokenSecret = getenv("TWITTER_ACCESS_TOKEN_SECRET")
)

const dataPath = "SETYOUROWNDATAPATH"

func getenv(name string) string {
	v := os.Getenv(name)
	if v == "" {
		panic("missing required environment variable " + name)
	}
	return v
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

// searchTweets search for a query String
func searchTweets(ctx context.Context, file *os.File, queryString string, api *anaconda.TwitterApi) {
	searchResult, err := api.GetSearch(queryString, nil)
	if err != nil {
		panic(err)
	}

	for _, t := range searchResult.Statuses {
		if !t.Retweeted {
			_, err := file.WriteString(fmt.Sprintf("%s\n", t.Text))
			check(err)
			file.Sync()
		} else {
			continue
		}
	}
}

// trackHashtag track a given hashtag
func trackHashtag(ctx context.Context, tweet chan string, hashtag string, api *anaconda.TwitterApi) {
	stream := api.PublicStreamFilter(url.Values{
		"track": []string{fmt.Sprintf("#%s", hashtag)},
	})
	defer stream.Stop()

	for v := range stream.C {
		t, ok := v.(anaconda.Tweet)
		if !ok {
			log.Printf("received unexpected value of type %T", v)
			continue
		}
		tweet <- t.Text
	}
}

func printToFile(file *os.File, text string) {
	_, err := file.WriteString(fmt.Sprintf("%s\n", text))
	check(err)
	file.Sync()
}

func fileStat(filename string) {
	for {
		file, err := os.Stat(filename)
		check(err)

		log.Printf("Current file size: %.2fMB", float64(file.Size())/float64(1048576))
		time.Sleep(5 * time.Second)
	}
}

func main() {
	hashtag := flag.String("hashtag", "", "a string")
	trackTime := flag.Int("time", 5, "a int")
	flag.Parse()

	if *hashtag == "" {
		log.Fatalln("must declare a hashtag")
	}

	fileName := fmt.Sprintf("%s/%s.txt", dataPath, *hashtag)

	file, err := os.OpenFile(fileName, os.O_APPEND|os.O_WRONLY, os.ModeAppend)
	if err != nil {
		file, err = os.Create(fileName)
		check(err)
	}
	defer file.Close()

	anaconda.SetConsumerKey(consumerKey)
	anaconda.SetConsumerSecret(consumerSecret)
	api := anaconda.NewTwitterApi(accessToken, accessTokenSecret)

	ctx := context.Background()
	tweet := make(chan string)

	ctx, cancel := context.WithTimeout(ctx, time.Duration(*trackTime)*time.Minute)
	defer cancel()

	go trackHashtag(ctx, tweet, *hashtag, api)
	go fileStat(fileName)
	// Run!
	var gracefulStop = make(chan os.Signal)
	signal.Notify(gracefulStop, syscall.SIGTERM, syscall.SIGINT)

	for {
		select {
		case text := <-tweet:
			printToFile(file, text)
		case <-ctx.Done():
			log.Println("Time out!")
			os.Exit(0)
		case sig := <-gracefulStop:
			log.Printf("==> caught sig: %+v\n", sig)
			file.Sync()
			file.Close()
			os.Exit(0)
		}
	}
	fileInfo, err := os.Stat(fileName)
	check(err)
	log.Println("Total file size: ", fileInfo.Size())
}
