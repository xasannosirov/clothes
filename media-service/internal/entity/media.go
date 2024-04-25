package entity

import "time"

type Media struct {
	Id         string
	Product_Id string
	Image_Url  string
	Created_at time.Time
	Updated_at time.Time
}
