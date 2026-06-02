package user

type LoginReq struct {
	Email string `json:"email" validate:"required,email"`
}

type UpdateProfileReq struct {
	Name     *string `json:"name,omitempty"`
	Avatar   *string `json:"avatar,omitempty"`
	Language *string `json:"language,omitempty"`
	Offset   *string `json:"offset,omitempty"`
}

type LoginResponse struct {
	Token string    `json:"token"`
	User  *UserInfo `json:"user"`
}

type UserInfo struct {
	ID       int64  `json:"id"`
	Email    string `json:"email"`
	Name     string `json:"name"`
	Avatar   string `json:"avatar"`
	Role     string `json:"role"`
	Language string `json:"language"`
	Offset   string `json:"offset"`
}

type BatchCancelFavouriteReq struct {
	BookIDs    *[]int64 `json:"bookIds,omitempty"`
	SelectAll  *bool    `json:"selectAll,omitempty"`
	ExcludeIDs *[]int64 `json:"excludeIds,omitempty"`
}

type BatchCancelFavouriteRes struct {
	Cancelled int64 `json:"cancelled"`
}

type FavouriteListRes struct {
	Items     []any `json:"items"`
	Total     int64 `json:"total"`
	PageIndex int   `json:"pageIndex"`
	PageSize  int   `json:"pageSize"`
}
