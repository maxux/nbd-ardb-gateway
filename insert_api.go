package main

import (
    "github.com/go-redis/redis"
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

    client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
    })

    pong, err := client.Ping().Result()
    fmt.Println(pong, err)

    r.ParseMultipartForm(8192)

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
