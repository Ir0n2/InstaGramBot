package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"context"
    	"log"
	"google.golang.org/genai"
)
var filerage string
func main() {

	var subs [3]string = [3]string{"memes", "newvegasmemes", "starwarsmemes"}

	for i := 0; i <= 3; i++ {
	seen := make(map[string]bool)

	
	err := downloadOnePost(subs[i], "image", seen)
	if err != nil {
		fmt.Println("Error:", err)
	}
	// Wait or break as needed
	
	//fmt.Println(filerage)
	command()
	
	errr := deleteFile(filerage)
	if errr != nil {
		fmt.Println(errr)
	}
	}
}

func deleteFile(path string) error {
	err := os.Remove(path)
	if err != nil {
		return fmt.Errorf("failed to delete file %s: %w", path, err)
	}
	fmt.Printf("Deleted: %s\n", path)
	return nil
}


func command() {
	cmd := exec.Command("bash", "instaPost.sh", geminiPrompt("make a caption of 5 hashtags for memes. Please say nothing except for hashtags"))

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Run()
	if err != nil {
		fmt.Println("Error:", err)
	}
}

func downloadOnePost(subreddit, mediaType string, seen map[string]bool) error {
	url := "https://www.reddit.com/r/" + subreddit + ".json"
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("User-Agent", "SimpleGoClient/1.0")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return err
	}

	children := result["data"].(map[string]interface{})["children"].([]interface{})

	for _, c := range children {
		post := c.(map[string]interface{})["data"].(map[string]interface{})
		postID := post["id"].(string)

		// Skip if we've seen it
		if seen[postID] {
			continue
		}

		if mediaType == "image" {
			u := post["url"].(string)
			if isImage(u) {
				seen[postID] = true
				return download(u) // Done
			}
		} else if mediaType == "video" && post["is_video"].(bool) {
			media := post["media"].(map[string]interface{})
			videoURL := media["reddit_video"].(map[string]interface{})["fallback_url"].(string)
			seen[postID] = true
			return download(videoURL) // Done
		}
	}

	return fmt.Errorf("no new %s post found", mediaType)
}


func geminiPrompt(fuck string) string {
    ctx := context.Background()
    client, err := genai.NewClient(ctx, &genai.ClientConfig{
        APIKey:  "AIzaSyDuCHSFZnqrG6e5hjUJm6MMM6acjCVx0fI",
        Backend: genai.BackendGeminiAPI,
    })
    if err != nil {
        log.Fatal(err)
    }

    result, err := client.Models.GenerateContent(
        ctx,
        "gemini-2.5-flash",
        genai.Text(fuck),
        nil,
    )
    if err != nil {
        log.Fatal(err)
    }
    return result.Text()
}

func isImage(u string) bool {
	ext := strings.ToLower(filepath.Ext(u))
	return ext == ".jpg" || ext == ".jpeg" || ext == ".png" || ext == ".gif"
}

func download(url string) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	filename := "Desktop/" + filepath.Base(url)
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = io.Copy(file, resp.Body)
	if err == nil {
		fmt.Println("Downloaded:", filename)
	}
	
	filerage = filename
	return err
}

