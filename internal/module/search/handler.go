package search

import (
	"net/http"

	"serica-go/internal/data"
	u "serica-go/internal/module/user"
	"serica-go/internal/pkg/httputil"
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
func (h *Handler) GetHistory(w http.ResponseWriter, r *http.Request) {
	userID := u.GetUserID(r)
	histories, _ := h.searchRepo.FindByUser(userID, 10)
	httputil.RespondJSON(w, 200, histories)
}

func (h *Handler) ClearHistory(w http.ResponseWriter, r *http.Request) {
	userID := u.GetUserID(r)
	h.searchRepo.DeleteAll(userID)
	httputil.RespondJSON(w, 200, map[string]bool{"ok": true})
}

// @Summary      获取热门搜索
// @Tags         Search
// @Produce      json
// @Success      200  {array}   map[string]interface{}
// @Router       /v1/client/search/hot [get]
func (h *Handler) GetHot(w http.ResponseWriter, r *http.Request) {
	keywords, _ := h.searchRepo.FindHotKeywords(10)
	httputil.RespondJSON(w, 200, keywords)
}

func (h *Handler) GetHotWords(w http.ResponseWriter, r *http.Request) {
	top := u.QueryInt(r, "top", 10)
	keywords, _ := h.searchRepo.FindHotKeywords(top)
	httputil.RespondJSON(w, 200, map[string]interface{}{
		"items": keywords,
		"total": len(keywords),
	})
}

func (h *Handler) GetRecommendWords(w http.ResponseWriter, r *http.Request) {
	top := u.QueryInt(r, "top", 10)
	keywords, _ := h.searchRepo.FindHotKeywords(top)
	httputil.RespondJSON(w, 200, keywords)
}
