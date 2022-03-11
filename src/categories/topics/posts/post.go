package posts

import (
	"pitstop-api/src/database"
	"pitstop-api/src/schemas"
	"pitstop-api/src/utils"
	"time"
)

type Post struct {
	id           uint
	parentId     uint
	content      string
	authorId     uint
	creationDate time.Time
	likes        []uint
	dislikes     []uint
}

func (post Post) GetId() uint {
	return post.id
}

func (post Post) GetParentId() uint {
	return post.parentId
}

func (post Post) GetContent() string {
	return post.content
}

func (post Post) GetAuthorId() uint {
	return post.authorId
}

func (post Post) GetCreationDate() time.Time {
	return post.creationDate
}

func (post Post) GetLikes() []uint {
	return post.likes
}

func (post Post) GetDislikes() []uint {
	return post.dislikes
}

func (post Post) ToSchema() schemas.PostSchema {
	return schemas.PostSchema{
		Id:           post.id,
		ParentId:     post.parentId,
		Content:      post.content,
		AuthorId:     post.authorId,
		CreationDate: utils.TimeToIso(post.creationDate),
		Likes:        post.likes,
		Dislikes:     post.dislikes,
	}
}

func (post Post) Delete(deleteFromParent bool) error {
	return database.DeletePost(post.id, deleteFromParent)
}

func LoadFromId(id uint) (Post, bool) {
	post, found := database.GetPostById(id)
	if !found {
		return Post{}, false
	}
	return Post{
		id:           post.Id,
		parentId:     post.ParentId,
		content:      post.Content,
		authorId:     post.AuthorId,
		creationDate: time.Unix(post.CreationDate, 0),
		likes:        post.Likes,
		dislikes:     post.Dislikes,
	}, true
}
