package requests

type GetBoxesByGtinRequest struct {
	Gtin string `json:"gtin" binding:"required,gtin"`
}
