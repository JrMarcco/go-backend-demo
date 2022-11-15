package db

import (
	"database/sql"
	"github.com/stretchr/testify/suite"
	"github/jrmarcco/go-backend-demo/util"
	"log"
	"testing"

	_ "github.com/go-sql-driver/mysql"
)

type mysqlTestSuite struct {
	suite.Suite

	dbDriver string
	dbSource string

	conn *sql.DB

	queries *Queries
}

func TestMySQL(t *testing.T) {
	config, err := util.LoadConfig("../..")
	if err != nil {
		log.Fatal("can not load config:", err)
	}

	suite.Run(t, &mysqlTestSuite{
		dbDriver: config.Db.Driver,
		dbSource: config.Db.Source,
	})
}

func (m *mysqlTestSuite) SetupSuite() {
	conn, err := sql.Open(m.dbDriver, m.dbSource)
	if err != nil {
		log.Fatalln(err)
	}

	m.conn = conn
	m.queries = New(conn)
}
