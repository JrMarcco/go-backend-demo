package enttest

import (
	"github.com/jrmarcco/go-backend-demo/db/ent"
	"github.com/jrmarcco/go-backend-demo/util"
	"github.com/stretchr/testify/suite"
	"log"
	"testing"

	_ "github.com/go-sql-driver/mysql"
)

type entTestSuite struct {
	suite.Suite

	driver string
	source string

	client *ent.Client
}

func TestEnt(t *testing.T) {
	cfg, err := util.LoadConfig("../../..")
	if err != nil {
		log.Fatal(err)
	}

	suite.Run(t, &entTestSuite{
		driver: cfg.Db.Driver,
		source: cfg.Db.Source,
	})
}

func (e *entTestSuite) SetupSuite() {
	client, err := ent.Open(e.driver, e.source)

	if err != nil {
		log.Fatal(err)
	}
	e.client = client
}
