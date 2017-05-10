package main

import (
    "encoding/json"
    "net/http"
    "fmt"
    "encoding/base64"
)

// ExistsAPI is API implementation of /exists root endpoint
type ExistsAPI struct {
}

// Post is the handler for POST /exists
func (api ExistsAPI) Post(w http.ResponseWriter, r *http.Request) {
    rootkey := r.FormValue("rootkey")
    var keys []string

    fmt.Println(rootkey)
    if rootkey != "" {
        fmt.Println("Using HSET:", rootkey)
    }

	// decode request
	if err := json.NewDecoder(r.Body).Decode(&keys); err != nil {
        fmt.Println(err)
		w.WriteHeader(400)
		return
	}

    var klen = len(keys)

    // connect ardb
    client, err := getArdb()
    if err != nil {
        fmt.Println(err)
		w.WriteHeader(500)
		return
	}

    // build list
    var notFoundList = make([]string, klen)
    var notFound = 0

    for i := 0; i < klen; i++ {
        fmt.Printf("Decoding and checking: %s\n", keys[i])

        decoded, err := base64.StdEncoding.DecodeString(keys[i])
        if err != nil {
            fmt.Println(err)
            w.WriteHeader(500)
            return
        }

        if keyExists(client, string(decoded), rootkey) == false {
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
