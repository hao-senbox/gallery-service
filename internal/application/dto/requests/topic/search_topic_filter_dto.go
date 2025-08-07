package topic

type SearchTopicFilterReqDto struct {
	Keyword string `json:"keyword,omitempty" validate:"required"`
}
