package restapi

import (
	"encoding/json"
	"io"
	"net/http"
	"s3-test/internal/application"
	"s3-test/internal/models"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

const (
	HeaderUser   = "user_id"
	HeaderRentId = "rent_id"
)

type Server struct {
	core *application.Core
}

func NewServer(core *application.Core) *Server {
	return &Server{core: core}
}

func (s *Server) Run(port string) error {
	router := mux.NewRouter()

	router.HandleFunc("/api/v1/rent", s.rent)
	router.HandleFunc("/api/v1/rent_start", s.rentStart)
	router.HandleFunc("/api/v1/rent_stop", s.rentStop)
	router.HandleFunc("/api/v1/rent_info", s.rentInfo)
	router.HandleFunc("/api/v1/upload_photo", s.uploadPhoto)

	http.Handle("/", router)
	return http.ListenAndServe(port, nil)
}

// rent Инициализирует сессию для пользователя
func (s *Server) rent(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	uid, err := strconv.Atoi(r.Header.Get(HeaderUser))
	if err != nil {
		logrus.Errorf("Can`t create session: %v\n", err)
		http.Error(w, "Wrong request!", http.StatusBadRequest)
		return
	}

	rentID, err := s.core.Rent(models.User{
		ID: uid,
	})
	if err != nil {
		logrus.Errorf("DB fail: %v\n", err)
		http.Error(w, "Can`t create session!", http.StatusInternalServerError)
	}

	rentInfo, err := s.core.RentInfo(rentID)
	if err != nil {
		logrus.Errorf("Cant get rent info: %v\n", err)
		http.Error(w, "Can`t get rent info for started session!", http.StatusAccepted)
		return
	}

	s.sendModel(w, rentInfo)
}

// rentStart Запуск сессии аренды
func (s *Server) rentStart(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	rentID, err := strconv.Atoi(r.Header.Get(HeaderRentId))
	if err != nil {
		logrus.Errorf("Can`t get rent id: %v : raw-header:%s\n", err, r.Header.Get(HeaderRentId))
		http.Error(w, "Wrong request!", http.StatusBadRequest)
		return
	}

	if s.core.RentStart(rentID) != nil {
		logrus.Errorf("Unmarshal fail: %v\n", err)
		http.Error(w, "Can`t start rent!", http.StatusInternalServerError)
		return
	}

	rentInfo, err := s.core.RentInfo(rentID)
	if err != nil {
		logrus.Errorf("Cant get rent info: %v\n", err)
		http.Error(w, "Can`t get rent info for started session!", http.StatusAccepted)
		return
	}

	s.sendModel(w, rentInfo)
}

// rentStop Завершение сессии аренды
func (s *Server) rentStop(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	rentID, err := strconv.Atoi(r.Header.Get(HeaderRentId))
	if err != nil {
		logrus.Errorf("Can`t get rent id: %v : raw-header:%s\n", err, r.Header.Get(HeaderRentId))
		http.Error(w, "Wrong request!", http.StatusBadRequest)
		return
	}

	if s.core.RentStop(rentID) != nil {
		http.Error(w, "Can`t stop rent!", http.StatusInternalServerError)
		return
	}

	rentInfo, err := s.core.RentInfo(rentID)
	if err != nil {
		logrus.Errorf("Cant get rent info: %v\n", err)
		http.Error(w, "Can`t get rent info for started session!", http.StatusAccepted)
		return
	}

	s.sendModel(w, rentInfo)
}

// uploadPhoto Загрузить фото из запроса в наш сторадж
func (s *Server) uploadPhoto(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	rentID, err := strconv.Atoi(r.Header.Get(HeaderRentId))
	if err != nil {
		logrus.Errorf("Can`t get rent id: %v\n", err)
		http.Error(w, "Wrong request!", http.StatusBadRequest)
		return
	}

	src, hdr, err := r.FormFile("photo")
	if err != nil {
		http.Error(w, "Wrong request!", http.StatusBadRequest)
		return
	}
	defer src.Close()

	rentInfo, err := s.core.RentInfo(rentID)
	if err != nil {
		logrus.Errorf("Cant get rent info: %v\n", err)
		http.Error(w, "Can`t get rent info for started session!", http.StatusAccepted)
		return
	}

	img := models.ImageUnit{
		Payload:     src,
		PayloadSize: hdr.Size,
		User:        rentInfo.User,
	}

	err = s.core.UploadPhoto(r.Context(), img, rentID)
	if err != nil {
		logrus.Errorf("Fail update img in database: %v\n", err)
		http.Error(w, "Can`t upload photo!", http.StatusInternalServerError)
	}
}

func (s *Server) rentInfo(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	rentID, err := strconv.Atoi(r.Header.Get(HeaderRentId))
	if err != nil {
		logrus.Errorf("Can`t get rent id: %v : raw-header:%s\n", err, r.Header.Get(HeaderRentId))
		http.Error(w, "Wrong request!", http.StatusBadRequest)
		return
	}

	rentInfo, err := s.core.RentInfo(rentID)
	if err != nil {
		logrus.Errorf("Cant get rent info: %v\n", err)
		http.Error(w, "Can`t get rent info for started session!", http.StatusAccepted)
		return
	}

	s.sendModel(w, rentInfo)
}

func (s Server) sendModel(w http.ResponseWriter, model interface{}) {
	b, err := json.Marshal(model)
	if err != nil {
		logrus.Errorf("Marshal fail: %v\n", err)
	}

	_, err = w.Write([]byte(b))
	if err != nil {
		logrus.Errorf("Write respinse fail: %v\n", err)
		return
	}
}

func (s Server) sendImage(w http.ResponseWriter, img io.Reader) {
	w.Header().Set("Content-Type", "image/jpeg")
	if _, err := io.Copy(w, img); err != nil {
		logrus.Errorf("Cant send photo: %v\n", err)
		http.Error(w, "Can`t download photo!", http.StatusInternalServerError)
	}
}
