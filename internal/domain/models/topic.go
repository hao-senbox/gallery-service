package models

import (
	"gallery-service/internal/pkg/constants"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type TopicVideoConfig struct {
	VideoKey  string `json:"video_key" bson:"video_key,omitempty"`
	VideoURL  string `json:"video_url" bson:"video_url,omitempty"`
	StartTime string `json:"start_time" bson:"start_time,omitempty"`
	EndTime   string `json:"end_time" bson:"end_time,omitempty"`
}

type TopicAudioConfig struct {
	AudioKey  string `json:"audio_key" bson:"audio_key,omitempty"`
	AudioURL  string `json:"audio_url" bson:"audio_url,omitempty"`
	StartTime string `json:"start_time" bson:"start_time,omitempty"`
	EndTime   string `json:"end_time" bson:"end_time,omitempty"`
}

type TopicImageConfig struct {
	ImageKey string `json:"image_key" bson:"image_key,omitempty"`
	ImageURL string `json:"image_url" bson:"image_url,omitempty"`
}

type TopicLanguageConfig struct {
	Language constants.Language `json:"language" bson:"language,omitempty"`
	Video    VideoConfig        `json:"video" bson:"video,omitempty"`
	Audio    AudioConfig        `json:"audio" bson:"audio,omitempty"`
}

type Topic struct {
	ID             primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	TopicName      string             `json:"topic_name" bson:"topic_name,omitempty"`
	Title          string             `json:"title" bson:"title,omitempty"`
	Note           string             `json:"note" bson:"note,omitempty"`
	Image          []ImageConfig      `json:"image" bson:"image,omitempty"`
	LanguageConfig []LanguageConfig   `json:"language_config" bson:"language_config,omitempty"`
	FolderID       primitive.ObjectID `json:"folder_id" bson:"folder_id,omitempty"`
	CreatedAt      time.Time          `json:"created_at" bson:"created_at,omitempty"`
	UpdatedAt      time.Time          `json:"updated_at" bson:"updated_at,omitempty"`
}

// GetName returns the name of the gallery
func (c Topic) GetName() string {
	return "topic gallery"
}
