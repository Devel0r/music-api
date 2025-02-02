package services

import (
	"encoding/json"
	"fmt"
	"io"
	"music/internal/model"
	"music/internal/repository"
	"music/pkg/config"
	"music/pkg/logger"
	"net/http"
	"net/url"
	"strings"

	"github.com/sirupsen/logrus"
)

type MainService struct {
	repo   *repository.MainRepository
	client *http.Client
	cfg    *config.Config
	log    *logger.Logger
}

func NewMainService(repo *repository.MainRepository, client *http.Client, cfg *config.Config, log *logger.Logger) *MainService {
	return &MainService{
		repo:   repo,
		client: client,
		cfg:    cfg,
		log:    log,
	}
}

func (s *MainService) GetAllSongs(filters map[string]string, limit, offset int) ([]model.Song, error) {
    s.log.Info("Getting filtered songs", logrus.Fields{
        "filters": filters,
        "limit":   limit,
        "offset":  offset,
    })
    
    return s.repo.GetAllSongs(filters, limit, offset)
}

func (s *MainService) AddSong(song model.Song) (int, error) {
	s.log.Info("Adding song", logrus.Fields{"group": song.GroupName, "song": song.SongTitle})
	resp, err := s.client.Get(fmt.Sprintf("%s/info?group=%s&song=%s", s.cfg.ExternalAPI, url.QueryEscape(song.GroupName), url.QueryEscape(song.SongTitle)))
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("external API error: status %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0, err
	}

	var songDetail model.SongDetail
	err = json.Unmarshal(body, &songDetail)
	if err != nil {
		return 0, err
	}

	song.Lyrics = songDetail.Text
	song.ReleaseDate = songDetail.ReleaseDate
	song.YouTubeLink = songDetail.Link

	id, err := s.repo.AddSong(song)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (s *MainService) UpdateSong(id int, song model.Song) error {
	s.log.Info("Updating song", logrus.Fields{"id": id, "group": song.GroupName, "song": song.SongTitle})
	return s.repo.UpdateSong(id, song)
}

func (s *MainService) DeleteSong(id int) error {
	s.log.Info("Deleting song", logrus.Fields{"id": id})
	return s.repo.DeleteSong(id)
}

func (s *MainService) GetSongByID(id int) (model.Song, error) {
	s.log.Info("Gettting song by id", logrus.Fields{"id": id})
	return s.repo.GetSongByID(id)
}

func (s *MainService) GetSongText(id, page, perPage int) (string, error) {
	s.log.Info("Getting a song text", logrus.Fields{"id": id})

	song, err := s.repo.GetSongByID(id)
	if err != nil {
		s.log.Error("get song text", logrus.Fields{"error": err.Error()})
	}

	verses := strings.Split(song.Lyrics, "\n\n")
	if len(verses) == 0 {
		return "", nil
	}

	start := (page - 1) * perPage
	end := start + perPage

	if start >= len(verses) {
		return "", nil
	}

	if end > len(verses) {
		end = len(verses)
	}

	return strings.Join(verses[start:end], "\n\n"), nil

}
