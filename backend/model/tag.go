package model

import (
	"time"
	"to-read/utils/logs"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

type Tag struct {
	ID        uint32         `json:"tag_id"      form:"tag_id"      query:"tag_id"      gorm:"primaryKey;unique;autoIncrement;not null"`
	CreatedAt time.Time      `json:"created_at"  form:"created_at"  query:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"  form:"updated_at"  query:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"deleted_at"  form:"deleted_at"  query:"deleted_at"`
	Name      string         `json:"tag_name"    form:"tag_name"    query:"tag_name"    gorm:"unique;not null"`
}

func FindTagByName(tagName string) (Tag, error) {
	m := GetModel()
	defer m.Close()

	var tag Tag
	result := m.tx.Model(&Tag{}).Where("name = ?", tagName).First(&tag)
	if result.Error != nil {
		logs.Info("Find tag by name failed.", zap.Error(result.Error))
		m.Abort()
		return tag, result.Error
	}

	m.tx.Commit()
	return tag, nil
}

func FindMaxTagID() (Tag, error) {
	m := GetModel()
	defer m.Close()

	var tag Tag
	result := m.tx.Order("id desc").First(&tag)
	if result.Error != nil {
		logs.Info("Find max tag id failed.", zap.Error(result.Error))
		m.Abort()
		return tag, result.Error
	}

	m.tx.Commit()
	return tag, nil
}

func FindOrCreateTag(tagName string) (Tag, error) {
	if tag, err := FindTagByName(tagName); err == nil {
		return tag, nil
	} else if err != gorm.ErrRecordNotFound {
		return Tag{}, err
	}

	m := GetModel()
	defer m.Close()

	tag := Tag{
		Name: tagName,
	}
	result := m.tx.Create(&tag)
	if result.Error != nil {
		logs.Warn("Create tag failed.", zap.Error(result.Error), zap.Any("tag", tag))
		m.Abort()
		return tag, result.Error
	}

	m.tx.Commit()
	return tag, nil
}

// GetTagList 获取指定用户的所有标签
func GetTagList(userID uint32) ([]Tag, error) {
	m := GetModel()
	defer m.Close()

	var tags []Tag

	// 使用子查询先获取用户的所有收藏ID
	collectionSubQuery := m.tx.Model(&Collection{}).
		Select("id").
		Where("user_id = ?", userID)

	// 再通过收藏ID获取所有关联的标签ID
	tagIDsSubQuery := m.tx.Model(&CollectionWithTag{}).
		Select("tag_id").
		Where("collection_id IN (?)", collectionSubQuery)

	// 最后获取这些标签的详细信息
	result := m.tx.Where("id IN (?)", tagIDsSubQuery).
		Order("name"). // 按标签名称排序
		Find(&tags)

	if result.Error != nil {
		logs.Error("Get tag list failed",
			zap.Error(result.Error),
			zap.Uint32("user_id", userID))
		m.Abort()
		return nil, result.Error
	}

	m.tx.Commit()
	return tags, nil
}
