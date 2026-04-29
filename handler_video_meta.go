package main

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"github.com/bootdotdev/learn-file-storage-s3-golang-starter/internal/auth"
	"github.com/bootdotdev/learn-file-storage-s3-golang-starter/internal/database"
	"github.com/google/uuid"
)

func (cfg *apiConfig) handlerVideoMetaCreate(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Title       string `json:"title"`
		Description string `json:"description"`
	}
	type response struct {
		Video
	}

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

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err = decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters", err)
		return
	}

	createVideoParams := database.CreateVideoParams{
		ID:    uuid.New().String(),
		Title: params.Title,
		Description: sql.NullString{
			String: params.Description,
			Valid:  true,
		},
		UserID: sql.NullString{
			String: userID,
			Valid:  true,
		},
	}
	dbVideo, err := cfg.db.CreateVideo(r.Context(), createVideoParams)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't create video", err)
		return
	}

	respondWithJSON(w, http.StatusCreated, response{
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

func (cfg *apiConfig) handlerVideoMetaDelete(w http.ResponseWriter, r *http.Request) {
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

	video, err := cfg.db.GetVideo(r.Context(), videoIDString)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "Couldn't get video", err)
		return
	}
	if video.UserID.String != userID {
		respondWithError(w, http.StatusForbidden, "You can't delete this video", err)
		return
	}

	err = cfg.db.DeleteVideo(r.Context(), videoIDString)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't delete video", err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (cfg *apiConfig) handlerVideoGet(w http.ResponseWriter, r *http.Request) {
	type response struct {
		Video
	}

	videoIDString := r.PathValue("videoID")
	dbVideo, err := cfg.db.GetVideo(r.Context(), videoIDString)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "Couldn't get video", err)
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

func (cfg *apiConfig) handlerVideosRetrieve(w http.ResponseWriter, r *http.Request) {
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

	dbVideos, err := cfg.db.GetVideos(r.Context(), sql.NullString{
		String: userID,
		Valid:  true,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't retrieve videos", err)
		return
	}

	videos := []Video{}

	for _, dbVideo := range dbVideos {
		v := Video{
			ID:           dbVideo.ID,
			CreatedAt:    dbVideo.CreatedAt.Time,
			UpdatedAt:    dbVideo.UpdatedAt.Time,
			Title:        dbVideo.Title,
			Description:  dbVideo.Description.String,
			ThumbnailURL: &dbVideo.ThumbnailUrl.String,
			VideoURL:     dbVideo.VideoUrl,
			UserID:       dbVideo.UserID.String,
		}
		videos = append(videos, v)
	}

	respondWithJSON(w, http.StatusOK, videos)
}
