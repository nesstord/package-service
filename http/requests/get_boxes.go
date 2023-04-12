package requests

type GetBoxesRequest struct {
	Sgtins []string `json:"gtins" binding:"required,gt=0,lte=500,dive,required,sgtin"`
}
