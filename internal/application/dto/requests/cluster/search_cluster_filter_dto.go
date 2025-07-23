package cluster

type SearchClusterFilterReqDto struct {
	Keyword string `json:"keyword,omitempty" validate:"required"`
}
