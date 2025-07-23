package folder

type SearchFolderFilterReqDto struct {
	Keyword string `json:"keyword,omitempty" validate:"required"`
}
