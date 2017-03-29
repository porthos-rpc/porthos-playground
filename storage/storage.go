package storage

import (
	"github.com/porthos-rpc/porthos-playground/models"
)

// Storage structure.
type Storage interface {
	SaveServiceSpecs(serviceSpecs *models.ServiceSpecs)
	GetSpecs() ([]*models.ServiceSpecs, error)
}
