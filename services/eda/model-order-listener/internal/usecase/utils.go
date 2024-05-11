package usecase

import (
	"go-knowledge/services/eda/model-order-listener/internal/entity"
	sharedDTO "go-knowledge/services/eda/model-order-listener/internal/usecase/dtos/shared"
)

func ConvertSubcontextsDTOToEntity(subcontexts []sharedDTO.SubcontextDTO) []entity.Subcontext {
	convertedSubcontexts := make([]entity.Subcontext, len(subcontexts))
	for i, subcontext := range subcontexts {
		convertedSubcontexts[i] = entity.Subcontext(subcontext)
	}
	return convertedSubcontexts
}

func GetFilesNameFromDTO(files []sharedDTO.FileDTO) []string {
    convertedFiles := make([]string, len(files))
    for i, file := range files {
        convertedFiles[i] = file.Name
    }
    return convertedFiles
}

func ConvertSubcontextsEntityToDTO(subcontexts []entity.Subcontext) []sharedDTO.SubcontextDTO {
    convertedSubcontexts := make([]sharedDTO.SubcontextDTO, len(subcontexts))
    for i, subcontext := range subcontexts {
        convertedSubcontexts[i] = sharedDTO.SubcontextDTO(subcontext)
    }
    return convertedSubcontexts
}
