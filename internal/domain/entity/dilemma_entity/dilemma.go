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

func NewDilemma(ID uuid.UUID, OwnerID uuid.UUID, Topic string, RootNode *DilemmaNode) *Dilemma {
	return &Dilemma{ID: ID, OwnerID: OwnerID, Topic: Topic, RootNode: RootNode}
}

func NewDilemmaNode(id uuid.UUID, Value string) *DilemmaNode {
	return &DilemmaNode{ID: id, Value: Value}
}
