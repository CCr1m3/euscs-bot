package db

import (
	"testing"
)

func Test_db_migrate(t *testing.T) {
	Clear()
	GetInstance()
	if err := migrate(); err != nil {
		t.Errorf("migrate() error: " + err.Error())
	}
}

func Test_db_getLatestMigration(t *testing.T) {
	Clear()
	GetInstance()
	_, err := getLatestMigration()
	if err == nil {
		t.Error("getLatestMigration() should be in error, no table migrations yet", err)
	}
	Init()
	index, err := getLatestMigration()
	if err != nil {
		t.Error("getLatestMigration() is in error: " + err.Error())
	}
	if index != len(migrations)-1 {
		t.Errorf("failed to get latest migration: got:%v , want:%v ", index, len(migrations)-1)
	}
}
