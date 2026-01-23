package dilemma_service

import (
	"context"

	"github.com/Woland-prj/dilemator/internal/domain/entity/dilemma_entity"
	"github.com/google/uuid"
)

//go:generate mockgen -destination=mocks_test.go -package dilemma_service_test . DilemmaRepositoryPort

type ChatGeneratorPort interface {
	GenerateNode(ctx context.Context, parentNode *dilemma_entity.DilemmaNode) (*dilemma_entity.DilemmaNode, error)
}

type DilemmaRepositoryPort interface {
	SaveDilemmaDescriber(ctx context.Context, dilemma *dilemma_entity.Dilemma) error
	GetDilemmaWithRoot(ctx context.Context, dilemmaID uuid.UUID) (*dilemma_entity.Dilemma, error)
	GetDilemmasByOwner(ctx context.Context, ownerID uuid.UUID, page, size int) ([]dilemma_entity.Dilemma, error)
	DeleteDilemma(ctx context.Context, dilemmaID uuid.UUID) error

	SaveNode(ctx context.Context, node *dilemma_entity.DilemmaNode) error
	GetNode(ctx context.Context, nodeID uuid.UUID) (*dilemma_entity.DilemmaNode, error)
	GetNodeWithParents(ctx context.Context, nodeID uuid.UUID) (*dilemma_entity.DilemmaNode, error)
	DeleteNode(ctx context.Context, nodeID uuid.UUID) error

	GetChildren(ctx context.Context, parentID uuid.UUID) ([]*dilemma_entity.DilemmaNode, error)
	LinkParentChild(ctx context.Context, parentID, childID uuid.UUID) error
}

// FileRepositoryPort defines a contract for interacting with file storage.
type FileRepositoryPort interface {
	Save(ctx context.Context, file []byte, contentType string) (string, error)
	DeleteByKey(ctx context.Context, key string) error
	GetDownloadLink(ctx context.Context, key string) (string, error)
	GetAuthorizedDownloadLink(ctx context.Context, key string) (string, error)
}
