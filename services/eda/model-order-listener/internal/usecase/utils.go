package usecase

import (
	"go-knowledge/services/eda/model-order-listener/internal/entity"
	sharedDTO "go-knowledge/services/eda/model-order-listener/internal/usecase/dtos/shared"
)

// func ConvertModelOrderDTOToEntity(modelOrderDTO inputDTO.ModelOrderDTO) entity.ModelOrder {
// 	modelOrder := entity.ModelOrder{
// 		Costumer:        modelOrderDTO.Costumer,
// 		Context:         modelOrderDTO.Context,
// 		Subcontexts:     ConvertSubcontextsDTOToEntity(modelOrderDTO.Subcontexts),
// 		BucketName:      modelOrderDTO.BucketName,
// 		FilesReferences: ConvertFilesReferencesDTOToEntity(modelOrderDTO.FilesReferences),
// 		Partition:       modelOrderDTO.Partition,
// 	}
// 	return modelOrder
// }

func ConvertSubcontextsDTOToEntity(subcontexts []sharedDTO.SubcontextDTO) []entity.Subcontext {
	convertedSubcontexts := make([]entity.Subcontext, len(subcontexts))
	for i, subcontext := range subcontexts {
		convertedSubcontexts[i] = entity.Subcontext(subcontext)
	}
	return convertedSubcontexts
}

func ConvertFilesReferencesDTOToEntity(filesReferences []sharedDTO.FileReferenceDTO) []entity.FileReference {
	convertedFilesReferences := make([]entity.FileReference, len(filesReferences))
	for i, fileReference := range filesReferences {
		convertedFilesReferences[i] = entity.FileReference(fileReference)
	}
	return convertedFilesReferences
}
