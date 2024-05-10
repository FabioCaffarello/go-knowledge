package entity

import (
	"errors"
	md5id "go-knowledge/libs/golang/shared/id/go-md5"
	"sort"
	"strings"
	"time"
)

type Subcontext struct {
	Name     string `json:"name",bson:"name"`
	Priority int    `json:"priority",bson:"priority"`
}

type FileReference struct {
	FileID md5id.ID `json:"file_id",bson:"file_id"`
	Name   string   `json:"name",bson:"name"`
}

type ModelOrder struct {
	ID              md5id.ID        `json:"_id",bson:"_id"`
	ModelID         md5id.ID        `json:"model_id",bson:"model_id"`
	Costumer        string          `json:"costumer",bson:"costumer"`
	Context         string          `json:"context",bson:"context"`
	Subcontexts     []Subcontext    `json:"subcontexts",bson:"subcontexts"`
	BucketName      string          `json:"bucket_name",bson:"bucket_name"`
	FilesReferences []FileReference `json:"files_references",bson:"files_references"`
	Partition       string          `json:"partition",bson:"partition"`
	CreatedAt       string          `json:"create_at",bson:"create_at"`
}

func NewModelOrder(
	costumer string,
	context string,
	subcontexts []Subcontext,
	bucketName string,
	files []string,
	partition string,
) (*ModelOrder, error) {
	modelOrder := &ModelOrder{
		Costumer:    costumer,
		Context:     context,
		Subcontexts: subcontexts,
		BucketName:  bucketName,
		Partition:   partition,
		CreatedAt:   time.Now().Format(time.RFC3339),
	}
	modelOrder.setFilesReferences(files)
	modelOrder.setModelID()
	modelOrder.setModelOrderID()

	if !modelOrder.isValid() {
		return nil, errors.New("ModelOrder is not valid")
	}
	return modelOrder, nil
}

func (m *ModelOrder) isValid() bool {
	if m.ID == "" {
		return false
	}
	if m.ModelID == "" {
		return false
	}
	if m.Costumer == "" {
		return false
	}
	if m.Context == "" {
		return false
	}
	if len(m.Subcontexts) == 0 {
		return false
	}
	if m.BucketName == "" {
		return false
	}
	if len(m.FilesReferences) == 0 {
		return false
	}
	if m.Partition == "" {
		return false
	}
	return true
}

func (m *ModelOrder) getBaseIDAsString() string {
	// Sort subcontexts by priority
	sort.Slice(m.Subcontexts, func(i, j int) bool {
		return m.Subcontexts[i].Priority < m.Subcontexts[j].Priority
	})

	var names []string
	for _, sc := range m.Subcontexts {
		names = append(names, sc.Name)
	}
	subcontextsNames := strings.Join(names, "")

	return m.Costumer + m.Context + subcontextsNames
}

func (m *ModelOrder) getModelIDAsString() string {
	baseID := m.getBaseIDAsString()
	return baseID + m.Partition
}

func (m *ModelOrder) getModelOrderIDAsString() string {
	baseID := m.getBaseIDAsString()
	return baseID + m.Partition
}

func (m *ModelOrder) setModelID() {
	modelID := m.getModelIDAsString()
	m.ModelID = md5id.GetIDFromString(modelID)
}

func (m *ModelOrder) setModelOrderID() {
	modelOrderID := m.getModelOrderIDAsString()
	m.ID = md5id.GetIDFromString(modelOrderID)
}

func (m *ModelOrder) setFilesReferences(files []string) {
	for _, file := range files {
		fileReference := FileReference{
			FileID: md5id.GetIDFromString(file),
			Name:   file,
		}
		m.FilesReferences = append(m.FilesReferences, fileReference)
	}
}
