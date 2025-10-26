package dilemma_dto

import "github.com/google/uuid"

type UpdateDilemmaNodeDto struct {
	ID    uuid.UUID
	Name  string
	Value string
}
