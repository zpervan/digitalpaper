package tests

import (
	"context"
	"digitalpaper/backend/core"
	"digitalpaper/backend/core/logger"
	"digitalpaper/backend/database"
	"fmt"
	mdb "github.com/ONSdigital/dp-mongodb-in-memory"
	"github.com/alexedwards/scs/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"go.mongodb.org/mongo-driver/bson"
)

type DatabaseTestSuite struct {
	suite.Suite
	dummyContext context.Context
	dbServerMock *mdb.Server
	dbHandlers   *database.Database
}

// Configuration variables
var selectAll = bson.D{} // No filter criteria

// Setup
func (dts *DatabaseTestSuite) SetupSuite() {
	fmt.Println("Setting up database test suite...")

	dts.dummyContext = context.Background()

	// Setup in-memory MongoDB server
	dbServer, err := mdb.StartWithOptions(dts.dummyContext, "5.0.2", mdb.WithPort(27017))
	if err != nil {
		assert.Fail(dts.T(), "error while starting dummy MongoDB server")
	}

	// Setup database handlers/functions
	dummyApp := core.Application{
		Log:            logger.New(),
		SessionManager: &scs.SessionManager{},
	}

	dts.dbServerMock = dbServer
	dts.dbHandlers, err = database.NewDatabase(&dummyApp, dts.dbServerMock.URI())

	// @TODO: Create database connection test?
	assert.Nil(dts.T(), err, "error while creating new database instance")

	fmt.Println("Setting up database test suite... COMPLETE")
}

func (dts *DatabaseTestSuite) TearDownSuite() {
	fmt.Println("Tearing down database test suite...")

	dts.dbServerMock.Stop(dts.dummyContext)

	fmt.Println("Tearing down database test suite... COMPLETE")
}
