package user

import (
	"encoding/json"
	"net/http"

	"serica-go/internal/data"
	"serica-go/internal/pkg/httputil"
	"serica-go/internal/utl/rand"
)

type Handler struct {
	userRepo *data.UserRepo
	redis    *data.RedisClient
}

func NewHandler(userRepo *data.UserRepo, redis *data.RedisClient) *Handler {
	return &Handler{userRepo: userRepo, redis: redis}
}

type LoginInput struct {
	Email string `json:"email"`
}

const sessionKeyPrefix = "session:token:"
const sessionTTL = 7 * 24 * 60 * 60

// @Summary      用户登录
// @Description  通过邮箱登录
// @Tags         User
// @Accept       json
// @Produce      json
// @Param        body body LoginInput true "登录信息"
// @Success      200  {object}  map[string]interface{}
// @Failure      400  {object}  map[string]string
// @Failure      401  {object}  map[string]string
// @Router       /v1/client/login [post]
func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	var input LoginInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		httputil.RespondJSON(w, 400, map[string]string{"error": "invalid body"})
		return
	}

	user, err := h.userRepo.FindByEmail(input.Email)
	if err != nil {
		httputil.RespondJSON(w, 401, map[string]string{"error": "user not found"})
		return
	}

	token := rand.Token(32)
	ip := r.Header.Get("x-forwarded-for")
	ua := r.Header.Get("user-agent")
	dv := r.Header.Get("x-device")

	h.userRepo.CreateSession(&data.UserSession{
		UserID:    user.ID,
		Token:     token,
		IP:        ip,
		UserAgent: ua,
		Device:    dv,
	})

	if h.redis != nil {
		h.redis.SetJSON(r.Context(), sessionKeyPrefix+token, map[string]interface{}{
			"userId": user.ID,
			"email":  user.Email,
		}, sessionTTL)
	}

	httputil.RespondJSON(w, 200, map[string]interface{}{
		"token": token,
		"user":  user,
	})
}

// @Summary      用户注册
// @Description  通过邮箱注册/登录
// @Tags         User
// @Accept       json
// @Produce      json
// @Param        body body LoginInput true "注册信息"
// @Success      200  {object}  map[string]string
// @Router       /v1/client/user/register [post]
func (h *Handler) Register(w http.ResponseWriter, r *http.Request) {
	var input LoginInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		httputil.RespondJSON(w, 400, map[string]string{"error": "invalid body"})
		return
	}

	user, err := h.userRepo.FindByEmail(input.Email)
	if err != nil {
		user = &data.User{Email: input.Email}
		h.userRepo.CreateUser(user)
	}

	token := rand.Token(32)
	ip := r.Header.Get("x-forwarded-for")
	ua := r.Header.Get("user-agent")
	dv := r.Header.Get("x-device")

	h.userRepo.CreateSession(&data.UserSession{
		UserID:    user.ID,
		Token:     token,
		IP:        ip,
		UserAgent: ua,
		Device:    dv,
	})

	if h.redis != nil {
		h.redis.SetJSON(r.Context(), sessionKeyPrefix+token, map[string]interface{}{
			"userId": user.ID,
			"email":  user.Email,
		}, sessionTTL)
	}

	httputil.RespondJSON(w, 200, map[string]string{"token": token})
}

// @Summary      获取用户信息
// @Description  获取当前登录用户信息
// @Tags         User
// @Produce      json
// @Success      200  {object}  map[string]interface{}
// @Security     BearerAuth
// @Router       /v1/client/user/profile [get]
func (h *Handler) GetProfile(w http.ResponseWriter, r *http.Request) {
	userID := GetUserID(r)
	user, err := h.userRepo.FindByID(userID)
	if err != nil {
		httputil.RespondJSON(w, 404, map[string]string{"error": "user not found"})
		return
	}
	httputil.RespondJSON(w, 200, user)
}

type UpdateProfileInput struct {
	Name     *string `json:"name,omitempty"`
	Avatar   *string `json:"avatar,omitempty"`
	Language *string `json:"language,omitempty"`
	Offset   *string `json:"offset,omitempty"`
}

