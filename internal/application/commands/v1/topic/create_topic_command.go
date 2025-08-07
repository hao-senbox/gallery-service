package topic

import (
	"gallery-service/internal/domain/models"
)

type CreateTopicCommand struct {
	TopicName      string
	Title          string
	Note           string
	Image          []models.TopicImageConfig
	LanguageConfig []models.TopicLanguageConfig
	FolderID       string
}

func NewCreateTopicCommand(
	topicName string,
	title string,
	note string,
	image []models.TopicImageConfig,
	languageConfig []models.TopicLanguageConfig,
	folderID string,
) *CreateTopicCommand {
	return &CreateTopicCommand{
		TopicName:      topicName,
		Title:          title,
		Note:           note,
		Image:          image,
		LanguageConfig: languageConfig,
		FolderID:       folderID,
	}
}
