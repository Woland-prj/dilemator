package dilemma_errors

import "errors"

var (
	ErrDilemmaAlreadyExists = errors.New("dilemma_already_exists")
	ErrDilemmaNotFound      = errors.New("dilemma not found")
	ErrNodeNotFound         = errors.New("node not found")
	ErrNodeAlreadyExists    = errors.New("dilemma node already exists")
	ErrNodeAlreadyHasParent = errors.New("dilemma node has parent")
)
