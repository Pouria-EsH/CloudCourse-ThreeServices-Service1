package storage

import (
	"mime/multipart"
)

type ImageStorage interface {
	Upload(imageFile multipart.File, size int64, key string) (string, error)
}
