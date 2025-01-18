package config

import (
	"github.com/sirupsen/logrus"
	"styl-monolith/internal/media/adapters/rest"
)

func BuildMediaCrudDomainHandler(log *logrus.Logger) rest.CrudPostHandler {
	mediaHandler := rest.NewCrudPostHandler(log)
	return mediaHandler
}
