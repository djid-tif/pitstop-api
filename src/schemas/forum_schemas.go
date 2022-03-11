package schemas

type CategorySchema struct {
	Id           uint   `json:"id"`
	Name         string `json:"name"`
	CreationDate string `json:"creation_date"`
	Icon         string `json:"icon"`
	Topics       []uint `json:"topics"`
	LastTopic    uint   `json:"last_topic"`
}

type TopicSchema struct {
	Id           uint     `json:"id"`
	ParentId     uint     `json:"parent_id"`
	Name         string   `json:"name"`
	AuthorId     uint     `json:"author_id"`
	CreationDate string   `json:"creation_date"`
	Posts        []uint   `json:"posts"`
	Pinned       bool     `json:"pinned"`
	PinnedBy     uint     `json:"pinned_by"`
	Locked       bool     `json:"locked"`
	LockDate     string   `json:"lock_date"`
	Tags         []string `json:"tags"`
}

type PostSchema struct {
	Id           uint   `json:"id"`
	ParentId     uint   `json:"parent_id"`
	Content      string `json:"content"`
	AuthorId     uint   `json:"author_id"`
	CreationDate string `json:"creation_date"`
	Likes        []uint `json:"likes"`
	Dislikes     []uint `json:"dislikes"`
}
