package search

import (
	"net/http"

	"serica-go/internal/data"
	u "serica-go/internal/module/user"
)

type Handler struct {
	searchRepo *data.SearchHistoryRepo
}

func NewHandler(searchRepo *data.SearchHistoryRepo) *Handler {
	return &Handler{searchRepo: searchRepo}
}

// @Summary      获取搜索历史
// @Tags         Search
// @Produce      json
// @Success      200  {array}   map[string]interface{}
// @Security     BearerAuth
// @Router       /v1/client/search/history [get]
func (h *Handler) GetHistory(w http.ResponseWriter, r *http.Request) (any, error) {
	userID := u.GetUserID(r)
	histories, _ := h.searchRepo.FindByUser(userID, 10)
	return histories, nil
}

func (h *Handler) ClearHistory(w http.ResponseWriter, r *http.Request) (any, error) {
	userID := u.GetUserID(r)
	h.searchRepo.DeleteAll(userID)
	return map[string]bool{"ok": true}, nil
}

// @Summary      获取热门搜索
// @Tags         Search
// @Produce      json
// @Success      200  {array}   map[string]interface{}
// @Router       /v1/client/search/hot [get]
func (h *Handler) GetHot(w http.ResponseWriter, r *http.Request) (any, error) {
	keywords, _ := h.searchRepo.FindHotKeywords(10)
	return keywords, nil
}

func (h *Handler) GetHotWords(w http.ResponseWriter, r *http.Request) (any, error) {
	top := u.QueryInt(r, "top", 10)
	keywords, _ := h.searchRepo.FindHotKeywords(top)
	return map[string]interface{}{
		"items": keywords,
		"total": len(keywords),
	}, nil
}

func (h *Handler) GetRecommendWords(w http.ResponseWriter, r *http.Request) (any, error) {
	top := u.QueryInt(r, "top", 10)
	keywords, _ := h.searchRepo.FindHotKeywords(top)
	return keywords, nil
}
