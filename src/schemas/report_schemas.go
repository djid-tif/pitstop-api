package schemas

type ReportSchema struct {
	Id           uint   `json:"id"`
	ReportType   string `json:"report_type"`
	TargetId     uint   `json:"target_id"`
	CreationDate string `json:"creation_date"`
	AuthorId     uint   `json:"author_id"`
	Messages     []uint `json:"messages"`
}

type ReportMessageSchema struct {
	Id           uint   `json:"id"`
	ParentId     uint   `json:"parent_id"`
	CreationDate string `json:"creation_date"`
	AuthorId     uint   `json:"author_id"`
	Content      string `json:"content"`
}
