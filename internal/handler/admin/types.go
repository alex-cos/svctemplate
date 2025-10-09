package admin

type LoglevelRequest struct {
	Level string `json:"level" xml:"level" binding:"required"`
}
