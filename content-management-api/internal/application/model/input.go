package model

type CreateVideoInput struct {
	Title            string
	Description      string
	ContentType      string
	FileExtension    string
	OriginalFileName string
	File             []byte
}
