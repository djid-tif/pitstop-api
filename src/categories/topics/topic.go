package topics

import (
	"encoding/json"
	"errors"
	"pitstop-api/src/categories/topics/posts"
	"pitstop-api/src/database"
	"pitstop-api/src/schemas"
	"pitstop-api/src/utils"
	"time"
)

type Topic struct {
	id           uint
	parentId     uint
	name         string
	authorId     uint
	creationDate time.Time
	posts        []uint
	pinned       bool
	pinnedBy     uint
	locked       bool
	lockDate     time.Time
	tags         []string
}

func (topic Topic) GetId() uint {
	return topic.id
}

func (topic Topic) GetParentId() uint {
	return topic.parentId
}

func (topic Topic) GetName() string {
	return topic.name
}

func (topic Topic) GetAuthorId() uint {
	return topic.authorId
}

func (topic Topic) GetCreationDate() time.Time {
	return topic.creationDate
}

func (topic Topic) GetPosts() []uint {
	return topic.posts
}

func (topic Topic) IsPinned() bool {
	return topic.pinned
}

func (topic Topic) GetPinnedBy() uint {
	return topic.pinnedBy
}

func (topic Topic) IsLocked() bool {
	return topic.locked
}

func (topic Topic) GetLockDate() time.Time {
	return topic.lockDate
}

func (topic Topic) GetTags() []string {
	return topic.tags
}

func (topic Topic) ToSchema() schemas.TopicSchema {
	return schemas.TopicSchema{
		Id:           topic.id,
		ParentId:     topic.parentId,
		Name:         topic.name,
		AuthorId:     topic.authorId,
		CreationDate: utils.TimeToIso(topic.creationDate),
		Posts:        topic.posts,
		Pinned:       topic.pinned,
		PinnedBy:     topic.pinnedBy,
		Locked:       topic.locked,
		LockDate:     utils.TimeToIso(topic.lockDate),
		Tags:         topic.tags,
	}
}

func (topic *Topic) appendChild(postId uint) error {
	if utils.ContainsUint(postId, topic.posts) {
		return errors.New("un post avec cet id existe déjà dans ce topic")
	}

	topic.posts = append(topic.posts, postId)
	jsonChildren, err := json.Marshal(topic.posts)
	if err != nil {
		return err
	}

	return database.UpdateTopicById(topic.id, "posts", jsonChildren)
}

func (topic *Topic) Delete(deleteFromParent bool) error {
	for _, postId := range topic.posts {
		post, found := posts.LoadFromId(postId)
		if !found {
			continue
		}
		err := post.Delete(false)
		if err != nil {
			return err
		}
	}

	return database.DeleteTopic(topic.id, deleteFromParent)
}

func LoadFromId(id uint) (Topic, bool) {
	topicSchema, found := database.GetTopicById(id)
	if !found {
		return Topic{}, false
	}
	return Topic{
		id:           topicSchema.Id,
		parentId:     topicSchema.ParentId,
		name:         topicSchema.Name,
		authorId:     topicSchema.AuthorId,
		creationDate: time.Unix(topicSchema.CreationDate, 0),
		posts:        topicSchema.Posts,
		pinned:       topicSchema.Pinned,
		pinnedBy:     topicSchema.PinnedBy,
		locked:       topicSchema.Locked,
		lockDate:     time.Unix(topicSchema.LockDate, 0),
		tags:         topicSchema.Tags,
	}, true
}
