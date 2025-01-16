package model

type MusicRequest struct {
	IsAuto       int    `json:"is_auto"`
	Prompt       string `json:"prompt"`
	Lyrics       string `json:"lyrics"`
	Title        string `json:"title"`
	Instrumental int    `json:"instrumental"`
}

type MusicResponse struct {
	Status  int         `json:"status"`
	Message string      `json:"message"`
	Data    []MusicData `json:"data"`
}

type MusicData struct {
	AudioFile string `json:"audio_file"`
	ImageFile string `json:"image_file"`
	ItemUUID  string `json:"item_uuid"`
	Title     string `json:"title"`
	Lyric     string `json:"lyric"`
	Tags      string `json:"tags"`
}
