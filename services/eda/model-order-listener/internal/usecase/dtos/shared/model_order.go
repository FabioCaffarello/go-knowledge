package shareddto

type SubcontextDTO struct {
    Name string `json:"name"`
    Priority int `json:"priority"`
}

type FileReferenceDTO struct {
	FileID string `json:"file_id"`
	Name   string `json:"name"`
}