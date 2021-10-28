package main

import (
	"context"
	"log"
	"strings"

	"github.com/pastelnetwork/storage-challenges/application/dto"
	appgrpc "github.com/pastelnetwork/storage-challenges/application/grpc"
	"google.golang.org/api/option"
	"google.golang.org/api/transport/grpc"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/status"
)

func main() {
	clientConn, err := grpc.DialInsecure(context.Background(), option.WithEndpoint(":8080"))
	if err != nil {
		log.Fatal(err)
	}
	ctx := context.Background()
	res, err := appgrpc.NewStorageChallengeClient(clientConn).StorageChallenge(ctx, &dto.StorageChallengeRequest{})
	if err != nil {
		st, ok := status.FromError(err)
		if ok {
			log.Printf("STATUS: %d %s %d\n", st.Code(), st.Message(), len(st.Details()))
			for _, d := range st.Details() {
				switch t := d.(type) {
				case *errdetails.BadRequest:
					for _, fieldViolation := range t.FieldViolations {
						log.Printf("BAD_REQUEST----- %s = %s\n", strings.Split(fieldViolation.Field, " "), fieldViolation.Description)
					}
				default:
					log.Println(t)
				}
			}
		}

		return
	}
	log.Println(res)
}
