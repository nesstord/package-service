package requests

type GetBoxesBySgtinsRequest struct {
	Sgtins []string `json:"sgtins" binding:"required,gt=0,lte=500,dive,required,sgtin"`
}
