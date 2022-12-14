package db

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jrmarcco/go-backend-demo/util"
	"github.com/stretchr/testify/suite"
	"log"
	"testing"
)

type sqlcTestSuite struct {
	suite.Suite

	dbDriver string
	dbSource string

	conn    *sql.DB
	queries *Queries
}

func TestMySQL(t *testing.T) {
	config, err := util.LoadConfig("../..")
	if err != nil {
		log.Fatal("can not load config:", err)
	}

	suite.Run(t, &sqlcTestSuite{
		dbDriver: config.Db.Driver,
		dbSource: config.Db.Source,
	})
}

func (m *sqlcTestSuite) SetupSuite() {
	conn, err := sql.Open(m.dbDriver, m.dbSource)
	if err != nil {
		log.Fatalln(err)
	}

	m.conn = conn
	m.queries = New(conn)
}
