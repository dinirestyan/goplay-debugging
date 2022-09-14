package models

import (
	"mime/multipart"

	uuid "github.com/satori/go.uuid"
)

type Files struct {
	BaseModel
	FileName     string    `json:"file_name"`
	Uploader     uuid.UUID `json:"-"`
	UploaderName Uploader  `gorm:"foreignKey:Uploader"`
}

func (Files) TableName() string {
	return "files"
}

type Uploader struct {
	ID   uuid.UUID `json:"id" gorm:"type:uuid;primary_key;"`
	Name string    `json:"name"`
}

func (Uploader) TableName() string {
	return "users"
}

type Upload struct {
	Documents []*multipart.FileHeader `form:"documents" json:"documents"`
}
