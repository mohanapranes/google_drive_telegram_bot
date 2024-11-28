package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"google.golang.org/api/drive/v3"
	"google.golang.org/api/option"
)

func main() {

	srv, err := initGoogleDrive()
	if err != nil {
		log.Fatalf("Unable to create Drive service: %v", err)
	}

	tel := initTelegram()
	sendMessage(tel, srv)
}

func initTelegram() *tgbotapi.BotAPI {
	botToken := os.Getenv("BOT_API_KEY")
	if botToken == "" {
		log.Fatalf("BOT_API_KEY not found in .env file")
	}

	// Initialize the bot
	bot, err := tgbotapi.NewBotAPI(botToken)
	if err != nil {
		log.Panicf("Failed to create bot: %v", err)
	}

	// Log bot info
	bot.Debug = true
	log.Printf("Authorized on account %s", bot.Self.UserName)

	return bot
}

func initGoogleDrive() (*drive.Service, error) {
	ctx := context.Background()

	GoolgeDriveApiKey := os.Getenv("GOOGLE_DRIVE_API_KEY")

	// Initialize the Drive service with API key
	srv, err := drive.NewService(ctx, option.WithAPIKey(GoolgeDriveApiKey))
	if err != nil {
		log.Fatalf("Unable to create Drive service: %v", err)
		return nil, err
	}

	return srv, nil
}

func sendMessage(bot *tgbotapi.BotAPI, srv *drive.Service) {

	// Set up updates
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 120
	updates := bot.GetUpdatesChan(u)

	// Download a specific file by ID
	fileId := os.Getenv("FILE_ID")
	downloadFolder := os.Getenv("DOWNLOAD_FOLDER")
	downloadFileName := os.Getenv("FILE")
	filePath := downloadFolder + "/" + downloadFileName

	if err := downloadFile(srv, fileId, downloadFolder, downloadFileName); err != nil {
		log.Fatalf("Failed to download file: %v", err)
	}

	var mainChatID int64 = 0

	go func() {
		ticker := time.NewTicker(24 * time.Hour)
		defer ticker.Stop()

		for range ticker.C {
			if mainChatID != 0 {
				sendFile(bot, filePath, mainChatID)
			} else {
				log.Println("Main chat ID is not set. Skipping file sending.")
			}
		}
	}()

	// Process updates
	for update := range updates {
		chatID := update.FromChat().ID
		if update.Message == nil { // Ignore non-message updates
			return
		}
		if mainChatID == 0 || update.Message.Text == "/updateId" {
			mainChatID = chatID
			log.Printf("Updated main chat ID to: %d", mainChatID)
			reply := tgbotapi.NewMessage(chatID, "Main chat ID updated successfully!")
			bot.Send(reply)
			sendFile(bot, filePath, mainChatID)
		}
	}
}

func sendFile(bot *tgbotapi.BotAPI, filePath string, chatID int64) {

	file, err := os.Open(filePath)
	if err != nil {
		log.Printf("Error opening file: %v", err)
	}
	defer file.Close()

	// Create a new document message
	doc := tgbotapi.NewDocument(chatID, tgbotapi.FileReader{
		Name:   filePath,
		Reader: file,
	})

	// Send the file
	_, err = bot.Send(doc)
	if err != nil {
		log.Printf("Error sending document: %v", err)
	}

	log.Println("File sent successfully.")
}

func downloadFile(srv *drive.Service, fileId string, downloadDir string, dowloadFile string) error {
	// Get the file metadata
	file, err := srv.Files.Get(fileId).Fields("name, mimeType").Do()
	if err != nil {
		return fmt.Errorf("unable to get file metadata: %v", err)
	}

	// Create download directory if it doesn't exist
	if err := os.MkdirAll(downloadDir, 0755); err != nil {
		return fmt.Errorf("unable to create download directory: %v", err)
	}

	var resp *http.Response
	// Check if the file is a Google Doc
	if file.MimeType == "testing" {
		// Export Google Doc as PDF
		var ope *drive.Operation
		ope, err = srv.Files.Download(fileId).MimeType("application/vnd.openxmlformats-officedocument.wordprocessingml.document").Do()
		fmt.Printf("%v \n", ope)
		if err != nil {
			return fmt.Errorf("unable to export file: %v", err)
		}
	} else {
		// Download regular file
		resp, err = srv.Files.Get(fileId).Download()
		// fmt.Printf("%v", resp)
		if err != nil {
			return fmt.Errorf("unable to download file: %v", err)
		}
	}

	defer resp.Body.Close()

	// Create the output file
	outputPath := filepath.Join(downloadDir, dowloadFile)
	out, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("unable to create output file: %v", err)
	}
	defer out.Close()

	// Copy the file content
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return fmt.Errorf("unable to save file: %v", err)
	}

	fmt.Printf("Downloaded '%s' successfully to %s\n", dowloadFile, outputPath)
	return nil
}
