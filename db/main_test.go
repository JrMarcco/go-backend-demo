package db

import (
	"database/sql"
	"github.com/stretchr/testify/suite"
	"log"
	"testing"

	_ "github.com/go-sql-driver/mysql"
)

type mysqlTestSuite struct {
	suite.Suite

	dbDriver string
	dbSource string

	db *sql.DB

	queries *Queries
}

func TestMySQL(t *testing.T) {
	suite.Run(t, &mysqlTestSuite{
		dbDriver: "mysql",
		dbSource: "root:@tcp(127.0.0.1:3306)/simple_bank?parseTime=true",
	})
}

func (m *mysqlTestSuite) SetupSuite() {
	db, err := sql.Open(m.dbDriver, m.dbSource)
	if err != nil {
		log.Fatalln(err)
	}

	m.db = db
	m.queries = New(db)
}
