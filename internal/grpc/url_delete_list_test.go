package grpc

import (
	"context"
	"errors"
	pb "github.com/atrush/pract_01.git/internal/grpc/proto"
	"github.com/atrush/pract_01.git/internal/service"
	mk "github.com/atrush/pract_01.git/internal/service/mock"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestURLsServer_DeleteList(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	tests := []struct {
		name        string
		svc         service.URLShortener
		request     *pb.DelListRequest
		reqResponse *pb.DelListResponse
	}{

		{
			name:        "wrong user id",
			svc:         mockNoRun(ctrl),
			request:     &pb.DelListRequest{UserId: "wrong user id", List: []string{url.ShortID, urlDeleted.ShortID}},
			reqResponse: &pb.DelListResponse{Error: ErrorWrongUserID.Error()},
		},
		{
			name:        "empty list",
			svc:         mockNoRun(ctrl),
			request:     &pb.DelListRequest{UserId: userID.String()},
			reqResponse: &pb.DelListResponse{Error: ErrorUrlListIsEmpty.Error()},
		},
		{
			name:        "delete ok",
			svc:         mockDeleteListOk(ctrl),
			request:     &pb.DelListRequest{UserId: userID.String(), List: []string{url.ShortID, urlDeleted.ShortID}},
			reqResponse: &pb.DelListResponse{},
		},
		{
			name:        "server error",
			svc:         mockDeleteListServerError(ctrl),
			request:     &pb.DelListRequest{UserId: userID.String(), List: []string{url.ShortID, urlDeleted.ShortID}},
			reqResponse: &pb.DelListResponse{Error: serverErrMessage},
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
			resp, err := client.DelList(ctx, tt.request)
			require.NoError(t, err)
			require.Equal(t, resp.Error, tt.reqResponse.Error)
		})
	}
}

//  service.URLShortener mocks
func mockDeleteListOk(ctrl *gomock.Controller) *mk.MockURLShortener {
	mock := mk.NewMockURLShortener(ctrl)
	mock.EXPECT().DeleteURLList(userID, []string{url.ShortID, urlDeleted.ShortID}).Return(nil)
	return mock
}
func mockDeleteListServerError(ctrl *gomock.Controller) *mk.MockURLShortener {
	mock := mk.NewMockURLShortener(ctrl)
	mock.EXPECT().DeleteURLList(userID, []string{url.ShortID, urlDeleted.ShortID}).Return(errors.New(serverErrMessage))
	return mock
}
