package model

import (
	"time"
	"to-read/utils/logs"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

type Collection struct {
	ID          uint32         `json:"collection_id"      form:"collection_id"      query:"collection_id"      gorm:"primaryKey;unique;autoIncrement;not null"`
	CreatedAt   time.Time      `json:"created_at"         form:"created_at"         query:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"         form:"updated_at"         query:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"deleted_at"         form:"deleted_at"         query:"deleted_at"`
	UserID      uint32         `json:"user_id"            form:"user_id"            query:"user_id"            gorm:"not null;uniqueIndex:idx_user_url"`
	Url         string         `json:"url"                form:"url"                query:"url"                gorm:"not null;uniqueIndex:idx_user_url"`
	Type        string         `json:"type"               form:"type"               query:"type"` // one of text/image/video
	Title       string         `json:"title"              form:"title"              query:"title"              gorm:"type text"`
	Description string         `json:"description"        form:"description"        query:"description"        gorm:"type text"`
}

func FindCollectionByUrl(userID uint32, url string) (Collection, error) {
	m := GetModel()
	defer m.Close()

	var collection Collection
	result := m.tx.Where("user_id = ? AND url = ?", userID, url).First(&collection)
	if result.Error != nil {
		logs.Info("Find collection by URL failed.", zap.Error(result.Error), zap.String("url", url))
		m.Abort()
		return collection, result.Error
	}

	m.tx.Commit()
	return collection, nil
}

func AddCollection(userID uint32, url string, type_ string, title string, description string, tags []string) (Collection, []CollectionWithTag, error) {
	m := GetModel()
	defer m.Close()

	collection := Collection{
		UserID:      userID,
		Url:         url,
		Type:        type_,
		Title:       title,
		Description: description,
	}
	cwt := make([]CollectionWithTag, 0, len(tags))
	result := m.tx.Create(&collection)
	if result.Error != nil {
		logs.Info("Create collection failed.", zap.Error(result.Error), zap.Any("collection", collection))
		m.Abort()
		return collection, cwt, result.Error
	}

	cwt, err := AddCollectionWithTags(collection.ID, tags)
	if err != nil {
		logs.Info("Add collection with tags failed.", zap.Error(result.Error), zap.Any("collection", collection))
		m.Abort()
		return collection, cwt, result.Error
	}

	m.tx.Commit()
	return collection, cwt, nil
}

// SearchCollection 根据用户ID、关键词和标签列表搜索收藏
func SearchCollection(userID uint32, keyword string, tags []string) ([]Collection, error) {
	m := GetModel()
	defer m.Close()

	// 构建基本查询
	query := m.tx.Model(&Collection{}).Where("collections.user_id = ?", userID)

	// 关键词过滤
	if keyword != "" {
		query = query.Where("(collections.title ILIKE ? OR collections.description ILIKE ?)",
			"%"+keyword+"%", "%"+keyword+"%")
	}

	// 标签过滤
	if len(tags) > 0 {
		// 先查找所有标签的ID
		var tagIDs []uint32
		for _, tagName := range tags {
			tag, err := FindTagByName(tagName)
			if err != nil {
				if err == gorm.ErrRecordNotFound {
					// 标签不存在，返回空结果
					return []Collection{}, nil
				}
				logs.Warn("Find tag by name failed", zap.Error(err), zap.String("tag_name", tagName))
				m.Abort()
				return nil, err
			}
			tagIDs = append(tagIDs, tag.ID)
		}

		// 对每个标签ID，我们需要确保收藏与之关联
		for _, tagID := range tagIDs {
			// 使用子查询确保收藏包含当前标签
			subQuery := m.tx.Model(&CollectionWithTag{}).
				Select("collection_id").
				Where("tag_id = ?", tagID)

			query = query.Where("collections.id IN (?)", subQuery)
		}
	}

	// 执行查询并去重
	var collections []Collection
	result := query.Distinct().Find(&collections)
	if result.Error != nil {
		logs.Error("Search collections failed",
			zap.Error(result.Error),
			zap.Uint32("user_id", userID),
			zap.String("keyword", keyword),
			zap.Any("tags", tags))
		m.Abort()
		return nil, result.Error
	}

	m.tx.Commit()
	return collections, nil
}
