package response

import "github.com/google/uuid"

type FileUploadResponse struct {
	ID         uuid.UUID `json:"id"`
	Name       string    `json:"name"`
	Location   string    `json:"location"`
	BucketName string    `json:"bucket_name"`
	URL        string    `json:"url"`
}
