package authapitests

import (
	"context"
	"database/sql"
	"encoding/json"
	"testing"

	"entgo.io/ent/dialect/sql/schema"
	"github.com/danielgtaylor/huma/v2/humatest"

	"api-go/ent/enttest"
	"api-go/routes/auth_api"
	"api-go/utils/auth"
	"api-go/utils/db"
	"api-go/utils/password"

	moderncsqlite "modernc.org/sqlite"
)

type signInResponse struct {
	Token string `json:"token"`
}

func init() {
	sql.Register("sqlite3", &moderncsqlite.Driver{})
}

func setEnt(t *testing.T) {
	t.Helper()

	client := enttest.Open(
		t,
		"sqlite3",
		"file:ent?mode=memory&cache=shared&_pragma=foreign_keys(1)",
		enttest.WithMigrateOptions(schema.WithForeignKeys(false)),
	)
	t.Cleanup(func() { client.Close() })
	db.EntDB = client

	auth.SetSecret("secret-secret")
}

func TestSignInAPI_Success(t *testing.T) {
	setEnt(t)

	ctx := context.Background()
	hashed, err := password.HashPassword("correct-horse-battery-staple")
	if err != nil {
		t.Fatalf("hash password: %v", err)
	}
	userObj, err := db.EntDB.User.Create().
		SetEmail("test@example.com").
		SetName("Test User").
		SetPassword(hashed).
		Save(ctx)
	if err != nil {
		t.Fatalf("create user: %v", err)
	}

	_, api := humatest.New(t)
	auth_api.Register(api)
	resp := api.Post("/auth/signin", map[string]any{
		"email":    "test@example.com",
		"password": "correct-horse-battery-staple",
	})
	if resp.Code != 200 {
		t.Fatalf("unexpected status: %d body=%s", resp.Code, resp.Body.String())
	}

	var out signInResponse
	if err := json.Unmarshal(resp.Body.Bytes(), &out); err != nil {
		t.Fatalf("unmarshal response: %v body=%s", err, resp.Body.String())
	}
	if out.Token == "" {
		t.Fatalf("expected token in response, got empty")
	}
	tokenUserID, ok := auth.GetJWT(out.Token)
	if !ok || tokenUserID != userObj.ID {
		t.Fatalf("expected valid token for user id=%d, got id=%d ok=%v", userObj.ID, tokenUserID, ok)
	}
}

func TestSignInAPI_WrongPassword(t *testing.T) {
	setEnt(t)

	ctx := context.Background()
	hashed, err := password.HashPassword("the-right-password")
	if err != nil {
		t.Fatalf("hash password: %v", err)
	}
	_, err = db.EntDB.User.Create().
		SetEmail("test@example.com").
		SetPassword(hashed).
		Save(ctx)
	if err != nil {
		t.Fatalf("create user: %v", err)
	}

	_, api := humatest.New(t)
	auth_api.Register(api)
	resp := api.Post("/auth/signin", map[string]any{
		"email":    "test@example.com",
		"password": "the-wrong-password",
	})

	if resp.Code != 400 {
		t.Fatalf("expected 400, got %d body=%s", resp.Code, resp.Body.String())
	}
}

func TestSignInAPI_UserNotFound(t *testing.T) {
	setEnt(t)

	_, api := humatest.New(t)
	auth_api.Register(api)
	resp := api.Post("/auth/signin", map[string]any{
		"email":    "does-not-exist@example.com",
		"password": "irrelevant",
	})

	if resp.Code != 404 {
		t.Fatalf("expected 404, got %d body=%s", resp.Code, resp.Body.String())
	}
}
