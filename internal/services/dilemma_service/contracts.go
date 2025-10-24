package dilemma_service

import (
	"context"

	"github.com/Woland-prj/dilemator/internal/domain/dto/dilemma_dto"
	"github.com/Woland-prj/dilemator/internal/domain/entity/dilemma_entity"
	"github.com/google/uuid"
)

//go:generate mockgen -destination=mock_service.go -package dilemma_service . DilemmaService

type DilemmaService interface {
	CreateDilemma(ctx context.Context, req *dilemma_dto.CreateDilemmaDto) (*dilemma_entity.Dilemma, error)
	GetByID(ctx context.Context, dilemmaID uuid.UUID) (*dilemma_entity.Dilemma, error)
	GetByOwner(ctx context.Context, ownerID uuid.UUID, page, size int) ([]dilemma_entity.Dilemma, error)
	UpdateDilemma(ctx context.Context, req *dilemma_dto.UpdateDilemmaDto) (*dilemma_entity.Dilemma, error)
	DeleteDilemma(ctx context.Context, dilemmaID uuid.UUID) error

	CreateDilemmaNode(ctx context.Context, req *dilemma_dto.CreateDilemmaNodeDto) (*dilemma_entity.DilemmaNode, error)
	GetChildren(ctx context.Context, parentID uuid.UUID) ([]*dilemma_entity.DilemmaNode, error)
	GetNodeByID(ctx context.Context, nodeID uuid.UUID) (*dilemma_entity.DilemmaNode, error)
	UpdateDilemmaNode(ctx context.Context, req dilemma_dto.UpdateDilemmaNodeDto) (*dilemma_entity.DilemmaNode, error)
	DeleteDilemmaNode(ctx context.Context, nodeID uuid.UUID) error
}
