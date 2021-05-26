package main

import (
    "log"
    "net/http"
    "os"
    "io/ioutil"
    "strings"
    "path"
)

var gldata = make(map[string]string)

func store(w http.ResponseWriter, req *http.Request) {
    log.Println("Received store request.")
    key := req.Header.Get("key")

    if key == "" {
        w.Write([]byte("Key not found"))
        return
    }

    value := req.Header.Get("value")
    if value == "" {
        w.Write([]byte("Value not found"))
        return
    }

    gldata[key] = value

    w.Write([]byte("Stored " + key + " -> " + value))
}

func retrieve(w http.ResponseWriter, req *http.Request) {
    log.Println("Received retrieve request.")
    key := req.Header.Get("key")

    if key == "" {
        w.Write([]byte("Key not found"))
        return
    }

    if value, ok := gldata[key]; ok {
        w.Write([]byte(value))
    } else {
        w.Write([]byte("No data found for key: " + key))
    }
}

func clean(w http.ResponseWriter, req *http.Request) {
    log.Println("Received clean request.")
    gldata = make(map[string]string)
}

func listAll(w http.ResponseWriter, req *http.Request) {
    log.Println("Received list request.")
    result := ""

    for k, v := range gldata {
        result += k + ":" + v + "\n"
    }
    w.Write([]byte(result))
}

func upload(w http.ResponseWriter, req *http.Request) {
	req.ParseMultipartForm(10 << 30)

	file, headers, err := req.FormFile("filter")
	if err != nil {
		log.Println("Error Retrieving the File\n")
		log.Println(err)
		return
	}

  filename := path.Base((*headers).Filename)
	defer file.Close()
	log.Printf("Uploaded File: %s\n", filename)

	f, err := os.Create(filename)
	if err != nil {
		log.Println(err)
		return
	}
	defer f.Close()

	fileBytes, err := ioutil.ReadAll(file)
	if err != nil {
		log.Println(err)
		return
	}

	f.Write(fileBytes)
	log.Println("Successfully Uploaded File")
}

func getFilter(w http.ResponseWriter, req *http.Request) {
	filter := strings.TrimPrefix(req.URL.Path, "/filter/")
	log.Printf("Serving file: %s\n", filter)
	http.ServeFile(w, req, filter)
}

func main() {
    http.HandleFunc("/store", store)
    http.HandleFunc("/retrieve", retrieve)
    http.HandleFunc("/list", listAll)
    http.HandleFunc("/clean", clean)
    http.HandleFunc("/upload", upload)
    http.HandleFunc("/filter/", getFilter)
    log.Println("Starting server...")
    http.ListenAndServe(":8080", nil)
}
