package grpc

import (
	"context"
	"errors"
	pb "github.com/atrush/pract_01.git/internal/grpc/proto"
	"github.com/atrush/pract_01.git/internal/model"
	"github.com/atrush/pract_01.git/internal/service"
	"github.com/atrush/pract_01.git/internal/shterrors"
	"github.com/google/uuid"
)

type URLsServer struct {
	pb.UnimplementedURLsServer
	svc     service.URLShortener
	baseURL string
}

func NewURLServer(svc service.URLShortener, baseURL string) *URLsServer {
	return &URLsServer{
		svc:     svc,
		baseURL: baseURL,
	}
}

func (u *URLsServer) Get(ctx context.Context, request *pb.GetRequest) (*pb.GetResponse, error) {
	var response pb.GetResponse

	url, err := u.svc.GetURL(ctx, request.ShortId)
	if err != nil {
		response.Error = err.Error()
		return &response, nil
	}

	if url == (model.ShortURL{}) {
		response.Error = ErrorURLNotFounded.Error()
		return &response, nil
	}

	if url.IsDeleted {
		response.Error = ErrorURLIsDeleted.Error()
		return &response, nil
	}

	response.SrcUrl = url.URL
	return &response, nil
}

func (u *URLsServer) GetList(ctx context.Context, request *pb.GetListRequest) (*pb.GetListResponse, error) {
	var response pb.GetListResponse
	userID, err := uuid.Parse(request.UserId)
	if err != nil {
		response.Error = ErrorWrongUserID.Error()
		return &response, nil
	}

	urlList, err := u.svc.GetUserURLList(ctx, userID)
	if err != nil {
		response.Error = err.Error()
		return &response, nil
	}

	response.List = make([]*pb.GetListItem, len(urlList))
	for i, v := range urlList {
		response.List[i] = &pb.GetListItem{
			ShortUrl: u.baseURL + "/" + v.ShortID,
			SrcUrl:   v.URL,
		}
	}

	return &response, nil
}

func (u *URLsServer) Save(ctx context.Context, request *pb.SaveRequest) (*pb.SaveResponse, error) {
	var response pb.SaveResponse

	userID, err := uuid.Parse(request.UserId)
	if err != nil {
		response.Error = ErrorWrongUserID.Error()
		return &response, nil
	}

	shortID, err := u.svc.SaveURL(ctx, request.SrcUrl, userID)
	if err != nil {
		// if url exist, return url with error
		if errors.Is(err, &shterrors.ErrorConflictSaveURL{}) {
			conflictErr, _ := err.(*shterrors.ErrorConflictSaveURL)
			response.ShortUrl = u.baseURL + "/" + conflictErr.ExistShortURL
			response.Error = ErrorURLIsExist.Error()

			return &response, nil
		}

		response.Error = err.Error()

		return &response, nil
	}

	response.ShortUrl = u.baseURL + "/" + shortID

	return &response, nil
}

func (u *URLsServer) SaveList(ctx context.Context, request *pb.SaveListRequest) (*pb.SaveListResponse, error) {
	var response pb.SaveListResponse

	userID, err := uuid.Parse(request.UserId)
	if err != nil {
		response.Error = ErrorWrongUserID.Error()
		return &response, nil
	}

	if len(request.List) == 0 {
		response.Error = ErrorURLListIsEmpty.Error()
		return &response, nil
	}

	//  make map id[url] to add
	listToAdd := make(map[string]string, len(request.List))
	for _, el := range request.List {
		listToAdd[el.CorrelationId] = el.Url
	}

	savedList, err := u.svc.SaveURLList(listToAdd, userID)
	if err != nil {
		response.Error = err.Error()
		return &response, nil
	}
	response.List = make([]*pb.SaveListItem, len(savedList))
	i := 0
	for k, v := range savedList {
		response.List[i] = &pb.SaveListItem{
			CorrelationId: k,
			Url:           u.baseURL + "/" + v,
		}
		i++
	}

	return &response, nil
}

func (u *URLsServer) DelList(ctx context.Context, request *pb.DelListRequest) (*pb.DelListResponse, error) {
	var response pb.DelListResponse
	userID, err := uuid.Parse(request.UserId)
	if err != nil {
		response.Error = ErrorWrongUserID.Error()
		return &response, nil
	}

	if len(request.List) == 0 {
		response.Error = ErrorURLListIsEmpty.Error()
		return &response, nil
	}

	if err := u.svc.DeleteURLList(userID, request.List...); err != nil {
		response.Error = err.Error()
		return &response, nil
	}

	return &pb.DelListResponse{}, nil
}
