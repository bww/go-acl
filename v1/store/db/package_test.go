package db

import (
	"os"
	"testing"

	"github.com/bww/go-dbx/v1/test"
	"github.com/bww/go-util/env"
	"github.com/bww/go-util/urls"
)

func TestMain(m *testing.M) {
	test.Init("acl_v1_store_db_test", test.WithMigrations(urls.File(env.Etc("db"))))
	os.Exit(m.Run())
}
