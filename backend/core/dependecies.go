package core

import (
	"digitalpaper/backend/core/logger"
	"github.com/alexedwards/scs/v2"
)

type Application struct {
	Log *logger.Logger
	SessionManager *scs.SessionManager
}
