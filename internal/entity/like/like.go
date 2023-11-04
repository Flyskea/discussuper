package like

import (
	"context"
	"encoding/json"
	"harmoni/internal/entity"
	"harmoni/internal/entity/paginator"
	"harmoni/internal/pkg/common"
	"harmoni/internal/types/events/like"
)

type LikeType uint8

const (
	LikePost LikeType = iota + 1
	LikeComment
	LikeUser
)

func (t *LikeType) ToEventLikeType() like.LikeType {
	switch *t {
	case LikePost:
		return like.LikePost
	case LikeComment:
		return like.LikeComment
	case LikeUser:
		return like.LikeUser
	}
	return like.LikePost
}

type Like struct {
	ID uint `gorm:"primarykey;type:BIGINT UNSIGNED not NULL AUTO_INCREMENT;"`
	entity.TimeMixin
	UserID       int64    `gorm:"not null"`
	TargetUserID int64    `gorm:"not null"`
	LikingID     int64    `gorm:"not null"`
	LikeType     LikeType `gorm:"not null;type:TINYINT UNSIGNED"`
	Canceled     bool     `gorm:"not null;default:0;"`
}

func (*Like) TableName() string {
	return "like"
}

type LikeCacheInfo struct {
	LikingID  int64
	UpdatedAt int64 `gorm:"serializer:unixtime"`
}

func (r *LikeCacheInfo) ToJSONString() string {
	codeBytes, _ := json.Marshal(r)
	return common.BytesToString(codeBytes)
}

func (r *LikeCacheInfo) FromJSONString(data string) error {
	return json.Unmarshal(common.StringToBytes(data), r)
}

var (
	LikeTypeList = []LikeType{LikeUser, LikePost, LikeComment}
)

type LikeRepository interface {
	Like(ctx context.Context, like *Like, targetUserID int64, isCancel bool) error
	Save(ctx context.Context, like *Like, isCancel bool) error
	LikeCount(ctx context.Context, like *Like) (int64, bool, error)
	BatchLikeCount(ctx context.Context, likeType LikeType) (map[int64]int64, error)
	BatchLikeCountByIDs(ctx context.Context, likingIDs []int64, likeType LikeType) (map[int64]int64, error)
	// UpdateLikeCount(ctx context.Context, like *Like, count int8) error
	ListLikingIDs(ctx context.Context, query *LikeQuery) (paginator.Page[int64], error)
	IsLiking(ctx context.Context, like *Like) (bool, error)
	CacheLikeCount(ctx context.Context, like *Like, count int64) error
}
