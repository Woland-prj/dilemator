package dilemma_service

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/Woland-prj/dilemator/internal/domain/dto/dilemma_dto"
	"github.com/Woland-prj/dilemator/internal/domain/entity/dilemma_entity"
	"github.com/Woland-prj/dilemator/internal/domain/errors/berrors"
	"github.com/Woland-prj/dilemator/internal/domain/errors/dilemma_errors"
	"github.com/Woland-prj/dilemator/pkg/logger"
	"github.com/google/uuid"
)

type dilemmaService struct {
	log           logger.Interface
	repo          DilemmaRepositoryPort
	fileRepo      FileRepositoryPort
	chatGenerator ChatGeneratorPort
}

var _ DilemmaService = (*dilemmaService)(nil)

func NewDilemmaService(
	log logger.Interface,
	repo DilemmaRepositoryPort,
	fileRepo FileRepositoryPort,
	chatGenerator ChatGeneratorPort,
) DilemmaService {
	return &dilemmaService{
		log:           log,
		repo:          repo,
		fileRepo:      fileRepo,
		chatGenerator: chatGenerator,
	}
}

func (s *dilemmaService) saveImage(ctx context.Context, image []byte, contentType string) (*string, error) {
	const op = "dilemma - dilemmaService - saveImage"
	if len(image) == 0 {
		return nil, nil
	}
	key, err := s.fileRepo.Save(ctx, image, contentType)
	if err != nil {
		s.log.Error(fmt.Sprintf("%s: failed to save image: %v", op, err))

		return nil, berrors.InternalFromErr(op, err)
	}
	return &key, nil
}

func (s *dilemmaService) getImageLink(ctx context.Context, imgKey *string) (*string, error) {
	const op = "dilemma - dilemmaService - getImageLink"
	if imgKey == nil {
		return nil, nil
	}
	link, err := s.fileRepo.GetAuthorizedDownloadLink(ctx, *imgKey)
	if err != nil {
		s.log.Error(fmt.Sprintf("%s: failed to get image link: %v", op, err))

		return nil, berrors.InternalFromErr(op, err)
	}
	return &link, nil
}

func (s *dilemmaService) deleteImage(ctx context.Context, imgKey *string) error {
	const op = "dilemma - dilemmaService - deleteImage"
	if imgKey == nil {
		return nil
	}
	s.log.Debug(fmt.Sprintf("imgKey %s", *imgKey))
	if err := s.fileRepo.DeleteByKey(ctx, *imgKey); err != nil {
		s.log.Error(fmt.Sprintf("%s: failed to delete image: %v", op, err))

		return berrors.InternalFromErr(op, err)
	}
	return nil
}

// CreateDilemma создаёт новую дилемму с корневым узлом.
func (s *dilemmaService) CreateDilemma(
	ctx context.Context,
	req *dilemma_dto.CreateDilemmaDto,
) (*dilemma_entity.Dilemma, error) {
	const op = "dilemma - dilemmaService - CreateDilemma"

	// Сохраняем изображение
	var imgKey *string
	if req.RootImage != nil {
		imgKey, err := s.saveImage(
			ctx,
			req.RootImage.Data,
			req.RootImage.ContentType,
		)
		if err != nil {
			return nil, err
		}

		s.log.Debug("saved image", slog.Any("key", imgKey))
	}

	rootNode := dilemma_entity.NewDilemmaNode(
		uuid.New(),
		uuid.Nil,
		req.RootName,
		req.RootValue,
		imgKey,
	)
	dilemma := dilemma_entity.NewDilemma(
		uuid.New(),
		req.OwnerID,
		req.Topic,
		rootNode,
	)

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

	// Получаем ссылку для отображения
	var err error
	dilemma.RootNode.Image, err = s.getImageLink(ctx, dilemma.RootNode.Image)
	if err != nil {
		return nil, err
	}

	return dilemma, nil
}

// GetByID возвращает полную дилемму (без рекурсивной загрузки дерева — только корень).
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

	// Ссылка на изображение
	dilemma.RootNode.Image, err = s.getImageLink(ctx, dilemma.RootNode.Image)
	if err != nil {
		return nil, err
	}

	return dilemma, nil
}

// GetByOwner ищет все дилеммы по ID пользователя.
func (s *dilemmaService) GetByOwner(
	ctx context.Context,
	ownerID uuid.UUID,
	page, size int,
) ([]dilemma_entity.Dilemma, error) {
	const op = "dilemma - dilemmaService - GetByOwner"

	dilemmas, err := s.repo.GetDilemmasByOwner(ctx, ownerID, page, size)
	if err != nil {
		return nil, berrors.InternalFromErr(op, err)
	}

	for _, dilemma := range dilemmas {
		// Ссылка на изображение
		dilemma.RootNode.Image, err = s.getImageLink(ctx, dilemma.RootNode.Image)
		if err != nil {
			return nil, err
		}
	}

	return dilemmas, nil
}

