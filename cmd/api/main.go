package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/hirotofuu/newsbyte/internal/repository"
	"github.com/hirotofuu/newsbyte/internal/repository/dbrepo"
)

type application struct {
	DSN          string
	DB           repository.Databaserepo
	ADB          repository.ArticleDatabaserepo
	Domain       string
	auth         Auth
	JWTSecret    string
	JWTIssuer    string
	JWTAudience  string
	CookieDomain string
	APIKey       string
}

const port = 8080

func main() {
	var app application

	flag.StringVar(&app.JWTSecret, "jwt-secret", "x6ur8ec6t5hf4kzjy7tjzdatkxbcwfsbufydcktkdpkghpek6922g65z9eyb6twj7aweb24sjpiiw33e5xtanhdecd25yfa8bpdr", "signing secret")
	flag.StringVar(&app.JWTIssuer, "jwt-issuer", "example.com", "signing issuer")
	flag.StringVar(&app.JWTAudience, "jwt-audience", "example.com", "signing audience")
	flag.StringVar(&app.CookieDomain, "cookie-domain", "localhost", "cookie domain")
	flag.StringVar(&app.Domain, "domain", "example.com", "domain")
	flag.StringVar(&app.APIKey, "api_key", "aed5486efae48ec4cfa93ac93c4afeaf", "api key")
	flag.Parse()

	app.DSN = "host=localhost port=5432 user=postgres password=postgres dbname=newsbyte sslmode=disable timezone=UTC connect_timeout=5"

	conn, err := app.connectToDB()
	if err != nil {
		log.Fatal(err)
	}
	app.DB = &dbrepo.PostgresDBRepo{DB: conn}
	app.ADB = &dbrepo.ArticlePostgresDBRepo{DB: conn}
	defer conn.Close()

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

	log.Println("db connection")
	err = http.ListenAndServe(fmt.Sprintf(":%d", port), app.routes())
	if err != nil {
		log.Fatal(err)
	}

}
