package entity

import (
	"github.com/Woland-prj/dilemator/internal/domain/entity/dilemma_entity"
	"github.com/google/uuid"
)

const (
	DilemmaTableName     = "dilemmas"
	DilemmaNodeTableName = "dilemma_nodes"
	NodeChildrenTable    = "node_childrens"
)

// DilemmaEntity — GORM-сущность для таблицы dilemmas.
type DilemmaEntity struct {
	ID         uuid.UUID `gorm:"primaryKey;column:id;type:uuid"`
	OwnerID    uuid.UUID `gorm:"column:owner_id;type:uuid;not null"`
	Topic      string    `gorm:"column:topic;type:varchar(256);not null"`
	RootNodeID uuid.UUID `gorm:"column:root_node_id;type:uuid;not null;uniqueIndex"`

	RootNode *DilemmaNodeEntity `gorm:"foreignKey:ID;references:RootNodeID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}

func (*DilemmaEntity) TableName() string {
	return DilemmaTableName
}

// ToModel преобразует GORM-сущность в доменную.
func (e *DilemmaEntity) ToModel() *dilemma_entity.Dilemma {
	return &dilemma_entity.Dilemma{
		ID:       e.ID,
		OwnerID:  e.OwnerID,
		Topic:    e.Topic,
		RootNode: e.RootNode.ToModel(),
	}
}

// DilemmaEntityFromModel создаёт GORM-сущность из доменной.
func DilemmaEntityFromModel(d *dilemma_entity.Dilemma) *DilemmaEntity {
	rootNodeEnt := DilemmaNodeEntityFromModel(d.RootNode)

	return &DilemmaEntity{
		ID:         d.ID,
		OwnerID:    d.OwnerID,
		Topic:      d.Topic,
		RootNodeID: d.RootNode.ID,
		RootNode:   rootNodeEnt,
	}
}

// DilemmaNodeEntity — GORM-сущность для таблицы dilemma_nodes.
type DilemmaNodeEntity struct {
	ID       uuid.UUID `gorm:"primaryKey;column:id;type:uuid"`
	Name     string    `gorm:"column:name;type:text;not null"`
	Value    string    `gorm:"column:value;type:text;not null"`
	ParentID uuid.UUID `gorm:"-"`

	// Связь "один ко многим" через join-таблицу node_children
	Children []*DilemmaNodeEntity `gorm:"many2many:node_childrens;foreignKey:ID;references:ID;joinForeignKey:node_id;joinReferences:child_id;constraint:OnDelete:CASCADE"`
}

type NodeChildren struct {
	NodeID  uuid.UUID `gorm:"primaryKey;column:node_id;type:uuid"`
	ChildID uuid.UUID `gorm:"primaryKey;column:child_id;type:uuid"`
}

func (*DilemmaNodeEntity) TableName() string {
	return DilemmaNodeTableName
}

func (e *DilemmaNodeEntity) ToModel() *dilemma_entity.DilemmaNode {
	node := &dilemma_entity.DilemmaNode{
		ID:       e.ID,
		Name:     e.Name,
		Value:    e.Value,
		ParentID: e.ParentID,
	}

	if len(e.Children) > 0 {
		node.Scenarios = make([]*dilemma_entity.Scenario, len(e.Children))
		for i, child := range e.Children {
			node.Scenarios[i] = &dilemma_entity.Scenario{
				ID:   child.ID,
				Name: child.Name,
			}
		}
	}

	return node
}

func DilemmaNodeEntityFromModel(n *dilemma_entity.DilemmaNode) *DilemmaNodeEntity {
	return &DilemmaNodeEntity{
		ID:    n.ID,
		Name:  n.Name,
		Value: n.Value,
	}
}

func ToModelList(entities []*DilemmaEntity) []dilemma_entity.Dilemma {
	result := make([]dilemma_entity.Dilemma, 0, len(entities))
	for _, e := range entities {
		result = append(result, *e.ToModel())
	}

	return result
}
