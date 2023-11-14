//go:build comment_integration

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
var testCommentRepo repository.CommentDatabaserepo

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
	testCommentRepo = &CommentPostgresDBRepo{DB: testDB}

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

func TestPostgresDBRepoInsertComment(t *testing.T) {
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
		Profile:   "",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	userID, err := testRepo.InsertUser(testUser)
	if err != nil {
		t.Errorf("insert user reports an error: %s", err)
	}

	testArticle := models.Article{
		Title:     "you know say",
		Content:   "tomorrow tomorrow i love yeah tomorrow",
		TagsIn:    []string{"釘崎野薔薇", "呪術廻戦"},
		Medium:    1,
		UserID:    1,
		CommentOK: true,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	articleID, err := testArticleRepo.InsertArticle(testArticle)
	if err != nil {
		t.Errorf("insert article reports an error: %s", err)
	}

	testComment := models.Comment{
		Comment:   "this article is saiko----",
		UserID:    userID,
		ArticleID: articleID,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	commentID, err := testCommentRepo.InsertComment(testComment)
	if err != nil {
		t.Errorf("insert comment reports an error: %s", err)
	}

	if commentID != 1 {
		t.Errorf("insert comment reports wrong id; expected 1, but got %d", commentID)
	}

	comments, err := testCommentRepo.ArticleComments(articleID)
	if err != nil {
		t.Errorf("article comments reports an error: %s", err)
	}

	if len(comments) != 1 {
		t.Errorf("article comment reports wrong size; expected 1, but got %d", len(comments))
	}

	comments, err = testCommentRepo.UserComments(userID, userID)
	if err != nil {
		t.Errorf("user comments reports an error: %s", err)
	}

	if len(comments) != 1 {
		t.Errorf("article comment reports wrong size; expected 1, but got %d", len(comments))
	}

	err = testCommentRepo.InsertGoodComment(commentID, userID)
	if err != nil {
		t.Errorf("insert good comment reports an error: %s", err)
	}

	come, err := testCommentRepo.OneComment(commentID, userID)
	if err != nil {
		t.Errorf("one comments reports an error: %s", err)
	}
	if come.IsGoodFlag != 1 {
		t.Errorf("one comment reports wrong value; expected 0, but got %d", come.IsGoodFlag)
	}

	err = testCommentRepo.DeleteGoodComment(commentID, userID)
	if err != nil {
		t.Errorf("delete good comment reports an error: %s", err)
	}

	come, err = testCommentRepo.OneComment(commentID, userID)
	if come.IsGoodFlag != 0 {
		t.Errorf("one comment reports wrong value; expected 0, but got %d", come.IsGoodFlag)
	}

	err = testCommentRepo.DeleteComment(commentID)
	if err != nil {
		t.Errorf("delete comment reports an error: %s", err)
	}

	comments, err = testCommentRepo.ArticleComments(articleID)
	if err != nil {
		t.Errorf("article comments reports an error: %s", err)
	}

	if len(comments) != 0 {
		t.Errorf("article comment reports wrong size; expected 0, but got %d", len(comments))
	}

}
