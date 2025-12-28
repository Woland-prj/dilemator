package dilemma_repo

import (
	"context"
	"errors"

	"github.com/Woland-prj/dilemator/internal/domain/entity/dilemma_entity"
	"github.com/Woland-prj/dilemator/internal/domain/errors/berrors"
	"github.com/Woland-prj/dilemator/internal/domain/errors/dilemma_errors"
	pentity "github.com/Woland-prj/dilemator/internal/repo/dilemma_repo/entity"
	"github.com/Woland-prj/dilemator/internal/services/dilemma_service"
	"github.com/Woland-prj/dilemator/pkg/logger"
	"github.com/Woland-prj/dilemator/pkg/postgres"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type DilemmaRepositoryAdapter struct {
	l logger.Interface
	*postgres.Postgres
}

var _ dilemma_service.DilemmaRepositoryPort = (*DilemmaRepositoryAdapter)(nil)

func NewDilemmaRepositoryAdapter(pg *postgres.Postgres, l logger.Interface) *DilemmaRepositoryAdapter {
	return &DilemmaRepositoryAdapter{
		Postgres: pg,
		l:        l,
	}
}

// SaveDilemmaDescriber сохраняет дилемму (без рекурсивного сохранения узлов).
func (r *DilemmaRepositoryAdapter) SaveDilemmaDescriber(ctx context.Context, dilemma *dilemma_entity.Dilemma) error {
	const op = "repo - dilemma_router - DilemmaRepositoryAdapter - SaveDilemmaDescriber"

	dilemmaEnt := pentity.DilemmaEntityFromModel(dilemma)

	var existing pentity.DilemmaEntity

	err := r.DB.WithContext(ctx).Select("id").First(&existing, "id = ?", dilemmaEnt.ID).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return berrors.InternalFromErr(op, err)
	}

	return r.DB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			if err := tx.Create(&dilemmaEnt).Error; err != nil {
				if errors.Is(err, gorm.ErrDuplicatedKey) {
					return berrors.FromErr(op, dilemma_errors.ErrDilemmaAlreadyExists)
				}
				return berrors.InternalFromErr(op, err)
			}
		} else {
			if err := tx.Model(&existing).Updates(&dilemmaEnt).Error; err != nil {
				return berrors.InternalFromErr(op, err)
			}
		}

		return nil
	})
}

// GetDilemmaWithRoot загружает дилемму вместе с корневым узлом.
func (r *DilemmaRepositoryAdapter) GetDilemmaWithRoot(ctx context.Context, dilemmaID uuid.UUID) (*dilemma_entity.Dilemma, error) {
	const op = "repo - dilemma_router - DilemmaRepositoryAdapter - GetDilemmaWithRoot"

	var dilemmaEnt pentity.DilemmaEntity
	if err := r.DB.WithContext(ctx).
		Model(&pentity.DilemmaEntity{}).
		Preload("RootNode").
		Preload("RootNode.Children").
		First(&dilemmaEnt, "id = ?", dilemmaID).
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

	var existing pentity.DilemmaNodeEntity

	err := r.DB.WithContext(ctx).Select("id").First(&existing, "id = ?", nodeEnt.ID).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return berrors.InternalFromErr(op, err)
	}

	return r.DB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			if err := tx.Create(&nodeEnt).Error; err != nil {
				if errors.Is(err, gorm.ErrDuplicatedKey) {
					return berrors.FromErr(op, dilemma_errors.ErrNodeAlreadyExists)
				}
				return berrors.InternalFromErr(op, err)
			}
		} else {
			if err := tx.Model(&existing).Updates(&nodeEnt).Error; err != nil {
				return berrors.InternalFromErr(op, err)
			}
		}

		return nil
	})
}

func (r *DilemmaRepositoryAdapter) GetNode(ctx context.Context, nodeID uuid.UUID) (*dilemma_entity.DilemmaNode, error) {
	const op = "repo - dilemma_router - DilemmaRepositoryAdapter - GetNode"

	tx := r.DB.WithContext(ctx).Begin()
	defer func() {
		if p := recover(); p != nil {
			tx.Rollback()
			panic(p)
		}
	}()

	var nodeEnt pentity.DilemmaNodeEntity

	err := tx.SetupJoinTable(&pentity.DilemmaNodeEntity{}, "Children", &pentity.NodeChildren{})
	if err != nil {
		tx.Rollback()
		return nil, berrors.InternalFromErr(op, err)
	}

	// Загружаем узел вместе с детьми (Children)
	if err := tx.
		Model(&pentity.DilemmaNodeEntity{}).
		Preload("Children").
		Where("id = ?", nodeID).
		First(&nodeEnt).
		Error; err != nil {
		tx.Rollback()
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, berrors.FromErr(op, dilemma_errors.ErrNodeNotFound)
		}
		return nil, berrors.InternalFromErr(op, err)
	}

	// Узнаем parent-id через join-таблицу node_children
	var parentID string
	if err := tx.
		Table(pentity.NodeChildrenTable).
		Select("node_id").
		Where("child_id = ?", nodeID).
		Find(&parentID).
		Error; err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		tx.Rollback()
		r.l.Debug("err here")
		return nil, berrors.InternalFromErr(op, err)
	}

	// Записываем parentID в сущность
	if parentID != "" {
		nodeEnt.ParentID = uuid.MustParse(parentID)
	} else {
		nodeEnt.ParentID = uuid.Nil
	}

	if err := tx.Commit().Error; err != nil {
		return nil, berrors.InternalFromErr(op, err)
	}

	r.l.Debug("node %+v", nodeEnt)
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
		Joins("JOIN dilemma_nodes dn ON node_childrens.child_id = dn.id").
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
