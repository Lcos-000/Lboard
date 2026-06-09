package repository

import (
	"context"
	"errors"
	"whiteboard/server/internal/model"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// RoomRepository 房间仓库
type RoomRepository struct {
	db *pgxpool.Pool
}

// NewRoomRepository 创建房间仓库
func NewRoomRepository(db *pgxpool.Pool) *RoomRepository {
	return &RoomRepository{db: db}
}

// CreateWithOwner 创建房间并添加所有者为房间成员
func (r *RoomRepository) CreateWithOwner(ctx context.Context, room *model.Room) error {
	// 开始事务
	tx, err := r.db.Begin(ctx)
	// 检查事务是否成功
	if err != nil {
		return err
	}
	// 提交事务
	defer tx.Rollback(ctx)
	_, err = tx.Exec(ctx, `
INSERT INTO rooms (id, name, description, owner_id)
VALUES ($1, $2, $3, $4)
`, room.ID, room.Name, room.Description, room.OwnerID)
	if err != nil {
		return err
	}
	_, err = tx.Exec(ctx, `
INSERT INTO room_members (room_id, user_id, role)
VALUES ($1, $2, $3)
`, room.ID, room.OwnerID, "owner")
	if err != nil {
		return err
	}
	return tx.Commit(ctx)
}

// ListByUser 获取用户创建的房间列表
func (r *RoomRepository) ListByUser(ctx context.Context, userID string) ([]*model.Room, error) {
	rows, err := r.db.Query(ctx, `
SELECT r.id, r.name, r.description, r.owner_id, r.created_at, r.updated_at
FROM rooms r
JOIN room_members rm ON rm.room_id = r.id
WHERE rm.user_id = $1
ORDER BY r.created_at DESC
`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var rooms []*model.Room
	for rows.Next() {
		var room model.Room
		if err := rows.Scan(
			&room.ID,
			&room.Name,
			&room.Description,
			&room.OwnerID,
			&room.CreatedAt,
			&room.UpdatedAt,
		); err != nil {
			return nil, err
		}
		rooms = append(rooms, &room)
	}
	return rooms, rows.Err()
}

// GetByID 根据ID获取房间
func (r *RoomRepository) GetByID(ctx context.Context, roomID string) (*model.Room, error) {
	row := r.db.QueryRow(ctx, `
SELECT id, name, description, owner_id, created_at, updated_at
FROM rooms
WHERE id = $1
`, roomID)
	var room model.Room
	if err := row.Scan(
		&room.ID,
		&room.Name,
		&room.Description,
		&room.OwnerID,
		&room.CreatedAt,
		&room.UpdatedAt,
	); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return &room, nil
}

// GetMemberRole 获取房间成员角色
func (r *RoomRepository) GetMemberRole(ctx context.Context, roomID, userID string) (string, error) {
	row := r.db.QueryRow(ctx, `
SELECT role
FROM room_members
WHERE room_id = $1 AND user_id = $2
`, roomID, userID)
	var role string
	if err := row.Scan(&role); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return "", ErrNotFound
		}
		return "", err
	}
	return role, nil
}

// AddMember 添加房间成员
func (r *RoomRepository) AddMember(ctx context.Context, roomID, userID, role string) error {
	_, err := r.db.Exec(ctx, `
INSERT INTO room_members (room_id, user_id, role)
VALUES ($1, $2, $3)
ON CONFLICT (room_id, user_id)
DO UPDATE SET role = EXCLUDED.role
`, roomID, userID, role)
	return err
}
