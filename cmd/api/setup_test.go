package main

import (
	"os"
	"testing"

	"github.com/hirotofuu/newsbyte/internal/repository/dbrepo"
)

var app application

func TestMain(m *testing.M) {
	app.DB = &dbrepo.TestDBRepo{}
	app.ADB = &dbrepo.TestADBRepo{}
	app.Domain = "example.com"
	app.JWTSecret = "x6ur8ec6t5hf4kzjy7tjzdatkxbcwfsbufydcktkdpkghpek6922g65z9eyb6twj7aweb24sjpiiw33e5xtanhdecd25yfa8bpdr"
	os.Exit(m.Run())
}
