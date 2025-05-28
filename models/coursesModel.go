package models

type Course struct {
    ID          int64  `json:"id"`
    Name        string `json:"name" validate:"required"`
    Description string `json:"description"`
    Duration    string `json:"duration"` 
    Credits     int    `json:"credits"`
	Price		int	   `json:"price"`
}
