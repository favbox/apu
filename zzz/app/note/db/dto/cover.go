package dto

type Cover struct {
	FileID   string `json:"file_id"`
	FileName string `json:"file_name"`
	Width    int    `json:"width"`
	Height   int    `json:"height"`
}
