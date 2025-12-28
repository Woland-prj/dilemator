package responses

import (
	"github.com/Woland-prj/dilemator/internal/domain/entity/dilemma_entity"
)

type DilemmaResponse struct {
	ID       string               `json:"id"`
	OwnerID  string               `json:"ownerId"`
	Topic    string               `json:"topic"`
	RootNode *DilemmaNodeResponse `json:"rootNode"`
}

type DilemmaNodeResponse struct {
	ID        string `json:"id"`
	ParentID  string `json:"parentId"`
	Name      string `json:"name"`
	Value     string `json:"value"`
	Scenarios []*ScenarioResponse
}

type ScenarioResponse struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

func NewScenarioResponse(scenario *dilemma_entity.Scenario) *ScenarioResponse {
	return &ScenarioResponse{
		ID:   scenario.ID.String(),
		Name: scenario.Name,
	}
}

func NewDilemmaNodeResponse(node *dilemma_entity.DilemmaNode) *DilemmaNodeResponse {
	scs := make([]*ScenarioResponse, len(node.Scenarios))
	for i, s := range node.Scenarios {
		scs[i] = NewScenarioResponse(s)
	}
	return &DilemmaNodeResponse{
		ID:        node.ID.String(),
		ParentID:  node.ParentID.String(),
		Name:      node.Name,
		Value:     node.Value,
		Scenarios: scs,
	}
}

func NewDilemmaResponse(dilemma dilemma_entity.Dilemma) *DilemmaResponse {
	return &DilemmaResponse{
		ID:       dilemma.ID.String(),
		OwnerID:  dilemma.OwnerID.String(),
		Topic:    dilemma.Topic,
		RootNode: NewDilemmaNodeResponse(dilemma.RootNode),
	}
}

func NewDilemmaResponseList(dilemma []dilemma_entity.Dilemma) []*DilemmaResponse {
	ds := make([]*DilemmaResponse, len(dilemma))
	for i, d := range dilemma {
		ds[i] = NewDilemmaResponse(d)
	}
	return ds
}
