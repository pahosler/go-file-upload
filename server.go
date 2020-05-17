package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

func main() {
	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/", fs)
	http.HandleFunc("/upload", uploadRoute)

	log.Println("Listening on :3000...")
	log.Fatal(http.ListenAndServe(":3000", nil))
}

func uploadRoute(w http.ResponseWriter, r *http.Request) {
	r.ParseMultipartForm(10 << 20)

	file, handle, err := r.FormFile("file")
	if err != nil {
		log.Fatal(err)
	}

	defer file.Close()

	timestamp := time.Now().Unix()
	asString := fmt.Sprintf("%v", timestamp)

	f, err := os.OpenFile("upload/"+asString+"_"+handle.Filename, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		log.Fatal(err)
	}

	defer f.Close()

	io.Copy(f, file)

	fmt.Fprintf(w, "File has been uploaded\n")
	fmt.Fprintf(w, "File name: %v\n", handle.Filename)
	fmt.Fprintf(w, "File size: %v\n", handle.Size)
	fmt.Fprintf(w, "File type: %v\n", handle.Header["Content-Type"])
}
