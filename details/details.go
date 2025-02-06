package details

// SongDetail представляет подробную информацию о песне.
type SongDetail struct {
    ID          uint   `gorm:"primary key"` // Уникальный идентификатор песни
    ReleaseDate string `json:"releaseDate" validate:"required"` // Дата выпуска песни
    Text        string `json:"text" validate:"required"` // Текст песни
    Link        string `json:"link" validate:"required"` // Ссылка на песню
}

// Song представляет основную информацию о песне.
type Song struct {
    ID    uint   `json:"id"` // Уникальный идентификатор
    Group string `json:"group" validate:"required"` // Исполнитель или группа
    Song  string `json:"song" validate:"required"` // Название песни
}