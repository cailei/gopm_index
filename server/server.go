package gopm_index_server

import (
    "appengine"
    "appengine/datastore"
    "fmt"
    "gopm_index"
    "net/http"
)

var kind string = "PackageMeta"

func init() {
    http.HandleFunc("/all", handler_all)
    http.HandleFunc("/name_exists", handler_name_exists)
    http.HandleFunc("/publish", handler_publish)
}

func handler_all(w http.ResponseWriter, r *http.Request) {
    fmt.Fprint(w, "Hello, world!")
}

func handler_name_exists(w http.ResponseWriter, r *http.Request) {
    name := r.FormValue("name")
    if name == "" {
        http.Error(w, "Bad Request: package name is empty", http.StatusBadRequest)
        return
    }
    ctx := appengine.NewContext(r)
    key := datastore.NewKey(ctx, kind, name, 0, nil)
    entity := new(gopm_index.PackageMeta)
    if err := datastore.Get(ctx, key, entity); err != nil {
        if err == datastore.ErrNoSuchEntity {
            fmt.Fprintf(w, "0")
        } else {
            http.Error(w, err.Error(), http.StatusInternalServerError)
        }
        return
    }

    // package name exists
    fmt.Fprintf(w, "1")
}

func handler_publish(w http.ResponseWriter, r *http.Request) {

}
