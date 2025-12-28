package file_errors

import "errors"

var (
	ErrUploadingFailed = errors.New("uploading failed")
	ErrRemovingFailed  = errors.New("removing failed")
	ErrFileNotFound    = errors.New("file not found")
)
