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

func TestURLsServer_SaveList(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	tests := []struct {
		name        string
		svc         service.URLShortener
		request     *pb.SaveListRequest
		reqResponse *pb.SaveListResponse
	}{

		{
			name: "wrong user id",
			svc:  mockNoRun(ctrl),
			request: &pb.SaveListRequest{UserId: "wrong user id", List: []*pb.SaveListItem{
				{CorrelationId: "01", Url: url.URL},
				{CorrelationId: "02", Url: urlDeleted.URL},
			}},
			reqResponse: &pb.SaveListResponse{Error: ErrorWrongUserID.Error()},
		},
		{
			name:        "empty list",
			svc:         mockNoRun(ctrl),
			request:     &pb.SaveListRequest{UserId: userID.String()},
			reqResponse: &pb.SaveListResponse{Error: ErrorUrlListIsEmpty.Error()},
		},
		{
			name: "save ok",
			svc:  mockSaveListOk(ctrl),
			request: &pb.SaveListRequest{UserId: userID.String(), List: []*pb.SaveListItem{
				{CorrelationId: "01", Url: url.URL},
				{CorrelationId: "02", Url: urlDeleted.URL},
			}},
			reqResponse: &pb.SaveListResponse{List: []*pb.SaveListItem{
				{CorrelationId: "01", Url: baseURL + "/" + url.ShortID},
				{CorrelationId: "02", Url: baseURL + "/" + urlDeleted.ShortID},
			}},
		},
		{
			name: "server error",
			svc:  mockSaveListServerError(ctrl),
			request: &pb.SaveListRequest{UserId: userID.String(), List: []*pb.SaveListItem{
				{CorrelationId: "01", Url: url.URL},
				{CorrelationId: "02", Url: urlDeleted.URL},
			}},
			reqResponse: &pb.SaveListResponse{Error: serverErrMessage},
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
			resp, err := client.SaveList(ctx, tt.request)
			require.NoError(t, err)
			require.Equal(t, resp.Error, tt.reqResponse.Error)

			for _, rs := range tt.reqResponse.List {
				isFounded := false
				for _, v := range resp.List {
					if v.CorrelationId == rs.CorrelationId && v.Url == rs.Url {
						isFounded = true
						break
					}
				}
				require.True(t, isFounded, "url: %v ; short url: %v not founded", rs.CorrelationId, rs.Url)
			}
		})
	}
}

//  service.URLShortener mocks
func mockSaveListOk(ctrl *gomock.Controller) *mk.MockURLShortener {
	mock := mk.NewMockURLShortener(ctrl)
	mock.EXPECT().SaveURLList(map[string]string{"01": url.URL, "02": urlDeleted.URL}, userID).Return(
		map[string]string{"01": url.ShortID, "02": urlDeleted.ShortID}, nil)
	return mock
}
func mockSaveListServerError(ctrl *gomock.Controller) *mk.MockURLShortener {
	mock := mk.NewMockURLShortener(ctrl)
	mock.EXPECT().SaveURLList(gomock.Any(), userID).Return(nil, errors.New(serverErrMessage))
	return mock
}
