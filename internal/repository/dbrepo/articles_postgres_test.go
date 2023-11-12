//go:build article_integration

package dbrepo

import (
	"database/sql"
	"fmt"
	"github.com/hirotofuu/newsbyte/internal/models"
	"github.com/hirotofuu/newsbyte/internal/repository"
	"log"
	"os"
	"testing"
	"time"

	_ "github.com/jackc/pgconn"
	_ "github.com/jackc/pgx/v4"
	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/ory/dockertest/v3"
	"github.com/ory/dockertest/v3/docker"
)

var (
	host     = "localhost"
	user     = "postgres"
	password = "postgres"
	dbName   = "user_test"
	port     = "5435"
	dsn      = "host=%s port=%s user=%s password=%s dbname=%s sslmode=disable timezone=UTC connect_timeout=5"
)

var resource *dockertest.Resource
var pool *dockertest.Pool
var testDB *sql.DB
var testRepo repository.Databaserepo
var testArticleRepo repository.ArticleDatabaserepo

func TestMain(m *testing.M) {
	p, err := dockertest.NewPool("")
	if err != nil {
		log.Fatalf("could not connect to docker; is it running? %s", err)
	}

	pool = p

	opts := dockertest.RunOptions{
		Repository: "postgres",
		Tag:        "14.5",
		Env: []string{
			"POSTGRES_USER=" + user,
			"POSTGRES_PASSWORD=" + password,
			"POSTGRES_DB=" + dbName,
		},
		ExposedPorts: []string{"5432"},
		PortBindings: map[docker.Port][]docker.PortBinding{
			"5432": {
				{HostIP: "0.0.0.0", HostPort: port},
			},
		},
	}

	resource, err = pool.RunWithOptions(&opts)
	if err != nil {
		_ = pool.Purge(resource)
		log.Fatalf("could not start resource: %s", err)
	}

	if err := pool.Retry(func() error {
		var err error
		testDB, err = sql.Open("pgx", fmt.Sprintf(dsn, host, port, user, password, dbName))
		if err != nil {
			log.Println("Error:", err)
			return err
		}
		return testDB.Ping()
	}); err != nil {
		_ = pool.Purge(resource)
		log.Fatalf("could not connect to database")
	}

	err = createTables()
	if err != nil {
		log.Fatalf("error creating tables: %s", err)
	}

	testRepo = &PostgresDBRepo{DB: testDB}
	testArticleRepo = &ArticlePostgresDBRepo{DB: testDB}

	code := m.Run()

	if err := pool.Purge(resource); err != nil {
		log.Fatalf("could not purge resource: %s", err)
	}
	os.Exit(code)
}

func createTables() error {
	tableSQL, err := os.ReadFile("./testdata/create_tables.sql")
	if err != nil {
		fmt.Println(err)
		return err
	}

	_, err = testDB.Exec(string(tableSQL))
	if err != nil {
		fmt.Println(err)
		return err
	}

	return nil
}

func TestPostgresDBRepoRegisterUsers(t *testing.T) {
	users, err := testRepo.AllUsers()

	if err != nil {
		t.Errorf("all users reports an error: %s", err)
	}

	if len(users) != 0 {
		t.Errorf("all users reports wrong size; expected 1, but got %d", len(users))
	}

	testUser := models.User{
		UserName:  "Jack",
		Email:     "jack@smith.com",
		Password:  "secret",
		AvatarImg: "http:/clap",
		Profile:   "",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	_, err = testRepo.InsertUser(testUser)
	if err != nil {
		t.Errorf("insert user reports an error: %s", err)
	}

	users, err = testRepo.AllUsers()
	if err != nil {
		t.Errorf("all users reports an error: %s", err)
	}

	if len(users) != 1 {
		t.Errorf("all users reports wrong size; expected 2, but got %d", len(users))
	}
}

func TestArticlePostgresDBRepoInsert(t *testing.T) {
	tag := "hiroto"
	testArticle := models.Article{
		Title:      "you know say",
		Content:    "tomorrow tomorrow i love yeah tomorrow",
		TagsIn:     []string{tag},
		Medium:     1,
		UserID:     1,
		IsOpenFlag: true,
		CommentOK:  true,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}
	_, err := testArticleRepo.InsertArticle(testArticle)
	if err != nil {
		t.Errorf("insert article reports an error: %s", err)
	}

}

func TestArticlePostgresDBRepoGetArticles(t *testing.T) {
	tag := "hiroto"

	testArticle := models.Article{
		Title:      "you know say",
		Content:    "null",
		TagsIn:     []string{tag, "呪術廻戦"},
		Medium:     0,
		UserID:     1,
		IsOpenFlag: true,
		CommentOK:  true,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}
	_, err := testArticleRepo.InsertArticle(testArticle)
	if err != nil {
		t.Errorf("insert article reports an error: %s", err)
	}

	articles, err := testArticleRepo.UserArticles(1)
	if err != nil {
		t.Errorf("all users reports an error: %s", err)
	}

	if len(articles) != 2 {
		t.Errorf("all users reports wrong size; expected 2, but got %d", len(articles))
	}

	articles, err = testArticleRepo.AllArticles()
	if err != nil {
		t.Errorf("all users reports an error: %s", err)
	}

	if len(articles) != 2 {
		t.Errorf("all users reports wrong size; expected 2, but got %d", len(articles))
	}
	Ids, err := testArticleRepo.StateGoodArticle(1)
	if err != nil {
		t.Errorf("state good reports an error: %s", err)
	}
	if len(Ids) != 0 {
		t.Errorf("state good reports wrong size; expected 1, but got %d", len(Ids))
	}
	fmt.Println(Ids)

	article, err := testArticleRepo.OneArticle(2)
	if err != nil {
		t.Errorf("all users reports an error: %s", err)
	}
	if article.Title != "you know say" {
		if err != nil {
			t.Errorf("article title should be you know say: but got %s", article.Title)
		}
	}
}
