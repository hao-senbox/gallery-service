package topic

import "gallery-service/internal/domain/models"

type UpdateTopicCommand struct {
	ID             string
	TopicName      string
	Title          string
	Note           string
	Image          []models.TopicImageConfig
	LanguageConfig []models.TopicLanguageConfig
	FolderID       string
}

func NewUpdateTopicCommand(
	id string,
	topicName string,
	title string,
	note string,
	image []models.TopicImageConfig,
	languageConfig []models.TopicLanguageConfig,
) *UpdateTopicCommand {
	return &UpdateTopicCommand{
		ID:             id,
		TopicName:      topicName,
		Title:          title,
		Note:           note,
		Image:          image,
		LanguageConfig: languageConfig,
	}
}
