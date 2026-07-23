package main

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

type Book struct {
	Title  string `json:"title"`
	Author string `json:"author"`
}

type BookResponse struct {
	ID     uint   `json:"id"`
	Title  string `json:"title"`
	Author string `json:"author"`
}

func main() {
	if len(os.Args) < 2 {
		log.Fatalf("Usage: go run scripts/bulk_insert.go <csv_file> [api_url]\n  csv_file: path to CSV file (header: title,author)\n  api_url:  API base URL (default: http://localhost:8080)")
	}

	csvPath := os.Args[1]
	apiURL := "http://localhost:8080"
	if len(os.Args) >= 3 {
		apiURL = os.Args[2]
	}

	file, err := os.Open(csvPath)
	if err != nil {
		log.Fatalf("Failed to open CSV file: %v", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		log.Fatalf("Failed to parse CSV file: %v", err)
	}

	if len(records) == 0 {
		log.Fatal("CSV file is empty")
	}

	header := records[0]
	if len(header) < 2 || (header[0] != "title" && header[0] != "Title") || (header[1] != "author" && header[1] != "Author") {
		log.Fatalf("Invalid CSV header: %v. Expected format: title,author", header)
	}

	books := records[1:]
	total := len(books)
	if total == 0 {
		log.Fatal("No data rows found in CSV file")
	}

	log.Println("=== Bulk Insert Started ===")
	log.Printf("CSV file: %s", csvPath)
	log.Printf("API URL:  %s", apiURL)
	log.Printf("Total records to insert: %d", total)
	log.Println("----------------------------")

	client := &http.Client{Timeout: 10 * time.Second}

	var successCount, failCount int
	startTime := time.Now()

	for i, row := range books {
		if len(row) < 2 {
			log.Printf("[%d/%d] SKIP - Row has less than 2 columns: %v", i+1, total, row)
			failCount++
			continue
		}

		book := Book{
			Title:  row[0],
			Author: row[1],
		}

		body, err := json.Marshal(book)
		if err != nil {
			log.Printf("[%d/%d] ERROR marshaling JSON for '%s' by '%s': %v", i+1, total, book.Title, book.Author, err)
			failCount++
			continue
		}

		resp, err := client.Post(apiURL+"/books", "application/json", bytes.NewBuffer(body))
		if err != nil {
			log.Printf("[%d/%d] ERROR sending request for '%s' by '%s': %v", i+1, total, book.Title, book.Author, err)
			failCount++
			continue
		}

		respBody, _ := io.ReadAll(resp.Body)
		resp.Body.Close()

		if resp.StatusCode == http.StatusCreated {
			var created BookResponse
			if err := json.Unmarshal(respBody, &created); err == nil {
				log.Printf("[%d/%d] OK - Created book (id=%d): '%s' by '%s'", i+1, total, created.ID, created.Title, created.Author)
			} else {
				log.Printf("[%d/%d] OK - Created book: '%s' by '%s'", i+1, total, book.Title, book.Author)
			}
			successCount++
		} else {
			log.Printf("[%d/%d] FAIL (HTTP %d) - '%s' by '%s': %s", i+1, total, resp.StatusCode, book.Title, book.Author, string(respBody))
			failCount++
		}
	}

	elapsed := time.Since(startTime)

	log.Println("----------------------------")
	log.Println("=== Bulk Insert Completed ===")
	log.Printf("Total:   %d", total)
	log.Printf("Success: %d", successCount)
	log.Printf("Failed:  %d", failCount)
	log.Printf("Duration: %s", elapsed.Round(time.Millisecond))
	log.Printf("Avg per record: %s", time.Duration(int64(elapsed)/int64(total)).Round(time.Millisecond))

	if failCount > 0 {
		os.Exit(1)
	}
}

func init() {
	log.SetFlags(0)
	log.SetOutput(os.Stdout)
}
