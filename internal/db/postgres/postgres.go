package postgres

import (
	"database/sql"
	"log"
	"s3-test/internal/db"
	"s3-test/internal/models"
	"strings"

	_ "github.com/lib/pq"

	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
)

type PGS struct {
	db *sqlx.DB
}

func NewDataBase(dbPath string) (db.DB, func() error, error) {
	db, err := sqlx.Connect(
		"postgres",
		dbPath)
	if err != nil {
		log.Fatalln(err)
	}

	return &PGS{db: db}, db.Close, err
}

func (s *PGS) CreateNewRentSession(rent models.Rent) error {
	_, err := s.db.Exec(createRentSession, rent.User.ID, rent.RentID)
	if err != nil {
		logrus.Errorf("Can`t create rent record in db: %v\n", err)
	}

	return err
}
func (s *PGS) StartRentSession(rentID int) error {
	_, err := s.db.Exec(startRentSession, rentID)

	return err
}
func (s *PGS) CompleteRentSession(rentID int) error {
	_, err := s.db.Exec(completeRentSession, rentID)
	if err != nil {
		logrus.Errorf("Can`t update rent record in db: %v\n", err)
	}

	return err
}

func (s *PGS) GetRentSessions(user models.User) ([]models.Rent, error) {
	rentSessionsList := make([]models.Rent, 0)
	rows, err := s.db.Query(getRentSessions, user.ID)
	if err != nil {
		logrus.Errorf("Can`t get rent records in db: %v\n", err)
		return rentSessionsList, err
	}
	defer rows.Close()

	for rows.Next() {
		r := models.Rent{}

		var imagesBefore, imagesAfter sql.NullString
		var startedAt, completedAt sql.NullTime

		err := rows.Scan(&r.RentID, &r.User.ID, &startedAt, &completedAt, &imagesBefore, &imagesAfter)
		if err != nil {
			logrus.Errorf("Can`t get rent record from rows: %v\n", err)
			continue
		}

		r.StartedAt = startedAt.Time
		r.CompletedAt = completedAt.Time

		r.RentInfo.ImagesBefore = getArrayFromStr(imagesBefore.String)
		r.RentInfo.ImagesAfter = getArrayFromStr(imagesAfter.String)

		rentSessionsList = append(rentSessionsList, r)
	}

	return rentSessionsList, nil
}

func (s *PGS) AddImageRecord(img string, rentID int) error {
	// Получаем информацию о сессии аренды
	rent, err := s.GetRentStatus(rentID)
	if err != nil {
		logrus.Errorf("Can`t get rent record in db: %v\n", err)
		return err
	}

	// В зависимости от того были ли загружены фото ДО начала аренды
	// Или после добавляем запись в соответствующее поле в БД
	if rent.StartedAt.IsZero() {
		return s.updateImages(rent.ImagesBefore, img, updateImagesBefore, rentID)
	}

	return s.updateImages(rent.ImagesAfter, img, updateImagesAfter, rentID)
}

func (s *PGS) GetRentStatus(rentID int) (models.Rent, error) {
	rent := models.Rent{}
	var imagesBefore, imagesAfter sql.NullString
	var startedAt, completedAt sql.NullTime

	rows, err := s.db.Query(getStartedAtStatus, rentID)
	if err != nil {
		logrus.Errorf("Can`t get rent record in db: %v\n", err)
		return rent, err
	}
	defer rows.Close()

	for rows.Next() {
		err := rows.Scan(&rent.RentID, &rent.User.ID, &startedAt, &completedAt, &imagesBefore, &imagesAfter)
		if err != nil {
			logrus.Errorf("Can`t get rent record from rows: %v\n", err)
			continue
		}
	}

	rent.StartedAt = startedAt.Time
	rent.CompletedAt = completedAt.Time

	rent.RentInfo.ImagesBefore = getArrayFromStr(imagesBefore.String)
	rent.RentInfo.ImagesAfter = getArrayFromStr(imagesAfter.String)

	return rent, err
}

func (s *PGS) updateImages(old []string, new, req string, rentID int) error {
	// Добавляем в список старых записей в поле
	// новую запись об изображении
	old = append(old, new)
	new = strings.Join(old, ",")

	_, err := s.db.Exec(req, new, rentID)
	if err != nil {
		logrus.Errorf("Can`t update image record in db: %v\n", err)
	}

	return err
}

func getArrayFromStr(in string) []string {
	if in != "" {
		return strings.Split(in, ",")
	}
	return []string{}
}
