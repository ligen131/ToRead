package model

import (
	"time"
	"to-read/utils/logs"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

type CollectionWithTag struct {
	ID           uint32         `json:"collection_with_tag_id"      form:"collection_with_tag_id"      query:"collection_with_tag_id"      gorm:"primaryKey;unique;autoIncrement;not null"`
	CreatedAt    time.Time      `json:"created_at"                  form:"created_at"                  query:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"                  form:"updated_at"                  query:"updated_at"`
	DeletedAt    gorm.DeletedAt `json:"deleted_at"                  form:"deleted_at"                  query:"deleted_at"`
	CollectionID uint32         `json:"collection_id"               form:"collection_id"               query:"collection_id"               gorm:"not null;uniqueIndex:idx_collection_tag"`
	TagID        uint32         `json:"tag_id"                      form:"tag_id"                      query:"tag_id"                      gorm:"not null;uniqueIndex:idx_collection_tag"`
}

func FindCollectionWithTagID(collectionID uint32, tagID uint32) (CollectionWithTag, error) {
	m := GetModel()
	defer m.Close()

	cwt := CollectionWithTag{}
	result := m.tx.Where("collection_id = ? AND tag_id = ?", collectionID, tagID).First(&cwt)
	if result.Error != nil {
		logs.Info("Find collection with tag failed.", zap.Error(result.Error), zap.Any("cwt", cwt))
		m.Abort()
		return cwt, result.Error
	}

	m.tx.Commit()
	return cwt, nil
}

func FindOrCreateCollectionWithTagID(collectionID uint32, tagID uint32) (CollectionWithTag, error) {
	if cwt, err := FindCollectionWithTagID(collectionID, tagID); err == nil {
		return cwt, nil
	} else if err != gorm.ErrRecordNotFound {
		return cwt, err
	}

	m := GetModel()
	defer m.Close()

	// collection_with_tag
	cwt := CollectionWithTag{
		CollectionID: collectionID,
		TagID:        tagID,
	}
	result := m.tx.Create(&cwt)
	if result.Error != nil {
		logs.Info("Create collection with tag failed.", zap.Error(result.Error), zap.Any("cwt", cwt))
		m.Abort()
		return cwt, result.Error
	}

	m.tx.Commit()
	return cwt, nil
}

func AddCollectionWithTags(collectionID uint32, tags []string) ([]CollectionWithTag, error) {
	result := make([]CollectionWithTag, 0, len(tags))

	for _, tagName := range tags {
		tag, err := FindOrCreateTag(tagName)
		if err != nil {
			return result, err
		}

		logs.Debug("AddCollectionWithTags: FindOrCreateTag", zap.Uint32("collectionID", collectionID), zap.Uint32("tagID", tag.ID), zap.String("tagName", tagName))
		cwt, err := FindOrCreateCollectionWithTagID(collectionID, tag.ID)
		if err == nil {
			result = append(result, cwt)
		} else {
			return result, err
		}
	}
	return result, nil
}


// 获取收藏关联的标签
func GetCollectionTags(collectionID uint32) ([]string, error) {
	m := GetModel()
	defer m.Close()

	var tagNames []string
	
	// 查询与收藏关联的所有标签
	err := m.tx.Model(&Tag{}).
		Select("tags.name").
		Joins("JOIN collection_with_tags ON collection_with_tags.tag_id = tags.id").
		Where("collection_with_tags.collection_id = ?", collectionID).
		Pluck("name", &tagNames).Error
	
	if err != nil {
		m.Abort()
		return nil, err
	}

	m.tx.Commit()
	return tagNames, nil
}
