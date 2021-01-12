package deployment

import (
	"github.com/pkg/errors"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	constants "nenvoy.com/pkg/constants"
	"nenvoy.com/pkg/host"
	"nenvoy.com/pkg/network"
)

//Deployment - Struct for the deployment data in the database
type Deployment struct {
	gorm.Model
	ID       uint
	Name     string
	Hosts    []host.Host
	Networks []network.Network
}

// GetDeploymentByID - Gets a deployment from the database by it's ID
func GetDeploymentByID(depID uint) (Deployment, error) {
	// Connect and open the database
	db, err := gorm.Open(sqlite.Open(constants.DBPath), &gorm.Config{})
	if err != nil {
		return Deployment{}, errors.Wrap(err, "failed to connect database")
	}

	// Get the deployment from the db
	var dep Deployment
	err = db.First(&dep, depID).Error
	if err != nil {
		return dep, errors.Wrap(err, "failed to find deployment")
	}

	return dep, nil
}

// GetDeploymentByName - Gets a deployment from the database by it's Name
func GetDeploymentByName(depName string) (Deployment, error) {
	// Connect and open the database
	db, err := gorm.Open(sqlite.Open(constants.DBPath), &gorm.Config{})
	if err != nil {
		return Deployment{}, errors.Wrap(err, "failed to connect database")
	}

	// Get the deployment from the db
	var dep Deployment
	err = db.Where("name = ?", depName).First(&dep).Error
	if err != nil {
		return dep, errors.Wrap(err, "failed to find deployment")
	}

	return dep, nil
}
