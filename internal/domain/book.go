package domain

type Book struct {
	ID     uint   `json:"id" gorm:"primaryKey"`
	Title  string `json:"title"`
	Author string `json:"author"`
}
