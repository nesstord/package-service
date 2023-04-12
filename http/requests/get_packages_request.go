package requests

type GetPackagesRequest struct {
	Gtin string `json:"gtin" binding:"required,gtin"`
}