// @Summary      更新用户信息
// @Tags         User
// @Accept       json
// @Produce      json
// @Param        body body UpdateProfileInput true "更新信息"
// @Success      200  {object}  map[string]bool
// @Security     BearerAuth
// @Router       /v1/client/user/profile [patch]
func (h *Handler) UpdateProfile(w http.ResponseWriter, r *http.Request) {
	userID := GetUserID(r)
	var input UpdateProfileInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		httputil.RespondJSON(w, 400, map[string]string{"error": "invalid body"})
		return
	}

	updates := make(map[string]interface{})
	if input.Name != nil {
		updates["name"] = *input.Name
	}
	if input.Avatar != nil {
		updates["avatar"] = *input.Avatar
	}
	if input.Language != nil {
		updates["language"] = *input.Language
	}
	if input.Offset != nil {
		updates["offset"] = *input.Offset
	}

	if len(updates) == 0 {
		httputil.RespondJSON(w, 200, map[string]bool{"ok": true})
		return
	}

	if err := h.userRepo.Update(userID, updates); err != nil {
		httputil.RespondJSON(w, 500, map[string]string{"error": err.Error()})
		return
	}
	httputil.RespondJSON(w, 200, map[string]bool{"ok": true})
}

// @Summary      获取收藏列表
// @Description  分页获取当前用户的收藏书籍
// @Tags         User
// @Produce      json
// @Param        pageIndex  query  int  false  "页码"  default(1)
// @Param        pageSize   query  int  false  "每页数量"  default(12)
// @Param        desc       query  bool false  "是否倒序"  default(true)
// @Success      200  {object}  map[string]interface{}
// @Security     BearerAuth
// @Router       /v1/client/user/favourites [get]
func (h *Handler) GetFavourites(w http.ResponseWriter, r *http.Request) {
	userID := GetUserID(r)
	pageIndex := QueryInt(r, "pageIndex", 1)
	pageSize := QueryInt(r, "pageSize", 12)
	desc := QueryBool(r, "desc", true)

	offset := (pageIndex - 1) * pageSize
	items, total, err := h.userRepo.FindFavouritesPaginated(userID, offset, pageSize, desc)
	if err != nil {
		httputil.RespondJSON(w, 500, map[string]string{"error": err.Error()})
		return
	}

	httputil.RespondJSON(w, 200, map[string]interface{}{
		"items":     items,
		"total":     total,
		"pageIndex": pageIndex,
		"pageSize":  pageSize,
	})
}

// @Summary      收藏/取消收藏
// @Tags         User
// @Produce      json
// @Param        bookId path int true "书籍ID"
// @Success      200  {object}  map[string]bool
// @Security     BearerAuth
// @Router       /v1/client/books/{bookId}/favourite [post]
func (h *Handler) ToggleFavourite(w http.ResponseWriter, r *http.Request) {
	userID := GetUserID(r)
	bookID := PathInt(r, "bookId")
	liked, err := h.userRepo.ToggleFavourite(userID, bookID)
	if err != nil {
		httputil.RespondJSON(w, 500, map[string]string{"error": err.Error()})
		return
	}
	httputil.RespondJSON(w, 200, map[string]bool{"liked": liked})
}

type BatchCancelInput struct {
	BookIDs []int64 `json:"bookIds"`
}

// @Summary      批量取消收藏
// @Tags         User
// @Accept       json
// @Produce      json
// @Param        body body BatchCancelInput true "书籍ID列表"
// @Success      200  {object}  map[string]bool
// @Security     BearerAuth
// @Router       /v1/client/books/cancelFavourite [post]
func (h *Handler) BatchCancelFavourite(w http.ResponseWriter, r *http.Request) {
	userID := GetUserID(r)
	var input BatchCancelInput
	json.NewDecoder(r.Body).Decode(&input)
	h.userRepo.BatchCancelFavourite(userID, input.BookIDs)
	httputil.RespondJSON(w, 200, map[string]bool{"ok": true})
}

// @Summary      获取书签列表
// @Tags         User
// @Produce      json
// @Param        bookId path int true "书籍ID"
// @Success      200  {array}   map[string]interface{}
// @Security     BearerAuth
// @Router       /v1/client/books/{bookId}/bookmarks [get]
func (h *Handler) GetBookmarks(w http.ResponseWriter, r *http.Request) {
	userID := GetUserID(r)
	bookID := PathInt(r, "bookId")
	bookmarks, _ := h.userRepo.FindBookmarks(userID, bookID)
	httputil.RespondJSON(w, 200, bookmarks)
}

