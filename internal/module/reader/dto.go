package reader

type ReaderConfigReq struct {
	Theme    *string `json:"theme,omitempty"`
	FontSize *int    `json:"fontSize,omitempty"`
}

type CreateBookmarkReq struct {
	BookID   int64  `json:"bookId" validate:"required"`
	Title    string `json:"title" validate:"required"`
	Position string `json:"position" validate:"required"`
}

type CreateAnnotationReq struct {
	BookID   int64  `json:"bookId" validate:"required"`
	Content  string `json:"text" validate:"required"`
	Color    string `json:"color"`
	Position string `json:"position"`
	Note     string `json:"note"`
}

type UpdateAnnotationReq struct {
	Note  *string `json:"note,omitempty"`
	Color *string `json:"color,omitempty"`
}
