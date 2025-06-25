package main

import (
  "strings"
  "context"
  "fmt"
  "os"
  "os/exec"
  "google.golang.org/genai"
  "time"
)

func main() {
	for {	
		time.Sleep(8 * time.Second)
		//take pic of chat
		scrotChat()
		//respond to pic of chat
		r, err := geminiPicPrompt("output.png")
		if err != nil {
			fmt.Println("Error:", err)
			deleteFile("output.png") // Always delete
			continue
		}
		if strings.Contains(r, "continue") {
			continue
		}
		//type response to chat
		typeResponse(r)
		//delete files
		deleteFile("output.png")
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

func scrotChat() {
        cmd := exec.Command("scrot", "output.png")

        cmd.Stdout = os.Stdout
        cmd.Stderr = os.Stderr

        err := cmd.Run()
        if err != nil {
                fmt.Println("Error:", err)
        }
}


func typeResponse(response string) {
	cmd := exec.Command("bash", "chatType.sh", response)

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Run()
	if err != nil {
		fmt.Println("Error:", err)
	}
}
/*
func geminiPicPrompt(file string) string {
  ctx := context.Background()
  client, _ := genai.NewClient(ctx, &genai.ClientConfig{
    APIKey:  os.Getenv("GEMINI_API_KEY"),
    Backend: genai.BackendGeminiAPI,
  })

  uploadedFile, _ := client.Files.UploadFromPath(ctx, file, nil)

  parts := []*genai.Part{
      genai.NewPartFromText("I'm sending you a screenshot of text messages, if the person on the left in the grey box is the last one to respond, come up with a response to continue the conversation(other wise just say the word continue), just say your response"),
      genai.NewPartFromURI(uploadedFile.URI, uploadedFile.MIMEType),
  }

  contents := []*genai.Content{
      genai.NewContentFromParts(parts, genai.RoleUser),
  }

  result, _ := client.Models.GenerateContent(
      ctx,
      "gemini-2.5-flash",
      contents,
      nil,
  )

  return result.Text()
}
*/

func geminiPicPrompt(file string) (string, error) {
	ctx := context.Background()
	client, err := genai.NewClient(ctx, &genai.ClientConfig{
		APIKey:  os.Getenv("GEMINI_API_KEY"),
		Backend: genai.BackendGeminiAPI,
	})
	if err != nil {
		return "", fmt.Errorf("failed to create Gemini client: %w", err)
	}

	uploadedFile, err := client.Files.UploadFromPath(ctx, file, nil)
	if err != nil {
		return "", fmt.Errorf("failed to upload file: %w", err)
	}

	parts := []*genai.Part{
		genai.NewPartFromText("I'm sending you a screenshot of text messages..."),
		genai.NewPartFromURI(uploadedFile.URI, uploadedFile.MIMEType),
	}
	contents := []*genai.Content{
		genai.NewContentFromParts(parts, genai.RoleUser),
	}

	result, err := client.Models.GenerateContent(
		ctx,
		"gemini-2.5-flash",
		contents,
		nil,
	)
	if err != nil {
		return "", fmt.Errorf("failed to generate content: %w", err)
	}
	return result.Text(), nil
}

