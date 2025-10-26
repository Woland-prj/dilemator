package dilemma_repo

import (
	"context"
	"errors"

	"github.com/Woland-prj/dilemator/internal/domain/entity/dilemma_entity"
	"github.com/Woland-prj/dilemator/internal/domain/errors/berrors"
	"github.com/Woland-prj/dilemator/internal/domain/errors/dilemma_errors"
	pentity "github.com/Woland-prj/dilemator/internal/repo/dilemma_repo/entity"
	"github.com/Woland-prj/dilemator/internal/services/dilemma_service"
	"github.com/Woland-prj/dilemator/pkg/postgres"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type DilemmaRepositoryAdapter struct {
	*postgres.Postgres
}

var _ dilemma_service.DilemmaRepositoryPort = (*DilemmaRepositoryAdapter)(nil)

func NewDilemmaRepositoryAdapter(pg *postgres.Postgres) *DilemmaRepositoryAdapter {
	return &DilemmaRepositoryAdapter{
		Postgres: pg,
	}
}

// SaveDilemmaDescriber сохраняет дилемму (без рекурсивного сохранения узлов).
func (r *DilemmaRepositoryAdapter) SaveDilemmaDescriber(ctx context.Context, dilemma *dilemma_entity.Dilemma) error {
	const op = "repo - dilemma_router - DilemmaRepositoryAdapter - SaveDilemmaDescriber"

	dilemmaEnt := pentity.DilemmaEntityFromModel(dilemma)

	if err := r.DB.WithContext(ctx).Create(&dilemmaEnt).Error; err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return berrors.FromErr(op, dilemma_errors.ErrDilemmaAlreadyExists)
		}

		return berrors.InternalFromErr(op, err)
	}

	return nil
}

// GetDilemmaWithRoot загружает дилемму вместе с корневым узлом.
func (r *DilemmaRepositoryAdapter) GetDilemmaWithRoot(ctx context.Context, dilemmaID uuid.UUID) (*dilemma_entity.Dilemma, error) {
	const op = "repo - dilemma_router - DilemmaRepositoryAdapter - GetDilemmaWithRoot"

	var dilemmaEnt pentity.DilemmaEntity
	if err := r.DB.WithContext(ctx).
		Model(&pentity.DilemmaEntity{}).
		Preload("RootNode").
		Preload("RootNode.Children").
		Where(&dilemmaEnt, " = ?", dilemmaID).
		Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, berrors.FromErr(op, dilemma_errors.ErrDilemmaNotFound)
		}

		return nil, berrors.InternalFromErr(op, err)
	}

	return dilemmaEnt.ToModel(), nil
}

func (r *DilemmaRepositoryAdapter) GetDilemmasByOwner(ctx context.Context, ownerID uuid.UUID, page, size int) ([]dilemma_entity.Dilemma, error) {
	const op = "repo - dilemma_router - DilemmaRepositoryAdapter - GetDilemmasByOwner"

	offset := (page - 1) * size

	var entities []*pentity.DilemmaEntity

	if err := r.DB.WithContext(ctx).
		Model(&pentity.DilemmaEntity{}).
		Preload("RootNode").
		Where("owner_id = ?", ownerID).
		Offset(offset).
		Limit(size).
		Order("id ASC").
		Find(&entities).
		Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, berrors.FromErr(op, dilemma_errors.ErrDilemmaNotFound)
		}

		return nil, berrors.InternalFromErr(op, err)
	}

	return pentity.ToModelList(entities), nil
}

// DeleteDilemma удаляет дилемму (каскадное удаление узлов — за счёт FK в БД).
func (r *DilemmaRepositoryAdapter) DeleteDilemma(ctx context.Context, dilemmaID uuid.UUID) error {
	const op = "repo - dilemma_router - DilemmaRepositoryAdapter - DeleteDilemma"

	result := r.DB.WithContext(ctx).Delete(&pentity.DilemmaEntity{}, "id = ?", dilemmaID)
	if result.Error != nil {
		return berrors.InternalFromErr(op, result.Error)
	}

	if result.RowsAffected == 0 {
		return berrors.FromErr(op, dilemma_errors.ErrDilemmaNotFound)
	}

	return nil
}

