package dilemma_dto

import "github.com/google/uuid"

type CreateDilemmaDto struct {
	OwnerID   uuid.UUID
	Topic     string
	RootName  string
	RootValue string
	RootImage *FileDto
}
