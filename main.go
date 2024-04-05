package main

import (
	"fmt"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run main.go [URL]")
		return
	}

	fileToFind := os.Args[1]
	stringChecksumMap := getLocalFileMap()

	resultMap, found := stringChecksumMap[fileToFind]
	if !found {
		videoID := os.Args[1]
		fmt.Printf("Downloading video %s...\n", videoID)
		downloadVideo(videoID)
	}

	stringChecksumMap = getLocalFileMap()
	resultMap, found = stringChecksumMap[fileToFind]

	fmt.Printf("Found video: %s\n", resultMap[1])
	err := playVideo(resultMap[1])
	if err != nil {
		fmt.Printf("Error playing video: %v\n", err)
	}
}
