package gopm_index

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
