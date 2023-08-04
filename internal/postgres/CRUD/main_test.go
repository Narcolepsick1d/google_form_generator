package CRUD

import (
	"fmt"
	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/jmoiron/sqlx"
	"os"
	"testing"
)

var (
	db *sqlx.DB
)

func TestMain(m *testing.M) {
	var err error
	conStr := fmt.Sprintf("host=%s port=%v user=%s password=%s dbname=%s sslmode=%s", "localhost",
		5432, "postgres", "21garik21", "postgres", "disable")

	//conStr := fmt.Sprintf("host=%s port=%v user=%s password=%s dbname=%s sslmode=%s", "10.100.1.71",
	//	30032, "postgres", "cffcc8@A27b7", "vendoo_go_loan_service", "disable")

	db, err = sqlx.Connect("pgx", conStr)
	if err != nil {
		panic(err)
	}

	os.Exit(m.Run())

}
