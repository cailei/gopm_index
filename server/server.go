package gopm_index_server

import (
    "appengine"
    "appengine/datastore"
    "bytes"
    "fmt"
    "gopm_index"
    "io"
    "net/http"
)

type FullIndex struct {
    Content []byte
}

func init() {
    http.HandleFunc("/all", handler_all)
    http.HandleFunc("/publish", handler_publish)
}

func handler_all(w http.ResponseWriter, r *http.Request) {
    ctx := appengine.NewContext(r)
    key := datastore.NewKey(ctx, "FullIndex", "full_index", 0, nil)
    index := new(FullIndex)
    err := datastore.Get(ctx, key, index)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    reader := bytes.NewReader(index.Content)
    io.Copy(w, reader)
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
        http.Error(w, fmt.Sprintf("The package name '%v' already exists in the index registry.\n", meta.Name), http.StatusInternalServerError)
        return
    }

    if err != datastore.ErrNoSuchEntity {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    _, err = datastore.Put(ctx, key, meta)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    update_full_index(ctx, w)

    fmt.Printf("Published a package '%v'\n", meta.Name)
}

func update_full_index(ctx appengine.Context, w http.ResponseWriter) {
    // collect all packages into a single string (full index)
    query := datastore.NewQuery("PackageMeta").Order("Name")
    buf := bytes.NewBuffer(nil)
    for it := query.Run(ctx); ; {
        var meta gopm_index.PackageMeta
        _, err := it.Next(&meta)

        if err == datastore.Done {
            break
        }

        if err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }

        // serialize PackageMeta to json
        json, err := meta.ToJson()
        if err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }

        // append to the buffer
        reader := bytes.NewReader(json)
        written, err := io.Copy(buf, reader)
        if err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }
        if written != int64(len(json)) {
            http.Error(w, fmt.Sprintf("Indexing package '%v' failed.", meta.Name), http.StatusInternalServerError)
            return
        }
    }

    // store full index to a special place
    index := &FullIndex{buf.Bytes()}
    key := datastore.NewKey(ctx, "FullIndex", "full_index", 0, nil)
    _, err := datastore.Put(ctx, key, index)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
}
