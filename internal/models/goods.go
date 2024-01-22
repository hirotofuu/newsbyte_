package models

// good用タイプ

type Good struct {
	IsGoodFlag bool `json:"is_good_flag"`
	GoodNum    int  `json:"good_num"`
}