// @Summary      添加书签
// @Tags         User
// @Accept       json
// @Produce      json
// @Param        bookId path int true "书籍ID"
// @Param        body body map[string]interface{} true "书签信息"
// @Success      200  {object}  map[string]interface{}
// @Security     BearerAuth
// @Router       /v1/client/books/{bookId}/bookmarks [post]
func (h *Handler) CreateBookmark(w http.ResponseWriter, r *http.Request) {
	userID := GetUserID(r)
	bookID := PathInt(r, "bookId")
	var bm data.Bookmark
	json.NewDecoder(r.Body).Decode(&bm)
	bm.UserID = userID
	bm.BookID = bookID
	h.userRepo.CreateBookmark(&bm)
	httputil.RespondJSON(w, 200, bm)
}

// @Summary      删除书签
// @Tags         User
// @Produce      json
// @Param        bookmarkId path int true "书签ID"
// @Success      200  {object}  map[string]bool
// @Security     BearerAuth
// @Router       /v1/client/bookmarks/{bookmarkId} [delete]
func (h *Handler) DeleteBookmark(w http.ResponseWriter, r *http.Request) {
	id := PathInt(r, "bookmarkId")
	h.userRepo.DeleteBookmark(id)
	httputil.RespondJSON(w, 200, map[string]bool{"ok": true})
}

// @Summary      获取批注列表
// @Tags         User
// @Produce      json
// @Param        bookId path int true "书籍ID"
// @Success      200  {array}   map[string]interface{}
// @Security     BearerAuth
// @Router       /v1/client/books/{bookId}/annotations [get]
func (h *Handler) GetAnnotations(w http.ResponseWriter, r *http.Request) {
	userID := GetUserID(r)
	bookID := PathInt(r, "bookId")
	annotations, _ := h.userRepo.FindAnnotations(userID, bookID)
	httputil.RespondJSON(w, 200, annotations)
}

// @Summary      添加批注
// @Tags         User
// @Accept       json
// @Produce      json
// @Param        bookId path int true "书籍ID"
// @Param        body body map[string]interface{} true "批注信息"
// @Success      200  {object}  map[string]interface{}
// @Security     BearerAuth
// @Router       /v1/client/books/{bookId}/annotations [post]
func (h *Handler) CreateAnnotation(w http.ResponseWriter, r *http.Request) {
	userID := GetUserID(r)
	bookID := PathInt(r, "bookId")
	var a data.Annotation
	json.NewDecoder(r.Body).Decode(&a)
	a.UserID = userID
	a.BookID = bookID
	h.userRepo.CreateAnnotation(&a)
	httputil.RespondJSON(w, 200, a)
}

// @Summary      获取阅读器配置
// @Tags         User
// @Produce      json
// @Success      200  {object}  map[string]interface{}
// @Security     BearerAuth
// @Router       /v1/client/reader-config [get]
func (h *Handler) GetReaderConfig(w http.ResponseWriter, r *http.Request) {
	userID := GetUserID(r)
	cfg, err := h.userRepo.GetReaderConfig(userID)
	if err != nil {
		httputil.RespondJSON(w, 200, map[string]interface{}{"theme": "light", "fontSize": 16})
		return
	}
	httputil.RespondJSON(w, 200, cfg)
}

// @Summary      保存阅读器配置
// @Tags         User
// @Accept       json
// @Produce      json
// @Param        body body map[string]interface{} true "配置信息"
// @Success      200  {object}  map[string]interface{}
// @Security     BearerAuth
// @Router       /v1/client/reader-config [post]
func (h *Handler) SaveReaderConfig(w http.ResponseWriter, r *http.Request) {
	userID := GetUserID(r)
	var cfg data.ReaderConfig
	json.NewDecoder(r.Body).Decode(&cfg)
	cfg.UserID = userID
	h.userRepo.SaveReaderConfig(&cfg)
	httputil.RespondJSON(w, 200, cfg)
}
