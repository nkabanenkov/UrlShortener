package grpc

import (
	context "context"
	"log"
	"urlshortener/internal/urlshortener/encoder"
	"urlshortener/internal/urlshortener/storage"
	"urlshortener/internal/urlshortener/validator"
	"urlshortener/pkg/urlshortener"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type UrlShortenerGrpcHandlers struct {
	UnimplementedUrlShortenerServer
	app *urlshortener.UrlShortener
}

func MakeUrlShortenerGrpcHandlers(app *urlshortener.UrlShortener) UrlShortenerGrpcHandlers {
	return UrlShortenerGrpcHandlers{
		UnimplementedUrlShortenerServer{},
		app,
	}
}

func (s *UrlShortenerGrpcHandlers) Shorten(ctx context.Context, req *ShortenRequest) (*ShortenResponse, error) {
	decodedUrl, err := s.app.Shorten(ctx, req.GetMessage())
	if err != nil {
		switch err.(type) {
		case validator.InvalidUrlError:
			return &ShortenResponse{}, status.Error(codes.InvalidArgument, err.Error())
		case storage.DatabaseError:
			log.Println("Database error occured: " + err.Error())
			return &ShortenResponse{}, status.Error(codes.Internal, err.Error())
		case encoder.EncodingOverflowError:
			log.Fatalln("Encoding overflow")
		default:
			panic(err.Error())
		}
	}

	return &ShortenResponse{Message: decodedUrl}, nil
}

func (s *UrlShortenerGrpcHandlers) Unshorten(ctx context.Context, req *UnshortenRequest) (*UnshortenResponse, error) {
	msg := req.GetMessage()
	if len(msg) == 0 {
		return &UnshortenResponse{Message: ""}, status.Error(codes.InvalidArgument, "empty URL")
	}

	decodedUrl, err := s.app.Unshorten(ctx, msg)
	if err != nil {
		switch err.(type) {
		case encoder.DecodingError:
			return &UnshortenResponse{}, status.Error(codes.InvalidArgument, err.Error())
		case storage.DatabaseError:
			log.Println("Database error occured: " + err.Error())
			return &UnshortenResponse{}, status.Error(codes.Internal, err.Error())
		case storage.UrlNotFoundError:
			return &UnshortenResponse{}, status.Error(codes.NotFound, err.Error())
		default:
			panic(err.Error())
		}
	}

	return &UnshortenResponse{Message: decodedUrl}, nil
}
