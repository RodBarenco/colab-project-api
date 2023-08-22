package res

// ImageSavedRes defines the response structure for a saved image.
type ImageSavedRes struct {
	ImageID uint   `json:"image_id"`
	Message string `json:"message"`
}

// ImageRetrievedRes defines the response structure for a retrieved image.
type ImageRetrievedRes struct {
	ImageBase64 string `json:"image_base64"`
	Message     string `json:"message"`
}
