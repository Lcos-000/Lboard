package service

import (
	"context"
	"errors"
	"strings"

	"whiteboard/server/internal/model"
	"whiteboard/server/internal/repository"
	"whiteboard/server/pkg/uuid"
)

// 依赖注入，需要room和user的操作仓库
type RoomService struct {
	rooms *repository.RoomRepository
	users *repository.UserRepository
}

func NewRoomService(
	rooms *repository.RoomRepository,
	users *repository.UserRepository,
) *RoomService {
	return &RoomService{
		rooms: rooms,
		users: users,
	}
}

type CreateRoomInput struct {
	Name        string
	Description string
}

type InviteMemberInput struct {
	Email string
	Role  string
}

// CreateRoom 创建房间
func (s *RoomService) CreateRoom(ctx context.Context, ownerID string, in CreateRoomInput) (*model.Room, error) {
	name := strings.TrimSpace(in.Name)
	if name == "" {
		return nil, ErrInvalidArgument
	}

	room := &model.Room{
		ID:          uuid.NewString(),
		Name:        name,
		Description: strings.TrimSpace(in.Description),
		OwnerID:     ownerID,
	}

	// 创建房间并关联所有者
	if err := s.rooms.CreateWithOwner(ctx, room); err != nil {
		return nil, err
	}

	// 返回创建的房间信息
	return s.rooms.GetByID(ctx, room.ID)
}

// ListRooms 获取用户房间列表
func (s *RoomService) ListRooms(ctx context.Context, userID string) ([]*model.Room, error) {
	return s.rooms.ListByUser(ctx, userID)
}

// GetRoom 获取房间详情
func (s *RoomService) GetRoom(ctx context.Context, userID, roomID string) (*model.Room, error) {
	// 检查用户是否是房间成员
	if _, err := s.rooms.GetMemberRole(ctx, roomID, userID); err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, ErrForbidden
		}
		return nil, err
	}

	return s.rooms.GetByID(ctx, roomID)
}

// InviteMember 邀请房间成员
func (s *RoomService) InviteMember(ctx context.Context, operatorID, roomID string, in InviteMemberInput) error {
	role, err := s.rooms.GetMemberRole(ctx, roomID, operatorID)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return ErrForbidden
		}
		return err
	}

	if role != "owner" && role != "admin" {
		return ErrForbidden
	}

	email := strings.ToLower(strings.TrimSpace(in.Email))
	if email == "" {
		return ErrInvalidArgument
	}

	memberRole := strings.TrimSpace(in.Role)
	if memberRole == "" {
		memberRole = "editor"
	}
	// 检查角色名称是否合法
	if memberRole != "viewer" && memberRole != "editor" && memberRole != "admin" {
		return ErrInvalidArgument
	}

	// 获取用户ID，检查用户是否存在
	user, err := s.users.GetByEmail(ctx, email)
	if err != nil {
		return err
	}

	// 添加房间成员
	return s.rooms.AddMember(ctx, roomID, user.ID, memberRole)
}
