package inputdto

import (
	sharedDTO "go-knowledge/services/eda/model-order-listener/internal/usecase/dtos/shared"
)

type ModelOrderDTO struct {
	Costumer        string                       `json:"costumer"`
	Context         string                       `json:"context"`
	Subcontexts     []sharedDTO.SubcontextDTO    `json:"subcontexts"`
	BucketName      string                       `json:"bucket_name"`
	FilesReferences []sharedDTO.FileReferenceDTO `json:"files_references"`
	Partition       string                       `json:"partition"`
}
