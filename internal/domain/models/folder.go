package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Folder struct {
	ID                 primitive.ObjectID  `json:"id" bson:"_id,omitempty"`
	FolderName         string              `json:"folder_name" bson:"folder_name,omitempty"`
	FolderThumbnailKey string              `json:"folder_thumbnail_key" bson:"folder_thumbnail_key,omitempty"`
	FolderThumbnailURL string              `json:"folder_thumbnail_url" bson:"folder_thumbnail_url,omitempty"`
	ParentID           *primitive.ObjectID `json:"parent_id" bson:"parent_id,omitempty"`
}

// GetName returns the name of the cluster
func (c Folder) GetName() string {
	return "cluster"
}
