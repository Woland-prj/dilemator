package dilemma_entity

import "github.com/google/uuid"

type Dilemma struct {
	ID       uuid.UUID
	OwnerID  uuid.UUID
	Topic    string
	RootNode *DilemmaNode
}

type DilemmaNode struct {
	ID       uuid.UUID
	ParentID uuid.UUID
	Name     string
	Value    string
	Scenrios []*Scenario
}

type Scenario struct {
	Name string
	Id   uuid.UUID
}

func NewDilemma(id, ownerID uuid.UUID, topic string, rootNode *DilemmaNode) *Dilemma {
	return &Dilemma{ID: id, OwnerID: ownerID, Topic: topic, RootNode: rootNode}
}

func NewDilemmaNode(id, pid uuid.UUID, name, value string) *DilemmaNode {
	return &DilemmaNode{ID: id, ParentID: pid, Name: name, Value: value}
}

func NewEmptyNode(pid uuid.UUID) *DilemmaNode {
	return &DilemmaNode{
		ID:       uuid.Nil,
		ParentID: pid,
		Name:     "",
		Value:    "",
	}
}
