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

func (h *Handler) GetConfig(w http.ResponseWriter, r *http.Request) (any, error) {
	userID := u.GetUserID(r)
	cfg, err := h.userRepo.GetReaderConfig(userID)
	if err != nil {
		return map[string]interface{}{"theme": "light", "fontSize": 16}, nil
	}
	return cfg, nil
}

func (h *Handler) UpdateConfig(w http.ResponseWriter, r *http.Request) (any, error) {
	userID := u.GetUserID(r)
	var input ReaderConfigReq
	if err := httputil.DecodeAndValidate(r, &input); err != nil {
		return nil, exception.InvalidBody
	}
	cfg := &data.ReaderConfig{UserID: userID}
	if input.Theme != nil {
		cfg.Theme = *input.Theme
	}
	if input.FontSize != nil {
		cfg.FontSize = *input.FontSize
	}
	if err := h.userRepo.SaveReaderConfig(cfg); err != nil {
		return nil, exception.InternalServerError.WithMsg(err.Error())
	}
	return cfg, nil
}

func (h *Handler) GetBookmarks(w http.ResponseWriter, r *http.Request) (any, error) {
	userID := u.GetUserID(r)
	bookID := u.PathInt(r, "id")
	bookmarks, _ := h.userRepo.FindBookmarks(userID, bookID)
	return bookmarks, nil
}

func (h *Handler) CreateBookmark(w http.ResponseWriter, r *http.Request) (any, error) {
	userID := u.GetUserID(r)
	var input CreateBookmarkReq
	if err := httputil.DecodeAndValidate(r, &input); err != nil {
		return nil, exception.InvalidBody
	}
	bm := &data.Bookmark{
		UserID:   userID,
		BookID:   input.BookID,
		Title:    input.Title,
		Position: input.Position,
	}
	if err := h.userRepo.CreateBookmark(bm); err != nil {
		return nil, exception.InternalServerError.WithMsg(err.Error())
	}
	return bm, nil
}

func (h *Handler) GetAnnotations(w http.ResponseWriter, r *http.Request) (any, error) {
	userID := u.GetUserID(r)
	bookID := u.PathInt(r, "id")
	annotations, _ := h.userRepo.FindAnnotations(userID, bookID)
	return annotations, nil
}

func (h *Handler) CreateAnnotation(w http.ResponseWriter, r *http.Request) (any, error) {
	userID := u.GetUserID(r)
	var input CreateAnnotationReq
	if err := httputil.DecodeAndValidate(r, &input); err != nil {
		return nil, exception.InvalidBody
	}
	a := &data.Annotation{
		UserID:   userID,
		BookID:   input.BookID,
		Content:  input.Content,
		Color:    input.Color,
		Position: input.Position,
	}
	if err := h.userRepo.CreateAnnotation(a); err != nil {
		return nil, exception.InternalServerError.WithMsg(err.Error())
	}
	return a, nil
}

func (h *Handler) UpdateAnnotation(w http.ResponseWriter, r *http.Request) (any, error) {
	id := u.PathInt(r, "annotationId")
	var input UpdateAnnotationReq
	if err := httputil.DecodeAndValidate(r, &input); err != nil {
		return nil, exception.InvalidBody
	}
	updates := make(map[string]interface{})
	if input.Note != nil {
		updates["content"] = *input.Note
	}
	if input.Color != nil {
		updates["color"] = *input.Color
	}
	if err := h.userRepo.UpdateAnnotation(id, updates); err != nil {
		return nil, exception.InternalServerError.WithMsg(err.Error())
	}
	return map[string]bool{"ok": true}, nil
}

func (h *Handler) DeleteAnnotation(w http.ResponseWriter, r *http.Request) (any, error) {
	id := u.PathInt(r, "annotationId")
	h.userRepo.DeleteAnnotation(id)
	return map[string]bool{"ok": true}, nil
}
