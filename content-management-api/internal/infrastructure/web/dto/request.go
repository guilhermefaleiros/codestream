package dto

type CreateEmptyMovieRequestDTO struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Genre       string `json:"genre"`
	LaunchYear  int    `json:"launch_year"`
	Duration    int    `json:"duration"`
}
