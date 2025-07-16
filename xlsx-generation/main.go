package main

import (
	"bytes"
	"fmt"
	"html/template"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/xuri/excelize/v2"
)

func routeHandler(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("index.html")
	if err != nil {
		http.Error(w, "Template not found", http.StatusInternalServerError)
		return
	}
	tmpl.Execute(w, nil)
}

func prepData() *excelize.File {
	f := excelize.NewFile()
	sheet := f.GetSheetName(0)
	columns := strings.Split("ABCDEFGHIJKLMNOPQRSTUVWXYZ", "")
	for i := range 10000 {
		for _, col := range columns {
			cell := fmt.Sprintf("%s%d", col, i+1)
			loremIpsum := "Lorem ipsum dolor sit amet, consectetur adipiscing elit. Sed do eiusmod tempor incididunt ut labore et dolore magna aliqua."
			f.SetCellValue(sheet, cell, fmt.Sprintf("%s : %s", cell, loremIpsum))
		}
	}
	return f
}

func generateWithoutBuffer(w http.ResponseWriter, r *http.Request) {
	// Create new file
	f := prepData()

	// 2. Write the file to a buffer
	buf, err := f.WriteToBuffer()
	if err != nil {
		http.Error(w, "Failed to generate Excel file", http.StatusInternalServerError)
		return
	}

	// 3. Set response headers
	w.Header().Set("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	w.Header().Set("Content-Disposition", `attachment; filename="report.xlsx"`)

	// 4. Write to response
	w.Write(buf.Bytes())
}

func generateWithBuffer(w http.ResponseWriter, r *http.Request) {
	// Create new file
	f := prepData()

	// Set headers for download
	w.Header().Set("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	w.Header().Set("Content-Disposition", `attachment; filename="report.xlsx"`)

	// Write to response writer
	if err := f.Write(w); err != nil {
		http.Error(w, "Unable to generate file", http.StatusInternalServerError)
		return
	}
}

func generateWithDelayedBuffer(w http.ResponseWriter, r *http.Request) {
	// Create new file
	f := prepData()
	buf, err := f.WriteToBuffer()
	if err != nil {
		http.Error(w, "Failed to generate file", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	w.Header().Set("Content-Disposition", `attachment; filename="report.xlsx"`)
	w.WriteHeader(http.StatusOK)

	chunkSize := 1024 // 1 KB
	reader := bytes.NewReader(buf.Bytes())
	chunk := make([]byte, chunkSize)

	for {
		n, err := reader.Read(chunk)
		if n > 0 {
			_, writeErr := w.Write(chunk[:n])
			if writeErr != nil {
				fmt.Println("Error writing chunk:", writeErr)
				break
			}

			// Simulate network delay
			time.Sleep(10 * time.Millisecond)
		}
		if err == io.EOF {
			break
		}
		if err != nil {
			fmt.Println("Error reading buffer:", err)
			break
		}
	}
}

func main() {
	http.HandleFunc("/", routeHandler)
	// http.HandleFunc("/generate", generateWithoutBuffer)
	// http.HandleFunc("/generate", generateWithBuffer)
	http.HandleFunc("/generate", generateWithDelayedBuffer)

	http.ListenAndServe(":8080", nil)
}
