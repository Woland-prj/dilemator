package dilemma_entity

import "github.com/google/uuid"

type Dilemma struct {
	ID       uuid.UUID
	OwnerID  uuid.UUID
	Topic    string
	RootNode *DilemmaNode
}

type DilemmaNode struct {
	ID    uuid.UUID
	Value string
}

func NewDilemma(id, ownerID uuid.UUID, topic string, rootNode *DilemmaNode) *Dilemma {
	return &Dilemma{ID: id, OwnerID: ownerID, Topic: topic, RootNode: rootNode}
}

func NewDilemmaNode(id uuid.UUID, value string) *DilemmaNode {
	return &DilemmaNode{ID: id, Value: value}
}
