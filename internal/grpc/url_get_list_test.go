package grpc

import (
	"context"
	"errors"
	pb "github.com/atrush/pract_01.git/internal/grpc/proto"
	"github.com/atrush/pract_01.git/internal/model"
	"github.com/atrush/pract_01.git/internal/service"
	mk "github.com/atrush/pract_01.git/internal/service/mock"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestURLsServer_GetList(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	tests := []struct {
		name        string
		svc         service.URLShortener
		request     *pb.GetListRequest
		reqErr      bool
		reqResponse *pb.GetListResponse
	}{
		{
			name:    "exist",
			svc:     mockGetListExistURL(ctrl),
			request: &pb.GetListRequest{UserId: userID.String()},
			reqResponse: &pb.GetListResponse{
				List: []*pb.GetListItem{
					{SrcUrl: url.URL, ShortUrl: baseURL + "/" + url.ShortID},
					{SrcUrl: urlDeleted.URL, ShortUrl: baseURL + "/" + urlDeleted.ShortID},
				},
			},
		},
		{
			name:        "wrong user id",
			svc:         mockGetListURLNoRun(ctrl),
			request:     &pb.GetListRequest{UserId: "wrong user id"},
			reqResponse: &pb.GetListResponse{Error: ErrorWrongUserID.Error()},
		},
		{
			name:        "not found",
			svc:         mockGetListEmpty(ctrl),
			request:     &pb.GetListRequest{UserId: userID.String()},
			reqResponse: &pb.GetListResponse{},
		},
		{
			name:        "server error",
			svc:         mockGetListURLServerError(ctrl),
			request:     &pb.GetListRequest{UserId: userID.String()},
			reqResponse: &pb.GetListResponse{Error: serverErrMessage},
		},
	}

	ctx := context.Background()

	urlServer, conn, err := initTestGRPCConn(ctx)
	require.NoError(t, err)
	defer conn.Close()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// set service mock
			urlServer.svc = tt.svc

			client := pb.NewURLsClient(conn)
			resp, err := client.GetList(ctx, tt.request)
			require.NoError(t, err)

			require.Equal(t, resp.Error, tt.reqResponse.Error)

			for _, rs := range tt.reqResponse.List {
				isFounded := false
				for _, v := range resp.List {
					if v.SrcUrl == rs.SrcUrl && v.ShortUrl == rs.ShortUrl {
						isFounded = true
						break
					}
				}
				require.True(t, isFounded, "url: %v ; short url: %v not founded", rs.SrcUrl, rs.ShortUrl)
			}
		})
	}
}

//  service.URLShortener mocks
func mockGetListEmpty(ctrl *gomock.Controller) *mk.MockURLShortener {
	mock := mk.NewMockURLShortener(ctrl)
	mock.EXPECT().GetUserURLList(gomock.Any(), userID).Return(nil, nil)
	return mock
}
func mockGetListExistURL(ctrl *gomock.Controller) *mk.MockURLShortener {
	mock := mk.NewMockURLShortener(ctrl)
	mock.EXPECT().GetUserURLList(gomock.Any(), userID).Return([]model.ShortURL{url, urlDeleted}, nil)
	return mock
}
func mockGetListURLNoRun(ctrl *gomock.Controller) *mk.MockURLShortener {
	mock := mk.NewMockURLShortener(ctrl)
	return mock
}
func mockGetListURLServerError(ctrl *gomock.Controller) *mk.MockURLShortener {
	mock := mk.NewMockURLShortener(ctrl)
	mock.EXPECT().GetUserURLList(gomock.Any(), userID).Return(nil, errors.New(serverErrMessage))
	return mock
}
