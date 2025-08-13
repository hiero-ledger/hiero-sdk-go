package main

import (
	"bufio"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"time"
)

var brokenlinks []string
var red = "\033[31m"
var reset = "\033[0m"
var broken int = 0
var working int = 0

func main() {
	err := filepath.Walk(".", func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && filepath.Ext(path) == ".md" {
			fmt.Println(red, path, reset)
			err := readMd(path)
			if err != nil {
				fmt.Println("Error")
			}
		}
		return nil
	})

	if err != nil {
		fmt.Println("Fehler beim Durchsuchen:", err)
	}

	fmt.Println(red, "\nBroken Links:", reset)
	for _, link := range brokenlinks {
		fmt.Println(link)
	}
	stats()
}

func readMd(path string) error {
	file, err := os.Open(path)
	if err != nil {
		fmt.Println("Error while reading:", path, "Err:", err)
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	scanner.Buffer(make([]byte, 0, 64*1024), 1024*1024)
	re := regexp.MustCompile(`https?://[^\)\s]+`)
	linenumber := 0
	for scanner.Scan() {
		linenumber++
		line := scanner.Text()
		matches := re.FindAllString(line, -1)
		for _, link := range matches {
			link, StatusH, StatusG := checkLink(link)
			index := strconv.Itoa(linenumber)
			if StatusH != "passed" || StatusG != "passed" {
				brokenlinks = append(brokenlinks, fmt.Sprintf("%s in Line %s Link: %s StatusCodeHead: %s StatusCodeGet: %s", path, index, link, StatusH, StatusG))
				broken++
			} else {
				working++
			}
		}
	}

	return nil
}

func checkLink(link string) (string, string, string) {
	client := http.Client{
		Timeout: 5 * time.Second,
	}
	req, _ := http.NewRequest("HEAD", link, nil)
	req.Header.Set("User-Agent", "Mozilla/5.0")
	resp, err := client.Do(req)
	if resp != nil {
		defer resp.Body.Close()
	}

	if err != nil || (resp != nil && resp.StatusCode >= 400) {
		reqG, _ := http.NewRequest("GET", link, nil)
		reqG.Header.Set("User-Agent", "Mozilla/5.0")
		respG, errG := client.Do(reqG)
		if respG != nil {
			defer respG.Body.Close()
		}
		if errG != nil || (respG != nil && respG.StatusCode >= 400) {
			statusHead := "?"
			if resp != nil {
				statusHead = strconv.Itoa(resp.StatusCode)
			}
			statusGet := "?"
			if respG != nil {
				statusGet = strconv.Itoa(respG.StatusCode)
			}

			return link, statusHead, statusGet
		} else {
			fmt.Println("Passed:", link)
			return "passed", "passed", "passed"
		}
	} else {
		fmt.Println("Passed:", link)
		return "passed", "passed", "passed"
	}
}

func stats() {
	fmt.Println("Working:", working, "Broken:", broken)
}
