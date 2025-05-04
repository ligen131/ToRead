package model

import (
	"time"
	"to-read/utils"
	"to-read/utils/logs"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

type User struct {
	ID          uint32         `json:"user_id"      form:"user_id"      query:"user_id"      gorm:"primaryKey;unique;autoIncrement;not null"`
	CreatedAt   time.Time      `json:"created_at"   form:"created_at"   query:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"   form:"updated_at"   query:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"deleted_at"   form:"deleted_at"   query:"deleted_at"`
	UserName    string         `json:"user_name"    form:"user_name"    query:"user_name"    gorm:"unique;not null"`
	Role        uint32         `json:"role"         form:"role"         query:"role"         gorm:"not null"`
	PasswordMD5 string         `json:"password_md5" form:"password_md5" query:"password_md5" gorm:"not null"`
	Deleted     bool           `json:"deleted"      form:"deleted"      query:"deleted"      gorm:"not null"`
}

func FindMaxUserID() (User, error) {
	m := GetModel()
	defer m.Close()

	var user User
	result := m.tx.Order("id desc").First(&user)
	if result.Error != nil {
		logs.Info("Find max user id failed.", zap.Error(result.Error))
		m.Abort()
		return user, result.Error
	}

	m.tx.Commit()
	return user, nil
}

func UserRegister(userName string, password string) (User, error) {
	m := GetModel()
	defer m.Close()

	user := User{
		UserName:    userName,
		Role:        1,
		PasswordMD5: utils.GetMD5(password),
		Deleted:     false,
	}
	result := m.tx.Create(&user)
	if result.Error != nil {
		logs.Warn("Create user failed.", zap.Error(result.Error), zap.Any("user", user))
		m.Abort()
		return user, result.Error
	}

	m.tx.Commit()
	return user, nil
}

func FindUserByID(userID uint32) (User, error) {
	m := GetModel()
	defer m.Close()

	var user User
	result := m.tx.First(&user, userID)
	if result.Error != nil {
		logs.Info("Find user by id failed.", zap.Error(result.Error))
		m.Abort()
		return user, result.Error
	}

	m.tx.Commit()
	return user, nil
}

func FindUserByName(userName string) (User, error) {
	m := GetModel()
	defer m.Close()

	var user User
	result := m.tx.Model(&User{}).Where("user_name = ?", userName).First(&user)
	if result.Error != nil {
		logs.Info("Find user by name failed.", zap.Error(result.Error))
		m.Abort()
		return user, result.Error
	}

	m.tx.Commit()
	return user, nil
}

func UpdateUserName(userID uint32, userName string) (User, error) {
	m := GetModel()
	defer m.Close()

	var user User
	result := m.tx.First(&user, userID).Update("user_name", userName)
	if result.Error != nil {
		logs.Info("Update user name failed.", zap.Error(result.Error))
		m.Abort()
		return user, result.Error
	}

	m.tx.Commit()
	return user, nil
}