// SaveNode сохраняет отдельный узел.
func (r *DilemmaRepositoryAdapter) SaveNode(ctx context.Context, node *dilemma_entity.DilemmaNode) error {
	const op = "repo - dilemma_router - DilemmaRepositoryAdapter - SaveNode"

	nodeEnt := pentity.DilemmaNodeEntityFromModel(node)

	if err := r.DB.WithContext(ctx).Create(&nodeEnt).Error; err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return berrors.FromErr(op, dilemma_errors.ErrNodeAlreadyExists)
		}

		return berrors.InternalFromErr(op, err)
	}

	return nil
}

// GetNode возвращает узел по ID.
func (r *DilemmaRepositoryAdapter) GetNode(ctx context.Context, nodeID uuid.UUID) (*dilemma_entity.DilemmaNode, error) {
	const op = "repo - dilemma_router - DilemmaRepositoryAdapter - GetNode"

	tx := r.DB.WithContext(ctx).Begin()

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			panic(r)
		}
	}()

	var nodeEnt pentity.DilemmaNodeEntity

	if err := tx.
		Model(&pentity.DilemmaNodeEntity{}).
		Preload("Children").
		First(&nodeEnt, "id = ?", nodeID).
		Error; err != nil {
		tx.Rollback()

		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, berrors.FromErr(op, dilemma_errors.ErrNodeNotFound)
		}

		return nil, berrors.InternalFromErr(op, err)
	}

	var parentID uuid.UUID

	err := tx.
		Table(pentity.NodeChildrenTable).
		Select("node_id").
		Where("child_id = ?", nodeID).
		Scan(&parentID).
		Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		tx.Rollback()

		return nil, berrors.InternalFromErr(op, err)
	}

	nodeEnt.ParentID = parentID

	if err := tx.Commit().Error; err != nil {
		return nil, berrors.InternalFromErr(op, err)
	}

	return nodeEnt.ToModel(), nil
}

// DeleteNode удаляет узел (и связи в node_children — каскадно).
func (r *DilemmaRepositoryAdapter) DeleteNode(ctx context.Context, nodeID uuid.UUID) error {
	const op = "repo - dilemma_router - DilemmaRepositoryAdapter - DeleteNode"

	result := r.DB.WithContext(ctx).Delete(&pentity.DilemmaNodeEntity{}, "id = ?", nodeID)
	if result.Error != nil {
		return berrors.InternalFromErr(op, result.Error)
	}

	if result.RowsAffected == 0 {
		return berrors.FromErr(op, dilemma_errors.ErrNodeNotFound)
	}

	return nil
}

// GetChildren возвращает дочерние узлы по parentID через join-таблицу.
func (r *DilemmaRepositoryAdapter) GetChildren(ctx context.Context, parentID uuid.UUID) ([]*dilemma_entity.DilemmaNode, error) {
	const op = "repo - dilemma_router - DilemmaRepositoryAdapter - GetChildren"

	var childrenEnt []*pentity.DilemmaNodeEntity
	if err := r.DB.WithContext(ctx).
		Table(pentity.NodeChildrenTable).
		Select("dn.id, dn.value").
		Joins("JOIN dilemma_nodes dn ON node_children.child_id = dn.id").
		Where("node_children.node_id = ?", parentID).
		Find(&childrenEnt).
		Error; err != nil {
		return nil, berrors.InternalFromErr(op, err)
	}

	children := make([]*dilemma_entity.DilemmaNode, len(childrenEnt))
	for i, ent := range childrenEnt {
		children[i] = ent.ToModel()
	}

	return children, nil
}

// LinkParentChild создаёт связь в таблице node_children.
func (r *DilemmaRepositoryAdapter) LinkParentChild(ctx context.Context, parentID, childID uuid.UUID) error {
	const op = "repo - dilemma_router - DilemmaRepositoryAdapter - LinkParentChild"

	link := map[string]interface{}{
		"node_id":  parentID,
		"child_id": childID,
	}

	if err := r.DB.WithContext(ctx).Table(pentity.NodeChildrenTable).Create(link).Error; err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return berrors.FromErr(op, dilemma_errors.ErrNodeAlreadyHasParent)
		}

		return berrors.InternalFromErr(op, err)
	}

	return nil
}
