package model

import (
	"time"
)

// User 用户模型
type User struct {
	ID           string    `json:"id"`
	Username     string    `json:"username"`
	Email        string    `json:"email"`
	PasswordHash string    `json:"-"`
	CreatedAt    time.Time `json:"createdAt"`
	UpdatedAt    time.Time `json:"updatedAt"`
}

// Room 房间模型
type Room struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	OwnerID     string    `json:"ownerId"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

// RoomMember 房间成员模型
type RoomMember struct {
	RoomID    string    `json:"roomId"`
	UserID    string    `json:"userId"`
	Role      string    `json:"role"`
	CreatedAt time.Time `json:"createdAt"`
}
