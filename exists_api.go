package main

import (
    "github.com/go-redis/redis"
	"encoding/json"
	"net/http"
    "fmt"
)

// ExistsAPI is API implementation of /exists root endpoint
type ExistsAPI struct {
}

// Post is the handler for POST /exists
func (api ExistsAPI) Post(w http.ResponseWriter, r *http.Request) {
    rootkey := r.FormValue("rootkey")
    var keys []string

    fmt.Println(rootkey)

	// decode request
	if err := json.NewDecoder(r.Body).Decode(&keys); err != nil {
        fmt.Println(err)
		w.WriteHeader(400)
		return
	}

    var klen = len(keys)

    client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
    })

    pong, err := client.Ping().Result()
    fmt.Println(pong, err)

    var notFoundList = make([]string, klen)
    var notFound = 0

    for i := 0; i < klen; i++ {
        fmt.Printf("Testing: %s\n", keys[i])
        exists, _ := client.Exists(keys[i]).Result()

        if exists == 0 {
            notFoundList[notFound] = keys[i]
            notFound += 1
        }
    }

    client.Close()

    fmt.Println(notFound)
    output, err := json.Marshal(notFoundList[0:notFound])

	// uncomment below line to add header
	w.Header().Set("Content-Type", "application/json")

    w.Write(output)
}
