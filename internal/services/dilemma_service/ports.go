package dilemma_service

import (
	"context"

	"github.com/Woland-prj/dilemator/internal/domain/entity/dilemma_entity"
	"github.com/google/uuid"
)

//go:generate mockgen -destination=mocks_test.go -package dilemma_service_test . DilemmaRepositoryPort

type DilemmaRepositoryPort interface {
	SaveDilemmaDescriber(ctx context.Context, dilemma *dilemma_entity.Dilemma) error
	GetDilemmaWithRoot(ctx context.Context, dilemmaID uuid.UUID) (*dilemma_entity.Dilemma, error)
	DeleteDilemma(ctx context.Context, dilemmaID uuid.UUID) error

	SaveNode(ctx context.Context, node *dilemma_entity.DilemmaNode) error
	GetNode(ctx context.Context, nodeID uuid.UUID) (*dilemma_entity.DilemmaNode, error)
	DeleteNode(ctx context.Context, nodeID uuid.UUID) error

	GetChildren(ctx context.Context, parentID uuid.UUID) ([]*dilemma_entity.DilemmaNode, error)
	LinkParentChild(ctx context.Context, parentID, childID uuid.UUID) error
}
