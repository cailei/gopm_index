package gopm_index

import (
    "encoding/json"
)

type PackageMeta struct {
    Name         string       `json:"name"`
    Description  string       `json:"description"`
    Category     string       `json:"category"`
    Keywords     []string     `json:"keywords"`
    Author       PersonMeta   `json:"author"`
    Contributors []PersonMeta `json:"contributors"`
    Repositories []string     `json:"repositories"`
}

type PersonMeta struct {
    Name  string `json:"name"`
    Email string `json:"email"`
}

func (meta *PackageMeta) toJsonString() (str string, err error) {
    var data []byte
    data, err = json.Marshal(meta)
    if err != nil {
        return
    }
    str = string(data)
    return
}
