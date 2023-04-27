package service

import (
	"context"
	"harmoni/internal/entity/paginator"
	tagentity "harmoni/internal/entity/tag"
	"harmoni/internal/pkg/errorx"
	"harmoni/internal/pkg/reason"
	"harmoni/internal/usecase"

	"go.uber.org/zap"
)

type TagService struct {
	tc     *usecase.TagUseCase
	logger *zap.SugaredLogger
}

func NewTagService(
	tc *usecase.TagUseCase,
	logger *zap.SugaredLogger,
) *TagService {
	return &TagService{
		tc:     tc,
		logger: logger,
	}
}

// GetTags TODO: Add condition
func (s *TagService) GetTags(ctx context.Context, req *tagentity.GetTagsRequest) (*tagentity.GetTagsReply, error) {
	tags, err := s.tc.GetPage(ctx, req.PageSize, req.Page)
	if err != nil {
		s.logger.Errorln(err)
		return nil, err
	}

	res := paginator.Page[tagentity.TagInfo]{
		CurrentPage: tags.CurrentPage,
		PageSize:    tags.PageSize,
		Total:       tags.Total,
		Pages:       tags.Pages,
		Data:        make([]tagentity.TagInfo, 0, len(tags.Data)),
	}

	for _, tag := range tags.Data {
		res.Data = append(res.Data, tagentity.ConvertTagToDisplay(&tag))
	}

	return &tagentity.GetTagsReply{
		Page: res,
	}, nil
}

func (s *TagService) Create(ctx context.Context, req *tagentity.CreateTagRequest) (*tagentity.CreateTagReply, error) {
	_, exist, err := s.tc.GetByTagName(ctx, req.Name)
	if err != nil {
		s.logger.Errorln(err)
		return nil, err
	} else if exist {
		s.logger.Infof("Create Tag attempt failed. Tag with name '%v' already exists.\n", req.Name)
		return nil, errorx.BadRequest(reason.TagAlreadyExist)
	}

	tag := &tagentity.Tag{
		TagName:      req.Name,
		Introduction: req.Introduction,
	}
	tag, err = s.tc.Create(ctx, tag)
	if err != nil {
		s.logger.Errorln(err)
		return nil, err
	}

	return &tagentity.CreateTagReply{
		TagDetail: tagentity.ConvertTagToDetailDisplay(tag),
	}, nil
}

func (s *TagService) GetByTagID(ctx context.Context, req *tagentity.GetTagDetailRequest) (*tagentity.GetTagDetailReply, error) {
	tag, exist, err := s.tc.GetByTagID(ctx, req.TagID)
	if err != nil {
		s.logger.Errorln(err)
		return nil, err
	} else if !exist {
		s.logger.Infof("Get Tag attempt failed. Tag with ID '%v' not found.\n", req.TagID)
		return nil, errorx.NotFound(reason.TagNotFound)
	}

	return &tagentity.GetTagDetailReply{
		TagDetail: tagentity.ConvertTagToDetailDisplay(tag),
	}, nil
}
