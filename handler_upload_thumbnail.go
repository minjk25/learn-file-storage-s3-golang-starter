package main

import (
	"database/sql"
	"io"
	"mime"
	"net/http"
	"os"
	"time"

	"github.com/bootdotdev/learn-file-storage-s3-golang-starter/internal/auth"
	"github.com/bootdotdev/learn-file-storage-s3-golang-starter/internal/database"
)

type Video struct {
	ID           string      `json:"id"`
	CreatedAt    time.Time   `json:"created_at"`
	UpdatedAt    time.Time   `json:"updated_at"`
	Title        string      `json:"title"`
	Description  string      `json:"description"`
	ThumbnailURL *string     `json:"thumbnail_url"`
	VideoURL     interface{} `json:"video_url"`
	UserID       string      `json:"user_id"`
}

func (cfg *apiConfig) handlerUploadThumbnail(w http.ResponseWriter, r *http.Request) {
	type response struct {
		Video
	}
	videoIDString := r.PathValue("videoID")

	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Couldn't find JWT", err)
		return
	}

	userID, err := auth.ValidateJWT(token, cfg.jwtSecret)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Couldn't validate JWT", err)
		return
	}

	const maxMemory = 10 << 20 // 10MB
	r.ParseMultipartForm(maxMemory)
	file, header, err := r.FormFile("thumbnail")
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Unable to parse form file", err)
		return
	}
	defer file.Close()

	typeContent, _, err := mime.ParseMediaType(header.Header.Get("Content-Type"))
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid Content-Type", err)
		return
	}
	if typeContent == "" {
		respondWithError(w, http.StatusBadRequest, "Missing Content-Type for thumbnail", nil)
		return
	}
	if typeContent != "image/jpeg" && typeContent != "image/png" {
		respondWithError(w, http.StatusBadRequest, "Invalid file type", nil)
		return
	}

	dbVideo, err := cfg.db.GetVideo(r.Context(), videoIDString)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "the video doesn't exist", err)
		return
	}
	if dbVideo.UserID.String != userID {
		respondWithError(w, http.StatusUnauthorized, "Not authorized to update this video", err)
		return
	}

	fileName, err := getAssetPath(typeContent)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Unable to create file on server", err)
		return
	}
	filePath := cfg.getAssetDiskPath(fileName)
	createdFile, err := os.Create(filePath)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Unable to create file on server", err)
		return
	}
	defer createdFile.Close()

	_, err = io.Copy(createdFile, file)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error saving file", err)
		return
	}

	thumbnailURL := cfg.getAssetURL(fileName)
	dbVideo.ThumbnailUrl = sql.NullString{
		String: thumbnailURL,
		Valid:  true,
	}
	updateVideoParams := database.UpdateVideoParams{
		Title:        dbVideo.Title,
		Description:  dbVideo.Description,
		ThumbnailUrl: dbVideo.ThumbnailUrl,
		VideoUrl:     dbVideo.VideoUrl,
		UserID:       dbVideo.UserID,
		ID:           dbVideo.ID,
	}

	err = cfg.db.UpdateVideo(r.Context(), updateVideoParams)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't update video", err)
		return
	}

	respondWithJSON(w, http.StatusOK, response{
		Video: Video{
			ID:           dbVideo.ID,
			CreatedAt:    dbVideo.CreatedAt.Time,
			UpdatedAt:    dbVideo.UpdatedAt.Time,
			Title:        dbVideo.Title,
			Description:  dbVideo.Description.String,
			ThumbnailURL: &dbVideo.ThumbnailUrl.String,
			VideoURL:     dbVideo.VideoUrl,
			UserID:       dbVideo.UserID.String,
		},
	})
}
