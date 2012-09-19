package gopm_index_server

import (
    "appengine"
    "appengine/datastore"
    "fmt"
    "gopm_index"
    "net/http"
)

func init() {
    http.HandleFunc("/all", handler_all)
    http.HandleFunc("/publish", handler_publish)
}

func handler_all(w http.ResponseWriter, r *http.Request) {
    fmt.Fprint(w, "Hello, world!")
}

func handler_publish(w http.ResponseWriter, r *http.Request) {
    json := r.FormValue("pkg")
    if json == "" {
        http.Error(w, "Form value pkg is empty!", http.StatusBadRequest)
        return
    }

    // unmarshal PackageMeta from json
    meta := new(gopm_index.PackageMeta)
    err := meta.FromJson([]byte(json))
    if err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    // check name uniqueness in the database
    ctx := appengine.NewContext(r)
    key := datastore.NewKey(ctx, "PackageMeta", meta.Name, 0, nil)
    entity := new(gopm_index.PackageMeta)

    err = datastore.Get(ctx, key, entity)
    if err == nil {
        fmt.Print("i am nil")
        http.Error(w, fmt.Sprintf("The package name '%v' already exists in the index registry.\n", meta.Name), http.StatusInternalServerError)
        return
    }

    fmt.Printf("%#v\n", err)

    if err != datastore.ErrNoSuchEntity {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    _, err = datastore.Put(ctx, key, meta)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    fmt.Printf("Published a package '%v'\n", meta.Name)
}
