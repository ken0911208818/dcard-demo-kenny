package model

type Pair struct {
	UserIdOne string `grom:"primaryKey" json:"userIdOne"`
	UserIdTwo string `grom:"primaryKey" json:"userIdTwo"`
}

func (p *Pair) TableName() string {
	return "pairs"
}
