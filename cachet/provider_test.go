package cachet

import (
	"fmt"
	"os"
	"testing"

	"database/sql"
	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	_ "github.com/lib/pq"
)

var testProvider = Provider()

func testPreCheck(t *testing.T) {
	if v := os.Getenv("CACHET_TOKEN"); v == "" {
		os.Setenv("CACHET_TOKEN", setupApiKey(t))
	}

	if v := os.Getenv("CACHET_URL"); v == "" {
		t.Fatal("CACHET_URL must be set for tests")
	}

}

func testProviderFactory() map[string]func() (*schema.Provider, error) {
	var out = make(map[string]func() (*schema.Provider, error))

	out["cachet"] = func() (*schema.Provider, error) {
		return testProvider, nil
	}

	return out
}

func setupApiKey(t *testing.T) string {
	psqlconn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", "localhost", 8001, "postgres", "postgres", "postgres")

	// open database
	db, err := sql.Open("postgres", psqlconn)
	if err != nil {
		t.Fatal(err)
	}

	// close database
	defer db.Close()

	// check db
	err = db.Ping()
	if err != nil {
		t.Fatal(err)
	}

	rows, err := db.Query(`SELECT "username", "api_key" FROM "chq_users" WHERE "username" = 'admin'`)
	if err != nil {
		t.Fatal(err)
	}

	defer rows.Close()
	for rows.Next() {
		var username, apiKey string
		err = rows.Scan(&username, &apiKey)
		if err != nil {
			t.Fatal(err)
		}

		if apiKey != "" {
			return apiKey
		}
	}

	insert := `insert into "chq_users" ("username", "password", "email", "api_key") values($1, $2, $3, $4)`
	apiKey := uuid.New().String()
	_, err = db.Exec(insert, "admin", "admin", "admin@admin.com", apiKey)
	if err != nil {
		t.Fatal(err)
	}

	return apiKey
}
