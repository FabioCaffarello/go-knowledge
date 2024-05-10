package outputdto

import sharedDTO "go-knowledge/services/eda/model-order-listener/internal/usecase/dtos/shared"

type FileOrderDTO struct {
	OrderID     string                    `json:"order_id"`
	ModelID     string                    `json:"model_id"`
	FileID      string                    `json:"file_id"`
	Costumer    string                    `json:"costumer"`
	Context     string                    `json:"context"`
	Subcontexts []sharedDTO.SubcontextDTO `json:"subcontexts"`
	BucketName  string                    `json:"bucket_name"`
	Partition   string                    `json:"partition"`
}
