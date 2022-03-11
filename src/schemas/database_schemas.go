package schemas

type DbUser struct {
	Id                  uint
	Username            string
	Email               string
	Password            string
	RegistrationDate    int64
	Rank                uint8
	LastPasswordRefresh int64
	Friends             []uint
	Settings            UserSettingsSchema
	Stats               UserStatsSchema
}

func (user DbUser) ToPublic() DbUserPublic {
	return DbUserPublic{
		Id:       user.Id,
		Username: user.Username,
	}
}

type DbUserPublic struct {
	Id       uint
	Username string
}

type DbAwaitUser struct {
	Id               uint
	Username         string
	Email            string
	Password         string
	RegistrationDate int64
}

type DbCategory struct {
	Id           uint
	Name         string
	Slug         string
	CreationDate int64
	Topics       []uint
	Icon         string
}

type DbTopic struct {
	Id           uint
	ParentId     uint
	Name         string
	Slug         string
	AuthorId     uint
	CreationDate int64
	Posts        []uint
	Pinned       bool
	PinnedBy     uint
	Locked       bool
	LockDate     int64
	Tags         []string
}

type DbPost struct {
	Id           uint
	ParentId     uint
	Content      string
	AuthorId     uint
	CreationDate int64
	Likes        []uint
	Dislikes     []uint
}

type DbReport struct {
	Id           uint
	ReportType   string
	TargetId     uint
	CreationDate int64
	AuthorId     uint
	Messages     []uint
}

type DbReportMessage struct {
	Id           uint
	ParentId     uint
	CreationDate int64
	AuthorId     uint
	Content      string
}
