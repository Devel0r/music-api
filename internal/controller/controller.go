package controller

import (
	"encoding/json"
	"music/internal/model"
	"music/internal/services"
	"music/pkg/logger"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

type MainController struct {
	service *services.MainService
	log     *logger.Logger
	router  *mux.Router
}

func NewMainController(service *services.MainService, m *mux.Router, log *logger.Logger) *MainController {
	return &MainController{
		service: service,
		log:     log,
		router:  m,
	}
}

func (c *MainController) RegisterHandlers() {
	c.router.HandleFunc("/songs", c.handleSongs).Methods("GET", "POST")
	c.router.HandleFunc("/songs/{id}", c.handleSongByID).Methods("GET", "PUT", "DELETE")
	c.router.HandleFunc("/songs/{id}/text", c.GetSongText).Methods("GET")
}

func (c *MainController) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	c.router.ServeHTTP(w, r)
}

func (c *MainController) handleSongs(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		c.GetAllSongs(w, r)
	case http.MethodPost:
		c.AddSong(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func (c *MainController) handleSongByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		c.log.Error("Invalid song ID", logrus.Fields{"error": err})
		http.Error(w, "Invalid song ID", http.StatusBadRequest)
		return
	}

	switch r.Method {
	case http.MethodGet:
		c.GetSong(w, r, id)
	case http.MethodPut:
		c.UpdateSong(w, r, id)
	case http.MethodDelete:
		c.DeleteSong(w, r, id)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// GetAllSongs godoc
// @Summary Get all songs
// @Description Get songs with filters and pagination
// @Tags songs
// @Param group query string false "Filter by group"
// @Param song query string false "Filter by song title"
// @Param release_date query string false "Filter by release date"
// @Param limit query int false "Limit (default 10)"
// @Param offset query int false "Offset (default 0)"
// @Success 200 {array} model.Song
// @Router /songs [get]
func (c *MainController) GetAllSongs(w http.ResponseWriter, r *http.Request) {
	c.log.Info("Handling GET all songs request", logrus.Fields{})

	filters := map[string]string{
		"group":   r.URL.Query().Get("group"),
		"song":    r.URL.Query().Get("song"),
		"release": r.URL.Query().Get("release_date"),
		"lyrics":  r.URL.Query().Get("lyrics"),
		"link":    r.URL.Query().Get("link"),
	}

	query := r.URL.Query()

	limit, _ := strconv.Atoi(query.Get("limit"))
	if limit <= 0 {
		limit = 10
	}

	offset, _ := strconv.Atoi(query.Get("offset"))
	if offset < 0 {
		offset = 0
	}

	songs, err := c.service.GetAllSongs(filters, limit, offset)
	if err != nil {
		c.log.Error("Failed to get songs", logrus.Fields{"error": err})
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(songs); err != nil {
		c.log.Error("Failed to encode response", logrus.Fields{"error": err})
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}

// GetSong godoc
// @Summary Get song by ID
// @Description Get song details by its ID
// @Tags songs
// @Produce json
// @Param id path int true "Song ID"
// @Success 200 {object} model.Song
// @Failure 404 {object} map[string]string
// @Router /songs/{id} [get]
func (c *MainController) GetSong(w http.ResponseWriter, r *http.Request, id int) {
	c.log.Info("Handling GET song request", logrus.Fields{"song_id": id})

	song, err := c.service.GetSongByID(id)
	if err != nil {
		c.log.Error("Failed to get song", logrus.Fields{"error": err, "song_id": id})
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(song); err != nil {
		c.log.Error("Failed to encode response", logrus.Fields{"error": err})
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}

// GetSongText godoc
// @Summary Get paginated song lyrics
// @Description Get song text with pagination by verses
// @Tags songs
// @Produce plain
// @Param id path int true "Song ID"
// @Param page query int false "Page number (default 1)"
// @Param per_page query int false "Verses per page (default 3)"
// @Success 200 {string} string
// @Failure 400 {object} map[string]string
// @Router /songs/{id}/text [get]
func (c *MainController) GetSongText(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		c.log.Error("Invalid song ID", logrus.Fields{"error": err})
		http.Error(w, "Invalid song ID", http.StatusBadRequest)
		return
	}

	c.log.Info("Handling GET song text request", logrus.Fields{"song_id": id})

	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	if page < 1 {
		page = 1
	}

	perPage, _ := strconv.Atoi(r.URL.Query().Get("per_page"))
	if perPage < 1 {
		perPage = 3
	}

	text, err := c.service.GetSongText(id, page, perPage)
	if err != nil {
		c.log.Error("Failed to get song text", logrus.Fields{"error": err})
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/plain")
	if _, err := w.Write([]byte(text)); err != nil {
		c.log.Error("Failed to write response", logrus.Fields{"error": err})
	}
}

// AddSong godoc
// @Summary Add new song
// @Description Add new song to library with data from external API
// @Tags songs
// @Accept json
// @Produce json
// @Param song body model.Song true "Song data"
// @Success 201 {object} map[string]int
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /songs [post]
func (c *MainController) AddSong(w http.ResponseWriter, r *http.Request) {
	c.log.Info("Handling POST song request", logrus.Fields{})

	var song model.Song
	if err := json.NewDecoder(r.Body).Decode(&song); err != nil {
		c.log.Error("Failed to decode request body", logrus.Fields{"error": err})
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if song.GroupName == "" || song.SongTitle == "" {
		c.log.Error("Missing required fields", logrus.Fields{})
		http.Error(w, "Group and song title are required", http.StatusBadRequest)
		return
	}

	id, err := c.service.AddSong(song)
	if err != nil {
		c.log.Error("Failed to add song", logrus.Fields{"error": err})
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(map[string]int{"id": id}); err != nil {
		c.log.Error("Failed to encode response", logrus.Fields{"error": err})
	}
}

// UpdateSong godoc
// @Summary Update song
// @Description Update existing song data
// @Tags songs
// @Accept json
// @Produce json
// @Param id path int true "Song ID"
// @Param song body model.Song true "Updated song data"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /songs/{id} [put]
func (c *MainController) UpdateSong(w http.ResponseWriter, r *http.Request, id int) {
	c.log.Info("Handling PUT song request", logrus.Fields{"song_id": id})

	var song model.Song
	if err := json.NewDecoder(r.Body).Decode(&song); err != nil {
		c.log.Error("Failed to decode request body", logrus.Fields{"error": err})
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if err := c.service.UpdateSong(id, song); err != nil {
		c.log.Error("Failed to update song", logrus.Fields{"error": err})
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(map[string]string{"status": "success"}); err != nil {
		c.log.Error("Failed to encode response", logrus.Fields{"error": err})
	}
}


// DeleteSong godoc
// @Summary Delete song
// @Description Delete song from library
// @Tags songs
// @Produce json
// @Param id path int true "Song ID"
// @Success 200 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /songs/{id} [delete]
func (c *MainController) DeleteSong(w http.ResponseWriter, r *http.Request, id int) {
	c.log.Info("Handling DELETE song request", logrus.Fields{"song_id": id})

	if err := c.service.DeleteSong(id); err != nil {
		c.log.Error("Failed to delete song", logrus.Fields{"error": err})
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(map[string]string{"status": "success"}); err != nil {
		c.log.Error("Failed to encode response", logrus.Fields{"error": err})
	}
}
