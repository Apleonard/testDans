package models

import "gorm.io/gorm"

const (
	Description = "&description="
	Location    = "&location="
	IsFullTime  = "&full_time=true"
	Page        = "&page="
)

type UserRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type User struct {
	gorm.Model
	Username string `gorm:"size:255;not null;unique" json:"username"`
	Password string `gorm:"size:255;not null;" json:"password"`
}

type JobList []struct {
	Page        int64  `json:"-" query:"page"`
	ID          string `json:"id"`
	Type        string `json:"type"`
	URL         string `json:"url"`
	CreatedAt   string `json:"created_at"`
	Company     string `json:"company"`
	CompanyURL  string `json:"company_url"`
	Location    string `json:"location"`
	Title       string `json:"title"`
	Description string `json:"description"`
	HowToApply  string `json:"how_to_apply"`
	CompanyLogo string `json:"company_logo"`
}

type JobListParams struct {
	Page        int    `form:"page"`
	Description string `form:"description"`
	Location    string `form:"location"`
	IsFullTime  bool   `form:"IsFullTime"`
}

type JobDetailParam struct {
	ID string `form:"id"`
}

type JobDetail struct {
	Page        int64  `json:"-" query:"page"`
	ID          string `json:"id"`
	Type        string `json:"type"`
	URL         string `json:"url"`
	CreatedAt   string `json:"created_at"`
	Company     string `json:"company"`
	CompanyURL  string `json:"company_url"`
	Location    string `json:"location"`
	Title       string `json:"title"`
	Description string `json:"description"`
	HowToApply  string `json:"how_to_apply"`
	CompanyLogo string `json:"company_logo"`
}
