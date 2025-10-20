package dilemma_dto

import "github.com/google/uuid"

type UpdateDilemmaDto struct {
	ID        uuid.UUID
	Topic     string
	RootValue string
}
