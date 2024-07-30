package apitest

import (
	"github.com/gin-gonic/gin"
	"github.com/nhannguyenacademy/ecommerce/internal/sdkapp/auth"
	"github.com/nhannguyenacademy/ecommerce/internal/sdkapp/mid"
	"github.com/nhannguyenacademy/ecommerce/internal/sdkbus/dbtest"
	"github.com/nhannguyenacademy/ecommerce/internal/user/userapp"
	"testing"
)

// todo: using github.com/stretchr/testify/suite

// StartTest initialized the system to run a test.
func StartTest(t *testing.T, testName string) *Test {
	db := dbtest.NewDatabase(t, testName)
	log := db.Log

	// -------------------------------------------------------------------------

	ath, err := auth.New(auth.Config{
		Log:       log,
		DB:        db.DB,
		KeyLookup: &KeyStore{},
	})
	if err != nil {
		t.Fatal(err)
	}

	gin.SetMode(gin.ReleaseMode)
	ginEngine := gin.New()
	v1Router := ginEngine.Group("v1")
	v1Router.Use(mid.Logging(log), mid.Panic(log))
	userapp.Routes(v1Router, userapp.Config{
		Log:     log,
		UserBus: db.BusDomain.User,
		Auth:    ath,
	})

	// -------------------------------------------------------------------------

	return New(db, ath, ginEngine.Handler())
}
