package handlers

import (
	"net/http"
	"strconv"

	"github.com/RodBarenco/colab-project-api/db"
	"github.com/RodBarenco/colab-project-api/res"
	"github.com/go-chi/chi"
)

// SaveImageToDBHandler saves the image to the database and returns the corresponding URL.
// Note: This handler does not have a direct route and is intended to be called by other functions.
func SaveImageToDBHandler(imageBase64 string) (string, error) {
	imageID, err := db.SaveImageToDB(dbAccessor, imageBase64)
	if err != nil {
		return "", err
	}

	imageURL := "image/" + strconv.Itoa(int(imageID))

	return imageURL, nil
}

func GetImageBase64ByIDHandler(w http.ResponseWriter, r *http.Request) {
	imageIDFromURL := chi.URLParam(r, "id")
	imageID, err := strconv.ParseUint(imageIDFromURL, 10, 32) // 32-bit unsigned integer
	if err != nil {
		RespondWithError(w, http.StatusBadRequest, "Invalid Image ID")
		return
	}

	imageBase64, err := db.GetImageBase64ByID(dbAccessor, uint(imageID))
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, "Failed to retrieve image")
		return
	}

	response := res.ImageRetrievedRes{
		ImageBase64: imageBase64,
		Message:     "Image retrieved successfully",
	}

	RespondWithJSON(w, http.StatusOK, response) // process image on the client side
}
