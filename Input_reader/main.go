//
//  Copyright 2017 The GoReduce Authors
//
//  Licensed under the Apache License, Version 2.0 (the "License");
//  you may not use this file except in compliance with the License.
//  You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//  Unless required by applicable law or agreed to in writing, software
//  distributed under the License is distributed on an "AS IS" BASIS,
//  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//  See the License for the specific language governing permissions and
//  limitations under the License.
//

package main

import (
	"flag"
	"fmt"
	"log"
	"math"
	"os"
	"path/filepath"
)

const (
	dataPath      = "/home/Eduardo/go/src/github.com/ArangoGutierrez/MapReduce/data"
	fileChunkSize = 2097152 // 2 MB
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func printToFile(file *os.File, text []byte) {
	_, err := file.WriteString(fmt.Sprintf("%s\n", text))
	check(err)
	file.Sync()
}

// splitFile reads the file line by line
// and split it into 2 MB file chunks
func splitFile(file *os.File, chunksNum uint64) {
	fileInfo, _ := file.Stat()
	fileSize := fileInfo.Size()
	fileName := filepath.Base(file.Name())
	ftracker := int64(0)

	for i := uint64(0); i <= chunksNum; i++ {
		// Create the chunk file
		chunkName := fmt.Sprintf("%s/dataChunck/chunk%d_of_%s.chunk", dataPath, i, fileName)
		f, err := os.Create(chunkName)
		check(err)
		defer f.Close()
		// Read a slice of the big file
		partSize := int(math.Min(fileChunkSize, float64(fileSize-int64(i*fileChunkSize))))
		partBuffer := make([]byte, partSize)
		_, err = file.Read(partBuffer)
		check(err)
		// write/save buffer to disk
		printToFile(f, partBuffer)
		finfo, _ := f.Stat()
		ftracker += finfo.Size()

		fmt.Printf("Chunk file=> %s\n", chunkName)
		// if ftracker >= fileSize {
		// 	break
		// }
	}
}

func totalPartsNum(file *os.File) uint64 {
	log.Println("Calculating number of file chunks for file", file.Name())

	fileInfo, _ := file.Stat()
	log.Println("File size:", (fileInfo.Size() / 1024), "MB")

	// calculate total number of parts the file will be chunked into
	totalPartsNum := uint64(fileInfo.Size() / fileChunkSize)
	fmt.Println(fileInfo.Size() / fileChunkSize)
	log.Println(totalPartsNum, "chunks files will be created from", file.Name())
	return totalPartsNum
}

func main() {
	hashtag := flag.String("hashtag", "", "a string")
	flag.Parse()
	fileToBeChunked := fmt.Sprintf("%s/%s.txt", dataPath, *hashtag)

	file, err := os.Open(fileToBeChunked)
	check(err)
	defer file.Close()

	// calculate total number of parts the file will be chunked into
	chunksNum := totalPartsNum(file)
	// Split the file on a calculated chunksNum
	splitFile(file, chunksNum)
}
