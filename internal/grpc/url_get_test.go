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

func TestURLsServer_Get(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	tests := []struct {
		name        string
		svc         service.URLShortener
		request     *pb.GetRequest
		reqErr      bool
		reqResponse *pb.GetResponse
	}{
		{
			name:        "exist",
			svc:         mockGetExistURL(ctrl),
			request:     &pb.GetRequest{ShortId: url.ShortID},
			reqResponse: &pb.GetResponse{SrcUrl: url.URL},
		},
		{
			name:        "not exist",
			svc:         mockGetNotExistURL(ctrl),
			request:     &pb.GetRequest{ShortId: "8xQ6p+JI"},
			reqResponse: &pb.GetResponse{Error: ErrorURLNotFounded.Error()},
		},
		{
			name:        "is deleted",
			svc:         mockGetDeletedURL(ctrl),
			request:     &pb.GetRequest{ShortId: urlDeleted.ShortID},
			reqResponse: &pb.GetResponse{Error: ErrorURLIsDeleted.Error()},
		},
		{
			name:        "server error",
			svc:         mockGetServerError(ctrl),
			request:     &pb.GetRequest{ShortId: urlDeleted.ShortID},
			reqResponse: &pb.GetResponse{Error: serverErrMessage},
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
			resp, err := client.Get(ctx, tt.request)
			require.NoError(t, err)

			require.Equal(t, resp.SrcUrl, tt.reqResponse.SrcUrl)
			require.Equal(t, resp.Error, tt.reqResponse.Error)
		})
	}
}

//  service.URLShortener mocks
func mockGetExistURL(ctrl *gomock.Controller) *mk.MockURLShortener {
	mock := mk.NewMockURLShortener(ctrl)
	mock.EXPECT().GetURL(gomock.Any(), url.ShortID).Return(url, nil)
	return mock
}
func mockGetDeletedURL(ctrl *gomock.Controller) *mk.MockURLShortener {
	mock := mk.NewMockURLShortener(ctrl)
	mock.EXPECT().GetURL(gomock.Any(), urlDeleted.ShortID).Return(urlDeleted, nil)
	return mock
}
func mockGetNotExistURL(ctrl *gomock.Controller) *mk.MockURLShortener {
	mock := mk.NewMockURLShortener(ctrl)
	mock.EXPECT().GetURL(gomock.Any(), gomock.Any()).Return(model.ShortURL{}, nil)
	return mock
}
func mockGetServerError(ctrl *gomock.Controller) *mk.MockURLShortener {
	mock := mk.NewMockURLShortener(ctrl)
	mock.EXPECT().GetURL(gomock.Any(), gomock.Any()).Return(model.ShortURL{}, errors.New(serverErrMessage))
	return mock
}
