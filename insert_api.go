package main

import (
	"net/http"
    "fmt"
    "io/ioutil"
)

// InsertAPI is API implementation of /insert root endpoint
type InsertAPI struct {
}

// Put is the handler for PUT /insert
func (api InsertAPI) Put(w http.ResponseWriter, r *http.Request) {
	// uncomment below line to add header
	// w.Header().Set("key","value")

    rootkey := r.FormValue("rootkey")
    fmt.Println(rootkey)

    r.ParseMultipartForm(8192)

    // connect ardb
    client, err := getArdb()
    if err != nil {
        fmt.Println(err)
		w.WriteHeader(500)
		return
	}

    m := r.MultipartForm
    files := m.File["files[]"]

    for i, _ := range files {
        file, err := files[i].Open()
        defer file.Close()

        if err != nil {
            w.WriteHeader(400)
            fmt.Println(err)
            return
        }

        fmt.Println("Trying to insert:", files[i].Filename)
        exists, _ := client.Exists(files[i].Filename).Result()

        if exists == 1 {
            w.WriteHeader(401)
            fmt.Println("Key exists, overwrite denied:", files[i].Filename)
            return
        }

        buffer, err := ioutil.ReadAll(file)
        if err != nil {
            w.WriteHeader(500)
            fmt.Println(err)
            return
        }

        client.Set(files[i].Filename, buffer, 0)
    }

    client.Close()
}
