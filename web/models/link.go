package models

import "time"

type Link struct {
	Id           uint      `gorm:"primarykey" param:"link_id" query:"link_id" form:"link_id" json:"link_id"`
	OriginalUrl  string    `param:"original_url" query:"original_url" form:"original_url" json:"original_url"`
	ShortenedUrl string    `param:"shortened_url" query:"shortened_url" form:"shortened_url" json:"shortened_url"`
	UserId       uint      `param:"user_id" query:"user_id" form:"user_id" json:"user_id"`
	User         User      `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Tag          string    `gorm:"index;unique" param:"tag" query:"tag" form:"tag" json:"tag"`
	CreatedAt    time.Time `param:"created_at" query:"created_at" form:"created_at" json:"created_at"`
	ExpiresAt    time.Time `param:"expires_at" query:"expires_at" form:"expires_at" json:"expires_at"`
	ClickCount   uint      `param:"click_count" query:"click_count" form:"click_count" json:"click_count"`
}
