package reader

import (
	"net/http"

	"serica-go/internal/data"
	u "serica-go/internal/module/user"
	"serica-go/internal/pkg/exception"
	"serica-go/internal/pkg/httputil"
)

type Handler struct {
	userRepo *data.UserRepo
	bookRepo *data.BookRepo
}

func NewHandler(userRepo *data.UserRepo, bookRepo *data.BookRepo) *Handler {
	return &Handler{userRepo: userRepo, bookRepo: bookRepo}
}

func (h *Handler) GetConfig(w http.ResponseWriter, r *http.Request) {
	userID := u.GetUserID(r)
	cfg, err := h.userRepo.GetReaderConfig(userID)
	if err != nil {
		httputil.RespondJSON(w, 200, map[string]interface{}{"theme": "light", "fontSize": 16})
		return
	}
	httputil.RespondJSON(w, 200, cfg)
}

func (h *Handler) UpdateConfig(w http.ResponseWriter, r *http.Request) {
	userID := u.GetUserID(r)
	var input ReaderConfigReq
	if err := httputil.DecodeAndValidate(r, &input); err != nil {
		exception.InvalidBody.Write(w)
		return
	}
	cfg := &data.ReaderConfig{UserID: userID}
	if input.Theme != nil {
		cfg.Theme = *input.Theme
	}
	if input.FontSize != nil {
		cfg.FontSize = *input.FontSize
	}
	if err := h.userRepo.SaveReaderConfig(cfg); err != nil {
		exception.InternalServerError.Write(w, err.Error())
		return
	}
	httputil.RespondJSON(w, 200, cfg)
}

func (h *Handler) GetBookmarks(w http.ResponseWriter, r *http.Request) {
	userID := u.GetUserID(r)
	bookID := u.PathInt(r, "id")
	bookmarks, _ := h.userRepo.FindBookmarks(userID, bookID)
	httputil.RespondJSON(w, 200, bookmarks)
}

func (h *Handler) CreateBookmark(w http.ResponseWriter, r *http.Request) {
	userID := u.GetUserID(r)
	var input CreateBookmarkReq
	if err := httputil.DecodeAndValidate(r, &input); err != nil {
		exception.InvalidBody.Write(w)
		return
	}
	bm := &data.Bookmark{
		UserID:   userID,
		BookID:   input.BookID,
		Title:    input.Title,
		Position: input.Position,
	}
	if err := h.userRepo.CreateBookmark(bm); err != nil {
		exception.InternalServerError.Write(w, err.Error())
		return
	}
	httputil.RespondJSON(w, 200, bm)
}

func (h *Handler) GetAnnotations(w http.ResponseWriter, r *http.Request) {
	userID := u.GetUserID(r)
	bookID := u.PathInt(r, "id")
	annotations, _ := h.userRepo.FindAnnotations(userID, bookID)
	httputil.RespondJSON(w, 200, annotations)
}

func (h *Handler) CreateAnnotation(w http.ResponseWriter, r *http.Request) {
	userID := u.GetUserID(r)
	var input CreateAnnotationReq
	if err := httputil.DecodeAndValidate(r, &input); err != nil {
		exception.InvalidBody.Write(w)
		return
	}
	a := &data.Annotation{
		UserID:   userID,
		BookID:   input.BookID,
		Content:  input.Content,
		Color:    input.Color,
		Position: input.Position,
	}
	if err := h.userRepo.CreateAnnotation(a); err != nil {
		exception.InternalServerError.Write(w, err.Error())
		return
	}
	httputil.RespondJSON(w, 200, a)
}

func (h *Handler) UpdateAnnotation(w http.ResponseWriter, r *http.Request) {
	id := u.PathInt(r, "annotationId")
	var input UpdateAnnotationReq
	if err := httputil.DecodeAndValidate(r, &input); err != nil {
		exception.InvalidBody.Write(w)
		return
	}
	updates := make(map[string]interface{})
	if input.Note != nil {
		updates["content"] = *input.Note
	}
	if input.Color != nil {
		updates["color"] = *input.Color
	}
	if err := h.userRepo.UpdateAnnotation(id, updates); err != nil {
		exception.InternalServerError.Write(w, err.Error())
		return
	}
	httputil.RespondJSON(w, 200, map[string]bool{"ok": true})
}

func (h *Handler) DeleteAnnotation(w http.ResponseWriter, r *http.Request) {
	id := u.PathInt(r, "annotationId")
	h.userRepo.DeleteAnnotation(id)
	httputil.RespondJSON(w, 200, map[string]bool{"ok": true})
}
