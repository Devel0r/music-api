package model

import "time"

type Song struct {
    ID          int       `json:"id" example:"1"`
    GroupName   string    `json:"group_name" example:"Muse"`
    SongTitle   string    `json:"song_title" example:"Supermassive Black Hole"`
    ReleaseDate string    `json:"release_date" example:"16.07.2006"`
    Lyrics      string    `json:"lyrics" example:"Ooh baby, don't you know I suffer?..."`
    YouTubeLink string    `json:"youtube_link" example:"https://youtu.be/Xsp3_a-PMTw"`
    CreatedAt   time.Time `json:"created_at" example:"2024-01-01T12:00:00Z"`
}

type SongDetail struct {
    ReleaseDate string `json:"releaseDate" example:"16.07.2006"`
    Text        string `json:"text" example:"Ooh baby, don't you know I suffer?..."`
    Link        string `json:"link" example:"https://youtu.be/Xsp3_a-PMTw"`
}
