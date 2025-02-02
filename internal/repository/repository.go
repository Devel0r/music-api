package repository

import (
    "database/sql"
    "fmt"
    "time"

    "music/internal/model"
)

type MainRepository struct {
    db *sql.DB
}

func NewMainRepository(db *sql.DB) *MainRepository {
    return &MainRepository{
        db: db,
    }
}

func (m *MainRepository) GetAllSongs(filters map[string]string, limit, offset int) ([]model.Song, error) {
    var songs []model.Song
    
    query := `
        SELECT id, group_name, song_title, release_date, lyrics, youtube_link, created_at
        FROM songs
        WHERE 1=1
    `
    
    args := make([]interface{}, 0)
    paramCounter := 1
    
    validFilters := map[string]string{
        "group":    "group_name",
        "song":     "song_title",
        "release":  "release_date",
        "lyrics":   "lyrics",
        "link":     "youtube_link",
    }
    
    for param, column := range validFilters {
        if value, exists := filters[param]; exists && value != "" {
            query += fmt.Sprintf(" AND %s = $%d", column, paramCounter)
            args = append(args, value)
            paramCounter++
        }
    }
    
    query += fmt.Sprintf(" LIMIT $%d OFFSET $%d", paramCounter, paramCounter+1)
    args = append(args, limit, offset)
    
    
    rows, err := m.db.Query(query, args...)
    if err != nil {
        return nil, fmt.Errorf("query error: %w", err)
    }
    defer rows.Close()
    
    
    for rows.Next() {
        var song model.Song
        if err := rows.Scan(&song.ID, &song.GroupName, &song.SongTitle, 
            &song.ReleaseDate, &song.Lyrics, &song.YouTubeLink, &song.CreatedAt); err != nil {
            return nil, fmt.Errorf("scan error: %w", err)
        }
        songs = append(songs, song)
    }
    
    if err := rows.Err(); err != nil {
        return nil, fmt.Errorf("rows error: %w", err)
    }
    
    return songs, nil
}

func (m *MainRepository) GetSongByID(id int) (model.Song, error) {
    var song model.Song
    query := `
        SELECT id, group_name, song_title, release_date, lyrics, youtube_link, created_at
        FROM songs
        WHERE id = $1
    `
    row := m.db.QueryRow(query, id)
    if err := row.Scan(
        &song.ID,
        &song.GroupName,
        &song.SongTitle,
        &song.ReleaseDate,
        &song.Lyrics,
        &song.YouTubeLink,
        &song.CreatedAt,
    ); err != nil {
        if err == sql.ErrNoRows {
            return song, fmt.Errorf("song not found")
        }
        return song, err
    }
    return song, nil
}

func (m *MainRepository) AddSong(song model.Song) (int, error) {
    var id int
    query := `
        INSERT INTO songs (group_name, song_title, release_date, lyrics, youtube_link, created_at)
        VALUES ($1, $2, $3, $4, $5, $6)
        RETURNING id
    `
    err := m.db.QueryRow(query,
        song.GroupName,
        song.SongTitle,
        song.ReleaseDate,
        song.Lyrics,
        song.YouTubeLink,
        time.Now(),
    ).Scan(&id)
    if err != nil {
        return 0, err
    }
    return id, nil
}

func (m *MainRepository) UpdateSong(id int, song model.Song) error {
    query := `
        UPDATE songs
        SET group_name = $1, song_title = $2, release_date = $3, lyrics = $4, youtube_link = $5
        WHERE id = $6
    `
    _, err := m.db.Exec(query,
        song.GroupName,
        song.SongTitle,
        song.ReleaseDate,
        song.Lyrics,
        song.YouTubeLink,
        time.Now(),
        id,
    )
    if err != nil {
        return err
    }
    return nil
}

func (m *MainRepository) DeleteSong(id int) error {
    query := `DELETE FROM songs WHERE id = $1`
    _, err := m.db.Exec(query, id)
    if err != nil {
        return err
    }
    return nil
}