package requests

type AggregateRequest struct {
	Sscc    string   `json:"sscc" binding:"required,sscc,len=18"`
	Created string   `json:"created" binding:"required,datetime=2006-01-02T15:04:05+07:00"`
	Sgtins  []string `json:"sgtins" binding:"required,gt=0,dive,required,sgtin,len=27"`
}
