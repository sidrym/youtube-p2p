package main

import (
    "crypto/md5"
    "encoding/hex"
    "fmt"
    "io"
    "io/ioutil"
    "os"
    "os/exec"
    "regexp"
    "strings"
)

func extractString(input string) string {
    re := regexp.MustCompile(`\[(.*?)\]`)
    match := re.FindStringSubmatch(input)
    if len(match) > 1 {
        return match[1]
    }
    return ""
}

func calculateChecksum(filePath string) (string, error) {
    file, err := os.Open(filePath)
    if err != nil {
        return "", err
    }
    defer file.Close()

    hash := md5.New()
    if _, err := io.Copy(hash, file); err != nil {
        return "", err
    }

    return hex.EncodeToString(hash.Sum(nil)), nil
}

func getLocalFileMap() map[string][2]string {
	// Define a map to store strings within brackets and their corresponding MD5 checksums
	stringChecksumMap := make(map[string][2]string)

	// Get the current directory
	currentDir, err := os.Getwd()
	if err != nil {
		fmt.Println("Error getting current directory:", err)
		return stringChecksumMap
	}

	// Read the files in the current directory
	files, err := ioutil.ReadDir(currentDir)
	if err != nil {
		fmt.Println("Error reading current directory:", err)
		return stringChecksumMap
	}

	// Iterate over the files
	for _, file := range files {
		if !file.IsDir() && strings.HasSuffix(file.Name(), ".webm") {
			fileName := file.Name()
			extractedString := extractString(fileName)
			if extractedString != "" {
				checksum, err := calculateChecksum(fileName)
				if err != nil {
					fmt.Printf("Error calculating checksum for file %s: %v\n", fileName, err)
					continue
				}
				// Store the extracted string and its corresponding MD5 checksum in the map
				stringChecksumMap[extractedString] = [2]string{checksum, fileName}
			}
		}
	}
	return stringChecksumMap
}

func playVideo(filePath string) error {
    // Check if the file exists
    _, err := os.Stat(filePath)
    if err != nil {
        if os.IsNotExist(err) {
            return fmt.Errorf("file %s does not exist", filePath)
        }
        return err
    }

    // File exists, play it using mpv
    cmd := exec.Command("mpv", filePath)
    cmd.Stdout = os.Stdout
    cmd.Stderr = os.Stderr
    fmt.Println("Playing video...")
    err = cmd.Run()
    if err != nil {
        return fmt.Errorf("error playing video: %w", err)
    }

    return nil
}

func downloadVideo(oid string) {
	videoURL := "https://www.youtube.com/watch?v=" + oid
	cmd := exec.Command("yt-dlp", videoURL)
	cmd.CombinedOutput()
}

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