package main

import (
    "log"
    "net/http"
    "time"

    "github.com/Jumpscale/go-raml/nbd-ardb-gateway/goraml"

    "github.com/gorilla/mux"
    "gopkg.in/validator.v2"

    "github.com/go-redis/redis"
)

func getArdb() (*redis.Client, error) {
    client := redis.NewClient(&redis.Options{
        Addr:     "localhost:6379",
        Password: "", // no password set
        DB:       0,  // use default DB
    })

    _, err := client.Ping().Result()

    return client, err
}

func keyExists(client *redis.Client, key string, rootkey string) bool {
    if rootkey != "" {
        existsBool, _ := client.HExists(rootkey, key).Result()
        return existsBool
    }

    exists, _ := client.Exists(key).Result()
    return (exists == 1)
}

func main() {
    // input validator
    validator.SetValidationFunc("multipleOf", goraml.MultipleOf)

    r := mux.NewRouter()

    // home page
    r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        http.ServeFile(w, r, "index.html")
    })

    // apidocs
    r.PathPrefix("/apidocs/").Handler(http.StripPrefix("/apidocs/", http.FileServer(http.Dir("./apidocs/"))))

    ExistsInterfaceRoutes(r, ExistsAPI{})

    InsertInterfaceRoutes(r, InsertAPI{})

    log.Println("starting server")
    s := &http.Server{
        Addr:           ":5000",
        Handler:        r,
        ReadTimeout:    3600 * time.Second,
        WriteTimeout:   3600 * time.Second,
        IdleTimeout:    3600 * time.Second,
        MaxHeaderBytes: 1 << 20,
    }

    s.ListenAndServe()
}
