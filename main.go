package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strings"
	"time"
)

// CONFIG holds various configuration options
type CONFIG struct {
	Editor           string
	NotesDir         string
	GitEnabled       bool
	GitRepoPath      string
	NotionEnabled    bool
	NotionToken      string
	NotionDatabaseID string
}

func main() {
	// Load configuration
	config := loadConfig()

	// Create a new note filename
	noteFile, err := createNewNoteFile(config.NotesDir)
	if err != nil {
		log.Fatalf("Error creating new note file: %v", err)
	}

	//Open the note in the configured editor
	err = openEditor(config.Editor, noteFile)
	if err != nil {
		log.Fatalf("Error opening editor: %v", err)
	}

	// 4. (Optional) Commit and push to Git
	if config.GitEnabled {
		if err := gitCommitAndPush(noteFile, config); err != nil {
			log.Printf("Error committing/pushing to git: %v", err)
		} else {
			fmt.Println("Note committed and pushed to Git.")
		}
	}

	// todo:Upload to Notion
	if config.NotionEnabled {
		content, err := ioutil.ReadFile(noteFile)
		if err != nil {
			log.Printf("Error reading note file for Notion upload: %v", err)
		} else {
			err := uploadToNotion(config, string(content))
			if err != nil {
				log.Printf("Error uploading to Notion: %v", err)
			} else {
				fmt.Println("Note uploaded to Notion.")
			}
		}
	}

	fmt.Println("Done!")
}

// loadConfig loads configuration from environment variables (or from a file, if desired).
func loadConfig() *CONFIG {

	return &CONFIG{
		Editor:           getEnv("NOTE_EDITOR", "vim"),
		NotesDir:         getEnv("NOTES_DIR", "./notes"),
		GitEnabled:       getEnvBool("GIT_ENABLED", false),
		GitRepoPath:      getEnv("GIT_REPO_PATH", "./notes"),
		NotionEnabled:    getEnvBool("NOTION_ENABLED", false),
		NotionToken:      os.Getenv("NOTION_TOKEN"),       // If needed
		NotionDatabaseID: os.Getenv("NOTION_DATABASE_ID"), // If needed
	}
}

func createNewNoteFile(notesDir string) (string, error) {
	// Ensure the notes directory exists
	err := os.MkdirAll(notesDir, 0755)
	if err != nil {
		return "", err
	}

	// Create a note filename based on timestamp
	timestamp := time.Now().Format("2006-01-02_150405")
	fileName := fmt.Sprintf("note_%s.md", timestamp)
	fullPath := notesDir + string(os.PathSeparator) + fileName

	// Create an empty file
	file, err := os.Create(fullPath)
	if err != nil {
		return "", err
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {

		}
	}(file)

	// You can optionally write a header or template content here
	// e.g. "# My Note\n\n"
	return fullPath, nil
}

func openEditor(editor, filePath string) error {
	cmd := exec.Command(editor, filePath)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func gitCommitAndPush(noteFile string, config *CONFIG) error {
	// cd into the Git repository path
	if err := os.Chdir(config.GitRepoPath); err != nil {
		return fmt.Errorf("could not chdir to repo path: %w", err)
	}

	// Stage the file
	if err := runCmd("git", "add", noteFile); err != nil {
		return err
	}

	// Commit
	message := fmt.Sprintf("Add note: %s", noteFile)
	if err := runCmd("git", "commit", "-m", message); err != nil {
		return err
	}

	// Push
	if err := runCmd("git", "push"); err != nil {
		return err
	}

	return nil
}

func uploadToNotion(config *CONFIG, content string) error {
	// Example of how you might use a Notion library like github.com/jomei/notionapi
	// Below is a conceptual snippet — you’ll need to adapt it to your usage.

	/*
	   client := notion.NewClient(notion.Token(config.NotionToken))
	   // Create a new page in a given database:
	   newPage := notion.PageCreateRequest{
	       Parent: notion.Parent{
	           DatabaseID: notion.DatabaseID(config.NotionDatabaseID),
	       },
	       Properties: notion.Properties{
	           "Title": notion.TitleProperty{
	               Title: []notion.RichText{
	                   {
	                       Text: notion.Text{
	                           Content: "My New Note",
	                       },
	                   },
	               },
	           },
	       },
	       Children: []notion.Block{
	           notion.ParagraphBlock{
	               BasicBlock: notion.BasicBlock{
	                   Type: notion.BlockTypeParagraph,
	               },
	               Paragraph: notion.RichTextBlock{
	                   Text: []notion.RichText{
	                       {
	                           Text: notion.Text{
	                               Content: content,
	                           },
	                       },
	                   },
	               },
	           },
	       },
	   }

	   _, err := client.Page.Create(context.Background(), &newPage)
	   return err
	*/

	// Since we’re not actually using the Notion client here, just simulate:
	fmt.Println("Simulating Notion upload with content:")
	fmt.Println(strings.Repeat("-", 40))
	fmt.Println(content)
	fmt.Println(strings.Repeat("-", 40))

	return nil
}

// Helper to run a command and get combined output or error
func runCmd(name string, args ...string) error {
	cmd := exec.Command(name, args...)
	cmd.Stdin = os.Stdin
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout

	return cmd.Run()
}

// use os now. todo: load maybe from .syt file in linux. also maybe can define a env file
func getEnv(key, defaultVal string) string {
	val := os.Getenv(key)
	if val == "" {
		return defaultVal
	}
	return val
}

func getEnvBool(key string, defaultVal bool) bool {
	val := os.Getenv(key)
	if val == "" {
		return defaultVal
	}
	val = strings.ToLower(val)
	return val == "true" || val == "1"
}
