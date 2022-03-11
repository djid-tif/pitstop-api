package categories

import (
	"encoding/json"
	"errors"
	"pitstop-api/src/categories/topics"
	"pitstop-api/src/database"
	"pitstop-api/src/schemas"
	"pitstop-api/src/utils"
	"time"
)

type Category struct {
	id           uint
	name         string
	creationDate time.Time
	topics       []uint
	icon         string
}

func (cat Category) GetId() uint {
	return cat.id
}

func (cat Category) GetName() string {
	return cat.name
}

func (cat Category) GetCreationDate() time.Time {
	return cat.creationDate
}

func (cat Category) GetIcon() string {
	return cat.icon
}

func (cat Category) GetTopics() []uint {
	return cat.topics
}

func (cat Category) GetLastTopic() uint {
	if len(cat.topics) == 0 {
		return 0
	}
	return database.GetLastTopic(cat.id)
}

func (cat Category) ToSchema() schemas.CategorySchema {
	return schemas.CategorySchema{
		Id:           cat.id,
		Name:         cat.name,
		CreationDate: cat.creationDate.String(),
		Icon:         cat.icon,
		Topics:       cat.topics,
		LastTopic:    cat.GetLastTopic(),
	}
}

func (cat *Category) appendChild(topicId uint) error {
	if utils.ContainsUint(topicId, cat.topics) {
		return errors.New("un topic avec cet id existe déjà dans cette catégorie")
	}

	cat.topics = append(cat.topics, topicId)
	jsonChildren, err := json.Marshal(cat.topics)
	if err != nil {
		return err
	}

	return database.UpdateCategoryById(cat.id, "topics", jsonChildren)
}

func (cat *Category) delete() error {
	for _, topicId := range cat.topics {
		topic, found := topics.LoadFromId(topicId)
		if !found {
			continue
		}
		err := topic.Delete(false)
		if err != nil {
			return err
		}
	}

	return database.DeleteCategory(cat.id)
}

func LoadFromId(id uint) (Category, bool) {
	categorySchema, found := database.GetCategory(id)
	if !found {
		return Category{}, false
	}
	return Category{
		id:           categorySchema.Id,
		name:         categorySchema.Name,
		creationDate: time.Unix(categorySchema.CreationDate, 0),
		topics:       categorySchema.Topics,
		icon:         categorySchema.Icon,
	}, true
}
