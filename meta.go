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

func (meta *PackageMeta) ToJsonString() (content []byte, err error) {
    content, err = json.Marshal(meta)
    return
}
