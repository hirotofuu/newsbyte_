package main

import (
	"os"
	"testing"
	"time"

	"github.com/hirotofuu/newsbyte/internal/repository/dbrepo"
)

// テスト環境

var app application

func TestMain(m *testing.M) {
	// ダミー関数のリポジトリ
	app.DB = &dbrepo.TestDBRepo{}
	app.ADB = &dbrepo.TestADBRepo{}
	app.CDB = &dbrepo.TestCDBRepo{}

	app.Domain = "example.com"
	app.JWTSecret = "x6ur8ec6t5hf4kzjy7tjzdatkxbcwfsbufydcktkdpkghpek6922g65z9eyb6twj7aweb24sjpiiw33e5xtanhdecd25yfa8bpdr"
	app.JWTIssuer = "newsbyte_master_dayo"
	app.JWTAudience = "example.com"
	app.CookieDomain = "localhost"
	app.Domain = "example.com"
	app.APIKey = "aed5486efae48ec4cfa93ac93c4afeaf"

	app.auth = Auth{
		Issuer:        app.JWTIssuer,
		Audience:      app.JWTAudience,
		Secret:        app.JWTSecret,
		TokenExpiry:   time.Minute * 10,
		RefreshExpiry: time.Hour * 24,
		CookiePath:    "/",
		CookieName:    "refresh-token",
		CookieDomain:  app.CookieDomain,
	}
	os.Exit(m.Run())
}
