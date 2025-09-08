package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

func uploadFile(w http.ResponseWriter, r *http.Request) {
	fmt.Println("uploading file ...")

	if r.Method != http.MethodPost {
		http.Error(w, "only POST method allowed", http.StatusMethodNotAllowed)
		return
	}

	err := r.ParseMultipartForm(10 << 20) // 10 MB
	if err != nil {
		http.Error(w, "failed to parse form: "+err.Error(), http.StatusBadRequest)
		return
	}

	file, handler, err := r.FormFile("newFile")
	if err != nil {
		http.Error(w, "error retrieving file: "+err.Error(), http.StatusBadRequest)
		return
	}
	defer file.Close()

	os.MkdirAll("./httpServer/", os.ModePerm)

	dst, err := os.Create("./httpServer/" + filepath.Base(handler.Filename))
	if err != nil {
		http.Error(w, "error creating file: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer dst.Close()

	_, err = io.Copy(dst, file)
	if err != nil {
		http.Error(w, "error saving file: "+err.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "filename: %s\nsize: %d\nstatus: uploaded", handler.Filename, handler.Size)
}

func downloadFile(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "only GET method allowed", http.StatusMethodNotAllowed)
		return
	}

	file := r.URL.Query().Get("file")
	if file == "" {
		http.Error(w, "file parameter is missing", http.StatusBadRequest)
		return
	}

	safeFile := filepath.Base(file)
	path := "./httpServer/" + safeFile

	f, err := os.Open(path)
	if err != nil {
		http.Error(w, "file not found", http.StatusNotFound)
		return
	}
	defer f.Close()

	w.Header().Set("Content-Disposition", "attachment; filename="+safeFile)
	w.Header().Set("Content-Type", "application/octet-stream")

	io.Copy(w, f)
}

func listAll(w http.ResponseWriter, r *http.Request) {
	file, err := os.ReadDir("./httpServer/")
	if err != nil {
		http.Error(w, "Error while fetching Directory", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html")
	fmt.Fprintln(w, "<h2>available files are listed below</h2>")

	for _, f := range file {
	name := f.Name()
	fmt.Fprintf(w, `<li><a href="/download?file=%s">%s</a></li>`, name, name)
	}

	fmt.Fprintln(w, "</ul>")
}

func deleteFile(w http.ResponseWriter, r *http.Request) {
	file := r.URL.Query().Get("file")

	if file == "" {
		http.Error(w, "file missing, enter valid file name", http.StatusBadRequest)
		return
	}

	err := os.Remove("./httpServer" + file)
	if err != nil {
		http.Error(w, "Error while deleting the file", http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "%s deleted successfully", file)
}

func main() {
	http.HandleFunc("/upload", uploadFile)
	http.HandleFunc("/download", downloadFile)
	http.HandleFunc("/", listAll)
	http.HandleFunc("/delete", deleteFile)
	fmt.Println("Running at http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