// UpdateDilemma обновляет тему дилеммы и значение корневого узла.
func (s *dilemmaService) UpdateDilemma(
	ctx context.Context,
	req *dilemma_dto.UpdateDilemmaDto,
) (*dilemma_entity.Dilemma, error) {
	const op = "dilemma - dilemmaService - UpdateDilemma"

	// Получаем текущую дилемму
	existing, err := s.repo.GetDilemmaWithRoot(ctx, req.ID)
	if err != nil {
		return nil, berrors.FromErr(op, err)
	}

	// Обновляем изображение
	if req.RootImage != nil && len(req.RootImage.Data) != 0 {
		newKey, err := s.saveImage(ctx, req.RootImage.Data, req.RootImage.ContentType)
		if err != nil {
			return nil, err
		}

		if err := s.deleteImage(ctx, existing.RootNode.Image); err != nil {
			return nil, err
		}

		existing.RootNode.Image = newKey
	}

	// Обновляем корневой узел
	existing.RootNode.Value = req.RootValue

	existing.RootNode.Name = req.RootName
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

	// Ссылка на изображение
	existing.RootNode.Image, err = s.getImageLink(ctx, existing.RootNode.Image)
	if err != nil {
		return nil, err
	}

	return existing, nil
}

// DeleteDilemma удаляет дилемму (каскадное удаление узлов обеспечивается БД).
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

// CreateDilemmaNode создаёт новый узел и связывает его с родителем через node_children.
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

	// Сохраеяем изображение
	var imgKey *string
	if req.Image != nil {
		imgKey, err = s.saveImage(
			ctx,
			req.Image.Data,
			req.Image.ContentType,
		)
		if err != nil {
			return nil, err
		}
	}

	node := dilemma_entity.NewDilemmaNode(
		uuid.New(),
		req.ParentID,
		req.Name,
		req.Value,
		imgKey,
	)

	if err := s.repo.SaveNode(ctx, node); err != nil {
		s.log.Error(fmt.Sprintf("%s: failed to save node: %v", op, err))

		return nil, berrors.InternalFromErr(op, err)
	}

	if err := s.repo.LinkParentChild(ctx, req.ParentID, node.ID); err != nil {
		s.log.Error(fmt.Sprintf("%s: failed to link parent-child: %v", op, err))

		return nil, berrors.InternalFromErr(op, err)
	}

	// Ссылка на изображение
	node.Image, err = s.getImageLink(ctx, node.Image)
	if err != nil {
		return nil, err
	}

	return node, nil
}

// GetNodeByID возвращает узел по ID.
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
		s.log.Debug(fmt.Sprintf("%s: node %+v, err: %s", op, node, err.Error()))
		return nil, berrors.InternalFromErr(op, err)
	}

	// Ссылка на изображение
	node.Image, err = s.getImageLink(ctx, node.Image)
	if err != nil {
		return nil, err
	}

	return node, nil
}

// GetChildren возвращает дочерние узлы для parentID.
func (s *dilemmaService) GetChildren(
	ctx context.Context,
	parentID uuid.UUID,
) ([]*dilemma_entity.DilemmaNode, error) {
	const op = "dilemma - dilemmaService - GetChildren"

	children, err := s.repo.GetChildren(ctx, parentID)
	if err != nil {
		if errors.Is(err, dilemma_errors.ErrNodeNotFound) {
			// Родитель может не иметь детей
			return []*dilemma_entity.DilemmaNode{}, nil
		}

		return nil, berrors.InternalFromErr(op, err)
	}

	return children, nil
}

// UpdateDilemmaNode обновляет значение узла.
func (s *dilemmaService) UpdateDilemmaNode(
	ctx context.Context,
	req dilemma_dto.UpdateDilemmaNodeDto,
) (*dilemma_entity.DilemmaNode, error) {
	const op = "dilemma - dilemmaService - UpdateDilemmaNode"

	node, err := s.repo.GetNode(ctx, req.ID)
	if err != nil {
		return nil, err
	}

	// Обновляем изображение
	if req.Image != nil && len(req.Image.Data) != 0 {
		newKey, err := s.saveImage(ctx, req.Image.Data, req.Image.ContentType)
		if err != nil {
			return nil, err
		}

		if err := s.deleteImage(ctx, node.Image); err != nil {
			return nil, err
		}

		node.Image = newKey
	}

	node.Name = req.Name
	node.Value = req.Value

	if err := s.repo.SaveNode(ctx, node); err != nil {
		s.log.Error(fmt.Sprintf("%s: failed to update node: %v", op, err))

		return nil, berrors.InternalFromErr(op, err)
	}

	return node, nil
}

func (s *dilemmaService) GenerateDilemmaNode(
	ctx context.Context,
	parentID uuid.UUID,
) (*dilemma_entity.DilemmaNode, error) {
	const op = "dilemma - dilemmaService - GenerateDilemmaNode"

	parentNode, err := s.repo.GetNodeWithParents(ctx, parentID)
	if err != nil {
		return nil, berrors.Wrap(op, "parent node not found", err)
	}

	generatedNode, err := s.chatGenerator.GenerateNode(ctx, parentNode)
	if err != nil {
		return nil, berrors.Wrap(op, "failed to generate node", err)
	}

	node := dilemma_entity.NewDilemmaNode(
		uuid.New(),
		parentID,
		generatedNode.Name,
		generatedNode.Value,
		nil,
	)

	if err := s.repo.SaveNode(ctx, node); err != nil {
		s.log.Error(fmt.Sprintf("%s: failed to save node: %v", op, err))

		return nil, berrors.InternalFromErr(op, err)
	}

	if err := s.repo.LinkParentChild(ctx, parentID, node.ID); err != nil {
		s.log.Error(fmt.Sprintf("%s: failed to link parent-child: %v", op, err))

		return nil, berrors.InternalFromErr(op, err)
	}

	return node, nil
}

// DeleteDilemmaNode удаляет узел (и всё поддерево, если БД поддерживает CASCADE).
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
