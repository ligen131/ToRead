package model

import (
	"context"
	"time"

	"gorm.io/gorm"
)

var (
	db *gorm.DB
)

type Model struct {
	tx      *gorm.DB
	context context.Context
	cancel  context.CancelFunc
}

func GetModel() *Model {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	return &Model{
		tx:      db.Begin().WithContext(ctx),
		context: ctx,
		cancel:  cancel,
	}
}

func InitModel() error {
	err := AutoMigrateTable(&User{})
	if err != nil {
		return err
	}
	err = AutoMigrateTable(&Collection{})
	if err != nil {
		return err
	}
	err = AutoMigrateTable(&Tag{})
	if err != nil {
		return err
	}
	err = AutoMigrateTable(&CollectionWithTag{})
	if err != nil {
		return err
	}

	return nil
}

func (m *Model) Close() {
	if r := recover(); r != nil {
		m.tx.Rollback()
	}
	m.cancel()
}

func (m *Model) Abort() {
	m.tx.Rollback()
	m.cancel()
}
