/*
This file is part of gopm (Go Package Manager)
Copyright (c) 2012 cailei (dancercl@gmail.com)

The MIT License (MIT)

Permission is hereby granted, free of charge, to any person obtaining a copy of
this software and associated documentation files (the "Software"), to deal in
the Software without restriction, including without limitation the rights to
use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies
of the Software, and to permit persons to whom the Software is furnished to do
so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
*/

package gopm_index_server

import (
    "appengine"
    "appengine/datastore"
    "bytes"
    "fmt"
    "gopm/index"
    "io"
    "net/http"
)

type FullIndex struct {
    content []byte
}

func init() {
    http.HandleFunc("/all", handlerGetFullIndex)
    http.HandleFunc("/publish", handlerPublishNewPackage)
}

func handlerGetFullIndex(w http.ResponseWriter, r *http.Request) {
    ctx := appengine.NewContext(r)
    key := getKeyOfFullIndexEntry(ctx)
    entity := new(FullIndex)
    err := datastore.Get(ctx, key, entity)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    reader := bytes.NewReader(entity.content)
    io.Copy(w, reader)
}

func getKeyOfFullIndexEntry(ctx appengine.Context) *datastore.Key {
    return datastore.NewKey(ctx, "FullIndex", "full_index", 0, nil)
}

func handlerPublishNewPackage(w http.ResponseWriter, r *http.Request) {
    json := r.FormValue("pkg")
    if json == "" {
        http.Error(w, "Form value pkg is empty!", http.StatusBadRequest)
        return
    }

    // unmarshal PackageMeta from json
    meta := new(index.PackageMeta)
    err := meta.FromJson([]byte(json))
    if err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    // check name uniqueness in the database
    ctx := appengine.NewContext(r)
    key := datastore.NewKey(ctx, "PackageMeta", meta.Name, 0, nil)
    entity := new(index.PackageMeta)

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
        var meta index.PackageMeta
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

        // append a newline
        io.Copy(buf, bytes.NewBuffer([]byte("\n")))
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
