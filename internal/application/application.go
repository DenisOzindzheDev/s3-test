// core use cases
package application

import (
	"context"
	"s3-test/internal/db"
	"s3-test/internal/models"
	"s3-test/internal/s3"
)

type Core struct {
	storage            s3.ImageStorage
	database           db.DB
	ImageStorageSVCDSN string
}

func NewSkater(storage s3.ImageStorage, database db.DB, minioEndpointURL string) (*Core, error) {
	return &Core{storage: storage, database: database, ImageStorageSVCDSN: minioEndpointURL}, nil
}

func (c *Core) Rent(usr models.User) (int, error) {
	rentOpts := models.Rent{
		User:   usr,
		RentID: generateRentID(),
	}
	return rentOpts.RentID, c.database.CreateNewRentSession(rentOpts)

}
func (c *Core) RentStart(rentID int) error {
	return c.database.StartRentSession(rentID)

}
func (c *Core) RentStop(rentID int) error {
	return c.database.CompleteRentSession(rentID)
}
func (c *Core) UploadPhoto(ctx context.Context, img models.ImageUnit, rentID int) error {
	imgPath, err := c.storage.UploadFile(ctx, img)
	if err != nil {
		return err
	}
	return c.database.AddImageRecord(imgPath, rentID)

}
func (c *Core) RentInfo(rentID int) (models.Rent, error) {
	session, err := c.database.GetRentStatus(rentID)
	if err != nil {
		return models.Rent{}, err
	}

	session = enrichImagesLinks(session, c.ImageStorageSVCDSN)
	return session, nil

}
