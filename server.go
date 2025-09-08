package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

func uploadFile(w http.ResponseWriter, r *http.Request) {
	fmt.Println("uploading file ...")
	r.ParseMultipartForm(10 << 20)

	if r.Method != http.MethodPost {
		http.Error(w, "only POST method allowed", http.StatusMethodNotAllowed)
		return
	}

	file, handler, err := r.FormFile("newFile")
	if err != nil {
		fmt.Println("Error retreiving the file : ", err)
		return
	}
	defer file.Close()

	os.MkdirAll("./httpServer/", os.ModePerm)
	dst, err := os.Create("./httpServer/" + handler.Filename)
	if err != nil {
		fmt.Println("Error creating file : ", err)
		return
	}
	defer dst.Close()

	_, err = io.Copy(dst, file)
	if err != nil {
		fmt.Println("Error saving the file : ", err)
		return
	}
	fmt.Fprintf(w, "filename : %s \n size : %d \n status : uploaded", handler.Filename, handler.Size)

}

func downloadFile(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Only GET method allowed : ", http.StatusMethodNotAllowed)
		return
	}

	file := r.URL.Query().Get("file")
	if file == "" {
		http.Error(w, "file parameter is missing", http.StatusBadRequest)
		return
	}
	
	path := "./httpServer/" + file
	f, err := os.Open(path)
	if err != nil {
		http.Error(w, "file not found", http.StatusNotFound)
		return
	} 
	defer f.Close()

	w.Header().Set("Content-Disposition", "attachement;filename="+file)
	w.Header().Set("Content-Type", "application/octet-stream")

	io.Copy(w, f)

}

func main() {
	http.HandleFunc("/upload", uploadFile)
	http.HandleFunc("/download", downloadFile)
	fmt.Print("Running at http://localhost:8080")
	log.Fatal((http.ListenAndServe(":8080", nil)))
}
