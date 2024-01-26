package application

import (
	"fmt"
	"math/rand"
	"s3-test/internal/models"
	"strconv"
	"time"
)

// Генерируем название файла исходя из userID
func GenerateObjectName(user models.User) string {
	t := time.Now()
	formatted := fmt.Sprintf("%d-%02d-%02dT%02d:%02d:%02d",
		t.Year(), t.Month(), t.Day(),
		t.Hour(), t.Minute(), t.Second())
	return fmt.Sprintf("%s/%s.%s",
		strconv.Itoa(user.ID),
		formatted,
		"png")
}
func generateRentID() int {
	return rand.Int()
}
func enrichImagesLinks(session models.Rent, minioEndpoint string) models.Rent {
	for i, image := range session.ImagesBefore {
		session.ImagesBefore[i] = fmt.Sprintf("%s/%s", minioEndpoint, image)
	}

	for i, image := range session.ImagesAfter {
		session.ImagesAfter[i] = fmt.Sprintf("%s/%s", minioEndpoint, image)
	}

	return session
}
