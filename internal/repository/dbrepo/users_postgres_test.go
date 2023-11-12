//go:build integration

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
	newID, err := testRepo.InsertUser(testUser)
	if err != nil {
		t.Errorf("insert user reports an error: %s", err)
	}

	testUser = models.User{
		UserName:  "hiroto",
		Email:     "hiroto@smith.com",
		Password:  "secret",
		AvatarImg: "http:/clap",
		Profile:   "",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	secondID, err := testRepo.InsertUser(testUser)
	if err != nil {
		t.Errorf("insert user reports an error: %s", err)
	}

	users, err = testRepo.AllUsers()
	if err != nil {
		t.Errorf("all users reports an error: %s", err)
	}

	if newID != 1 {
		t.Errorf("all users reports wrong size; expected 2, but got %d", len(users))
	}

	user, err := testRepo.GetUserByID(newID)
	if err != nil {
		t.Errorf("getuserbyid reports an error: %s", err)
	}

	if user.UserName != "Jack" {
		t.Errorf("expect jack but got %v", user.UserName)
	}

	user, err = testRepo.GetUserIdName("hiroto")
	if err != nil {
		t.Errorf("getuserbyid reports an error: %s", err)
	}

	if user.UserName != "hiroto" {
		t.Errorf("expect jack but got %v", user.UserName)
	}

	err = testRepo.InsertFollow(newID, secondID)
	if err != nil {
		t.Errorf("insert good comment reports an error: %s", err)
	}

	num, err := testRepo.GetFollowingUserIDs(newID)
	if err != nil {
		t.Errorf("get followings reports an error: %s", err)
	}
	if len(num) != 0 {
		t.Errorf("fet following reports wrong size; expected 2, but got %d", len(num))
	}

	user, err = testRepo.OneUser(newID)
	if err != nil {
		t.Errorf("get followings reports an error: %s", err)
	}
	if user.FollowingsCount != 0 {
		t.Errorf("one user reports wrong following num; expected 0, but got %d", user.FollowingsCount)
	}
	if user.FollowedsCount != 1 {
		t.Errorf("one user reports wrong followed num; expected 0, but got %d", user.FollowingsCount)
	}

	fUsers, err := testRepo.FollowingUsers(secondID)
	if err != nil {
		t.Errorf("following reports an error: %s", err)
	}
	if len(fUsers) != 1 {
		t.Errorf("expect 1 but got %v", len(fUsers))
	}

	fUsers, err = testRepo.FollowedUsers(newID)
	if err != nil {
		t.Errorf("followed reports an error: %s", err)
	}
	if len(fUsers) != 1 {
		t.Errorf("expect 1 but got %v", len(fUsers))
	}

	err = testRepo.DeleteFollow(newID, secondID)
	if err != nil {
		t.Errorf("delete good comment reports an error: %s", err)
	}
}

func TestPostgresDBRepoGetUsersByEmail(t *testing.T) {
	testEmail := "jack@smith.com"

	users, err := testRepo.GetUserByEmail(testEmail)
	if err != nil {
		t.Errorf("getuserbyemail reports an error: %s", err)
	}
	if users.UserName != "Jack" {
		t.Errorf("expect jack but got %v", users.UserName)
	}

}
