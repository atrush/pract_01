package grpc

import (
	"context"
	pb "github.com/atrush/pract_01.git/internal/grpc/proto"
	"github.com/atrush/pract_01.git/internal/model"
	mk "github.com/atrush/pract_01.git/internal/service/mock"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/test/bufconn"
	"log"
	"net"
)

var (
	baseURL = "http://localhost:8080"

	userID           = uuid.New()
	serverErrMessage = "server error"

	url = model.ShortURL{
		ShortID:   "1xQ6p+JI",
		URL:       "https://yandex.ru/",
		UserID:    userID,
		IsDeleted: false,
		ID:        uuid.New(),
	}
	urlDeleted = model.ShortURL{
		ShortID:   "2xQ6p+JI",
		URL:       "https://yandex.ru/12",
		UserID:    userID,
		IsDeleted: true,
		ID:        uuid.New(),
	}
)

func initTestGRPCConn(ctx context.Context) (*URLsServer, *grpc.ClientConn, error) {
	urlServer := &URLsServer{baseURL: baseURL}

	conn, err := grpc.DialContext(ctx, "",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithContextDialer(dialer(urlServer)))

	return urlServer, conn, err
}

func dialer(s *URLsServer) func(context.Context, string) (net.Conn, error) {
	listener := bufconn.Listen(1024 * 1024)

	server := grpc.NewServer()

	pb.RegisterURLsServer(server, s)

	go func() {
		if err := server.Serve(listener); err != nil {
			log.Fatal(err)
		}
	}()

	return func(context.Context, string) (net.Conn, error) {
		return listener.Dial()
	}
}

//  service.URLShortener mocks
func mockNoRun(ctrl *gomock.Controller) *mk.MockURLShortener {
	mock := mk.NewMockURLShortener(ctrl)
	return mock
}
