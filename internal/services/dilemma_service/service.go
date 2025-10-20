package dilemma_service

import (
	"context"
	"errors"
	"fmt"

	"github.com/Woland-prj/dilemator/internal/domain/dto/dilemma_dto"
	"github.com/Woland-prj/dilemator/internal/domain/entity/dilemma_entity"
	"github.com/Woland-prj/dilemator/internal/domain/errors/berrors"
	"github.com/Woland-prj/dilemator/internal/domain/errors/dilemma_errors"
	"github.com/Woland-prj/dilemator/pkg/logger"
	"github.com/google/uuid"
)

type dilemmaService struct {
	log  logger.Interface
	repo DilemmaRepositoryPort
}

var _ DilemmaService = (*dilemmaService)(nil)

func NewDilemmaService(
	log logger.Interface,
	repo DilemmaRepositoryPort,
) DilemmaService {
	return &dilemmaService{
		log:  log,
		repo: repo,
	}
}

// CreateDilemma создаёт новую дилемму с корневым узлом
func (s *dilemmaService) CreateDilemma(
	ctx context.Context,
	req *dilemma_dto.CreateDilemmaDto,
) (*dilemma_entity.Dilemma, error) {
	const op = "dilemma - dilemmaService - CreateDilemma"

	rootNode := dilemma_entity.NewDilemmaNode(uuid.New(), req.RootValue)
	dilemma := dilemma_entity.NewDilemma(uuid.New(), req.OwnerID, req.Topic, rootNode)

	// Сохраняем корневой узел
	if err := s.repo.SaveNode(ctx, rootNode); err != nil {
		s.log.Error(fmt.Sprintf("%s: failed to save root node: %v", op, err))
		return nil, berrors.InternalFromErr(op, err)
	}

	// Сохраняем дилемму (ссылается на root_node_id)
	if err := s.repo.SaveDilemmaDescriber(ctx, dilemma); err != nil {
		s.log.Error(fmt.Sprintf("%s: failed to save dilemma: %v", op, err))
		return nil, berrors.InternalFromErr(op, err)
	}

	return dilemma, nil
}

// GetByID возвращает полную дилемму (без рекурсивной загрузки дерева — только корень)
func (s *dilemmaService) GetByID(
	ctx context.Context,
	dilemmaID uuid.UUID,
) (*dilemma_entity.Dilemma, error) {
	const op = "dilemma - dilemmaService - GetByID"

	// Предполагается, что репозиторий загружает дилемму + корневой узел
	// Если репозиторий не поддерживает это — нужно будет расширить его метод
	dilemma, err := s.repo.GetDilemmaWithRoot(ctx, dilemmaID)
	if err != nil {
		if errors.Is(err, dilemma_errors.ErrDilemmaNotFound) {
			s.log.Debug(fmt.Sprintf("%s: dilemma %s not found", op, dilemmaID))
			return nil, berrors.Wrap(op, fmt.Sprintf("dilemma %s not found", dilemmaID), err)
		}
		return nil, berrors.InternalFromErr(op, err)
	}

	return dilemma, nil
}

// UpdateDilemma обновляет тему дилеммы и значение корневого узла
func (s *dilemmaService) UpdateDilemma(
	ctx context.Context,
	req *dilemma_dto.UpdateDilemmaDto,
) (*dilemma_entity.Dilemma, error) {
	const op = "dilemma - dilemmaService - UpdateDilemma"

	// Получаем текущую дилемму
	existing, err := s.GetByID(ctx, req.ID)
	if err != nil {
		return nil, berrors.FromErr(op, err)
	}

	// Обновляем корневой узел
	existing.RootNode.Value = req.RootValue
	if err := s.repo.SaveNode(ctx, existing.RootNode); err != nil {
		s.log.Error(fmt.Sprintf("%s: failed to update root node: %v", op, err))
		return nil, berrors.InternalFromErr(op, err)
	}

	// Обновляем дилемму
	existing.Topic = req.Topic
	if err := s.repo.SaveDilemmaDescriber(ctx, existing); err != nil {
		s.log.Error(fmt.Sprintf("%s: failed to update dilemma: %v", op, err))
		return nil, berrors.InternalFromErr(op, err)
	}

	return existing, nil
}

