package dilemma_errors

import "errors"

var (
	ErrDilemmaAlreadyExists = errors.New("dilemma_already_exists")
	ErrDilemmaNotFound      = errors.New("dilemma_router not found")
	ErrNodeNotFound         = errors.New("node not found")
	ErrNodeAlreadyExists    = errors.New("dilemma_router node already exists")
	ErrNodeAlreadyHasParent = errors.New("dilemma_router node has parent")
)
