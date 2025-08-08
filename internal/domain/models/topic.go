package models

import (
	"gallery-service/internal/pkg/constants"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type TopicImageConfig struct {
	PicName   string `json:"pic_name" bson:"pic_name,omitempty"`
	ImageKey  string `json:"image_key" bson:"image_key,omitempty"`
	ImageURL  string `json:"image_url" bson:"image_url,omitempty"`
	OnlineURL string `json:"online_url" bson:"online_url,omitempty"`
}

type TopicVideoConfig struct {
	VideoName string `json:"video_name" bson:"video_name,omitempty"`
	VideoKey  string `json:"video_key" bson:"video_key,omitempty"`
	VideoURL  string `json:"video_url" bson:"video_url,omitempty"`
	OnlineURL string `json:"online_url" bson:"online_url,omitempty"`
	StartTime string `json:"start_time" bson:"start_time,omitempty"`
	EndTime   string `json:"end_time" bson:"end_time,omitempty"`
}

type TopicAudioConfig struct {
	AudioName string `json:"audio_name" bson:"audio_name,omitempty"`
	AudioKey  string `json:"audio_key" bson:"audio_key,omitempty"`
	AudioURL  string `json:"audio_url" bson:"audio_url,omitempty"`
	OnlineURL string `json:"online_url" bson:"online_url,omitempty"`
	StartTime string `json:"start_time" bson:"start_time,omitempty"`
	EndTime   string `json:"end_time" bson:"end_time,omitempty"`
}

type TopicLanguageConfig struct {
	Language    constants.Language `json:"language" bson:"language,omitempty"`
	Component   string             `json:"component" bson:"component,omitempty"`
	Title       string             `json:"title" bson:"title,omitempty"`
	Note        string             `json:"note" bson:"note,omitempty"`
	Description string             `json:"description" bson:"description,omitempty"`
	Images      []TopicImageConfig `json:"images" bson:"images,omitempty"`
	Videos      []TopicVideoConfig `json:"videos" bson:"videos,omitempty"`
	Audios      []TopicAudioConfig `json:"audios" bson:"audios,omitempty"`
}

type Topic struct {
	ID             primitive.ObjectID    `json:"id" bson:"_id,omitempty"`
	FileName       string                `json:"file_name" bson:"file_name,omitempty"`
	IsPublished    bool                  `json:"is_published" bson:"is_published,omitempty"`
	LanguageConfig []TopicLanguageConfig `json:"language_config" bson:"language_config,omitempty"`
	CreatedAt      time.Time             `json:"created_at" bson:"created_at,omitempty"`
	UpdatedAt      time.Time             `json:"updated_at" bson:"updated_at,omitempty"`
}

// GetName returns the name of the gallery
func (c Topic) GetName() string {
	return "topic gallery"
}