// DeleteDilemma удаляет дилемму (каскадное удаление узлов обеспечивается БД)
func (s *dilemmaService) DeleteDilemma(
	ctx context.Context,
	dilemmaID uuid.UUID,
) error {
	const op = "dilemma - dilemmaService - DeleteDilemma"

	if err := s.repo.DeleteDilemma(ctx, dilemmaID); err != nil {
		if errors.Is(err, dilemma_errors.ErrDilemmaNotFound) {
			s.log.Debug(fmt.Sprintf("%s: dilemma %s not found", op, dilemmaID))
			return berrors.Wrap(op, fmt.Sprintf("dilemma %s not found", dilemmaID), err)
		}
		return berrors.InternalFromErr(op, err)
	}

	return nil
}

// CreateDilemmaNode создаёт новый узел и связывает его с родителем через node_children
func (s *dilemmaService) CreateDilemmaNode(
	ctx context.Context,
	req *dilemma_dto.CreateDilemmaNodeDto,
) (*dilemma_entity.DilemmaNode, error) {
	const op = "dilemma - dilemmaService - CreateDilemmaNode"

	// Проверим, существует ли родитель (опционально, если БД не проверяет FK)
	_, err := s.GetNodeByID(ctx, req.ParentID)
	if err != nil {
		return nil, berrors.Wrap(op, "parent node not found", dilemma_errors.ErrNodeNotFound)
	}

	node := dilemma_entity.NewDilemmaNode(uuid.New(), req.Value)

	if err := s.repo.SaveNode(ctx, node); err != nil {
		s.log.Error(fmt.Sprintf("%s: failed to save node: %v", op, err))
		return nil, berrors.InternalFromErr(op, err)
	}

	if err := s.repo.LinkParentChild(ctx, req.ParentID, node.ID); err != nil {
		s.log.Error(fmt.Sprintf("%s: failed to link parent-child: %v", op, err))
		return nil, berrors.InternalFromErr(op, err)
	}

	return node, nil
}

// GetNodeByID возвращает узел по ID
func (s *dilemmaService) GetNodeByID(
	ctx context.Context,
	nodeID uuid.UUID,
) (*dilemma_entity.DilemmaNode, error) {
	const op = "dilemma - dilemmaService - GetNodeByID"

	node, err := s.repo.GetNode(ctx, nodeID)
	if err != nil {
		if errors.Is(err, dilemma_errors.ErrNodeNotFound) {
			s.log.Debug(fmt.Sprintf("%s: node %s not found", op, nodeID))
			return nil, berrors.Wrap(op, fmt.Sprintf("node %s not found", nodeID), err)
		}
		return nil, berrors.InternalFromErr(op, err)
	}

	return node, nil
}

// GetChildren возвращает дочерние узлы для parentID
func (s *dilemmaService) GetChildren(
	ctx context.Context,
	parentID uuid.UUID,
) ([]*dilemma_entity.DilemmaNode, error) {
	const op = "dilemma - dilemmaService - GetChildren"

	children, err := s.repo.GetChildren(ctx, parentID)
	if err != nil {
		if errors.Is(err, dilemma_errors.ErrNodeNotFound) {
			// Родитель может не иметь детей — это не ошибка
			return []*dilemma_entity.DilemmaNode{}, nil
		}
		return nil, berrors.InternalFromErr(op, err)
	}

	return children, nil
}

// UpdateDilemmaNode обновляет значение узла
func (s *dilemmaService) UpdateDilemmaNode(
	ctx context.Context,
	req dilemma_dto.UpdateDilemmaNodeDto,
) (*dilemma_entity.DilemmaNode, error) {
	const op = "dilemma - dilemmaService - UpdateDilemmaNode"

	node, err := s.GetNodeByID(ctx, req.ID)
	if err != nil {
		return nil, err
	}

	node.Value = req.Value

	if err := s.repo.SaveNode(ctx, node); err != nil {
		s.log.Error(fmt.Sprintf("%s: failed to update node: %v", op, err))
		return nil, berrors.InternalFromErr(op, err)
	}

	return node, nil
}

// DeleteDilemmaNode удаляет узел (и всё поддерево, если БД поддерживает CASCADE)
func (s *dilemmaService) DeleteDilemmaNode(
	ctx context.Context,
	nodeID uuid.UUID,
) error {
	const op = "dilemma - dilemmaService - DeleteDilemmaNode"

	if err := s.repo.DeleteNode(ctx, nodeID); err != nil {
		if errors.Is(err, dilemma_errors.ErrNodeNotFound) {
			s.log.Debug(fmt.Sprintf("%s: node %s not found", op, nodeID))
			return berrors.Wrap(op, fmt.Sprintf("node %s not found", nodeID), err)
		}
		return berrors.InternalFromErr(op, err)
	}

	return nil
}
