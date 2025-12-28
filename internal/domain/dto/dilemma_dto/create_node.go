package dilemma_dto

import "github.com/google/uuid"

type CreateDilemmaNodeDto struct {
	ParentID uuid.UUID
	Name     string
	Value    string
	Image    *FileDto
}
