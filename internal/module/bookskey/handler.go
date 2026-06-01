package bookskey

import (
	"net/http"

	"serica-go/internal/data"
	u "serica-go/internal/module/user"
	"serica-go/internal/pkg/httputil"
)

type Handler struct {
	bookKeyRepo *data.BookKeyRepo
}

func NewHandler(bookKeyRepo *data.BookKeyRepo) *Handler {
	return &Handler{bookKeyRepo: bookKeyRepo}
}

// @Summary      获取书籍AES密钥
// @Description  返回加密后的书籍密钥
// @Tags         BooksKey
// @Produce      json
// @Param        bookId query string true "书籍ID"
// @Success      200  {object}  map[string]interface{}
// @Security     BearerAuth
// @Router       /v1/client/booksKey [get]
func (h *Handler) GetKey(w http.ResponseWriter, r *http.Request) {
	bookID := u.PathInt(r, "bookId")
	key, err := h.bookKeyRepo.FindByBookID(bookID)
	if err != nil {
		httputil.RespondJSON(w, 404, map[string]string{"error": "key not found"})
		return
	}
	httputil.RespondJSON(w, 200, key)
}
