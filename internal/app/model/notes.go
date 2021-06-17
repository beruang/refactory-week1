package model

type Notes struct {
	Id       int
	UserId   int
	Type     string
	Title    string
	Body     string
	Secret   string
	IsActive string
}

func NewNotes(id int, userId int,
	types string, title string, body string, secret string) *Notes {
	return &Notes{Id: id, UserId: userId, Type: types, Title: title, Body: body, Secret: secret}
}

type NotesRequest struct {
	Type   string `json:"type" validate:"required"`
	Title  string `json:"title" validate:"required"`
	Body   string `json:"body" validate:"required"`
	Secret string `json:"secret"`
}

type NotesResponse struct {
	Id     int    `json:"id"`
	Type   string `json:"type"`
	Title  string `json:"title"`
	Body   string `json:"body"`
	Secret string `json:"secret"`
}

func NewNotesResponse(id int,
	types string, title string, body string, secret string) *NotesResponse {
	return &NotesResponse{Id: id, Type: types, Title: title, Body: body, Secret: secret}
}

type SecretRequest struct {
	Secret string `json:"secret"`
}
