package reader

import (
	"encoding/json"
	"net/http"

	"serica-go/internal/data"
	u "serica-go/internal/module/user"
	"serica-go/internal/pkg/httputil"
)

type Handler struct {
	userRepo *data.UserRepo
	bookRepo *data.BookRepo
}

func NewHandler(userRepo *data.UserRepo, bookRepo *data.BookRepo) *Handler {
	return &Handler{userRepo: userRepo, bookRepo: bookRepo}
}

// @Summary      获取阅读器配置
// @Tags         Reader
// @Produce      json
// @Success      200  {object}  map[string]interface{}
// @Security     BearerAuth
// @Router       /v1/client/reader/config [get]
func (h *Handler) GetConfig(w http.ResponseWriter, r *http.Request) {
	userID := u.GetUserID(r)
	cfg, err := h.userRepo.GetReaderConfig(userID)
	if err != nil {
		httputil.RespondJSON(w, 200, map[string]interface{}{"theme": "light", "fontSize": 16})
		return
	}
	httputil.RespondJSON(w, 200, cfg)
}

// @Summary      更新阅读器配置
// @Tags         Reader
// @Accept       json
// @Produce      json
// @Param        body body map[string]interface{} true "配置信息"
// @Success      200  {object}  map[string]interface{}
// @Security     BearerAuth
// @Router       /v1/client/reader/config [patch]
func (h *Handler) UpdateConfig(w http.ResponseWriter, r *http.Request) {
	userID := u.GetUserID(r)
	var input map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		httputil.RespondJSON(w, 400, map[string]string{"error": "invalid body"})
		return
	}
	cfg := &data.ReaderConfig{UserID: userID}
	if theme, ok := input["theme"].(string); ok {
		cfg.Theme = theme
	}
	if fontSize, ok := input["fontSize"].(float64); ok {
		cfg.FontSize = int(fontSize)
	}
	if err := h.userRepo.SaveReaderConfig(cfg); err != nil {
		httputil.RespondJSON(w, 500, map[string]string{"error": err.Error()})
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

type CreateBookmarkInput struct {
	BookID   int64  `json:"bookId"`
	Title    string `json:"title"`
	Position string `json:"position"`
}

func (h *Handler) CreateBookmark(w http.ResponseWriter, r *http.Request) {
	userID := u.GetUserID(r)
	var input CreateBookmarkInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		httputil.RespondJSON(w, 400, map[string]string{"error": "invalid body"})
		return
	}
	bm := &data.Bookmark{
		UserID:   userID,
		BookID:   input.BookID,
		Title:    input.Title,
		Position: input.Position,
	}
	if err := h.userRepo.CreateBookmark(bm); err != nil {
		httputil.RespondJSON(w, 500, map[string]string{"error": err.Error()})
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

type CreateAnnotationInput struct {
	BookID   int64  `json:"bookId"`
	Content  string `json:"text"`
	Color    string `json:"color"`
	Position string `json:"position"`
	Note     string `json:"note"`
}

func (h *Handler) CreateAnnotation(w http.ResponseWriter, r *http.Request) {
	userID := u.GetUserID(r)
	var input CreateAnnotationInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		httputil.RespondJSON(w, 400, map[string]string{"error": "invalid body"})
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
		httputil.RespondJSON(w, 500, map[string]string{"error": err.Error()})
		return
	}
	httputil.RespondJSON(w, 200, a)
}

type UpdateAnnotationInput struct {
	Note  *string `json:"note,omitempty"`
	Color *string `json:"color,omitempty"`
}

func (h *Handler) UpdateAnnotation(w http.ResponseWriter, r *http.Request) {
	id := u.PathInt(r, "annotationId")
	var input UpdateAnnotationInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		httputil.RespondJSON(w, 400, map[string]string{"error": "invalid body"})
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
		httputil.RespondJSON(w, 500, map[string]string{"error": err.Error()})
		return
	}
	httputil.RespondJSON(w, 200, map[string]bool{"ok": true})
}

func (h *Handler) DeleteAnnotation(w http.ResponseWriter, r *http.Request) {
	id := u.PathInt(r, "annotationId")
	h.userRepo.DeleteAnnotation(id)
	httputil.RespondJSON(w, 200, map[string]bool{"ok": true})
}
