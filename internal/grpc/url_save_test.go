package grpc

import (
	"context"
	"errors"
	pb "github.com/atrush/pract_01.git/internal/grpc/proto"
	"github.com/atrush/pract_01.git/internal/service"
	mk "github.com/atrush/pract_01.git/internal/service/mock"
	"github.com/atrush/pract_01.git/internal/shterrors"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestURLsServer_Save(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	tests := []struct {
		name        string
		svc         service.URLShortener
		request     *pb.SaveRequest
		reqErr      bool
		reqResponse *pb.SaveResponse
	}{

		{
			name:        "wrong user id",
			svc:         mockNoRun(ctrl),
			request:     &pb.SaveRequest{UserId: "wrong user id"},
			reqResponse: &pb.SaveResponse{Error: ErrorWrongUserID.Error()},
		},
		{
			name:        "save ok",
			svc:         mockSaveOk(ctrl),
			request:     &pb.SaveRequest{UserId: url.UserID.String(), SrcUrl: url.URL},
			reqResponse: &pb.SaveResponse{ShortUrl: baseURL + "/" + url.ShortID},
		},
		{
			name:        "save exist",
			svc:         mockSaveExist(ctrl),
			request:     &pb.SaveRequest{UserId: url.UserID.String(), SrcUrl: url.URL},
			reqResponse: &pb.SaveResponse{ShortUrl: baseURL + "/" + url.ShortID, Error: ErrorURLIsExist.Error()},
		},

		{
			name:        "server error",
			svc:         mockSaveServerError(ctrl),
			request:     &pb.SaveRequest{UserId: url.UserID.String(), SrcUrl: url.URL},
			reqResponse: &pb.SaveResponse{Error: serverErrMessage},
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
			resp, err := client.Save(ctx, tt.request)
			require.NoError(t, err)

			require.Equal(t, resp.ShortUrl, tt.reqResponse.ShortUrl)
			require.Equal(t, resp.Error, tt.reqResponse.Error)
		})
	}
}

//  service.URLShortener mocks
func mockSaveOk(ctrl *gomock.Controller) *mk.MockURLShortener {
	mock := mk.NewMockURLShortener(ctrl)
	mock.EXPECT().SaveURL(gomock.Any(), url.URL, userID).Return(url.ShortID, nil)
	return mock
}
func mockSaveServerError(ctrl *gomock.Controller) *mk.MockURLShortener {
	mock := mk.NewMockURLShortener(ctrl)
	mock.EXPECT().SaveURL(gomock.Any(), url.URL, userID).Return("", errors.New(serverErrMessage))
	return mock
}
func mockSaveExist(ctrl *gomock.Controller) *mk.MockURLShortener {
	mock := mk.NewMockURLShortener(ctrl)
	var errExist error = &shterrors.ErrorConflictSaveURL{ExistShortURL: url.ShortID}
	mock.EXPECT().SaveURL(gomock.Any(), url.URL, userID).Return("", errExist)
	return mock
}
