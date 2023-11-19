package gateway_test

import (
	"gorm.io/gorm"

	"github.com/pecolynx/golang-webapi-boilerplate/src/sqls"
	testlibgateway "github.com/pecolynx/golang-webapi-boilerplate/testlib/gateway"
)

func init() {
	fns := []func() (*gorm.DB, error){
		func() (*gorm.DB, error) {
			return testlibgateway.InitMySQL(sqls.SQL, "127.0.0.1", 3307)
		},
		// func() (*gorm.DB, error) {
		// 	return testlibgateway.InitSQLiteInFile(sqls.SQL)
		// },
	}

	for _, fn := range fns {
		db, err := fn()
		if err != nil {
			panic(err)
		}
		sqlDB, err := db.DB()
		if err != nil {
			panic(err)
		}
		sqlDB.Close()
	}
}
