package proxy

import (
	"context"

	filespb "github.com/Game-WW/gateway-be/gen/files"
	userspb "github.com/Game-WW/gateway-be/gen/users"
	"google.golang.org/grpc"
)

type usersClient struct {
	userspb.UsersServiceClient
	media filespb.MediaServiceClient
}

func newUsersClient(
	users userspb.UsersServiceClient,
	media filespb.MediaServiceClient,
) *usersClient {
	return &usersClient{
		UsersServiceClient: users,
		media:              media,
	}
}

func (c *usersClient) UpdateAvatar(
	ctx context.Context,
	req *userspb.UpdateAvatarRequest,
	opts ...grpc.CallOption,
) (*userspb.UpdateAvatarResponse, error) {
	file, err := c.media.UploadFile(ctx, &filespb.UploadFileRequest{
		Image: req.GetImage(),
	}, opts...)
	if err != nil {
		return nil, err
	}

	_, err = c.UsersServiceClient.SetAvatarFileID(
		ctx,
		&userspb.SetAvatarFileIDRequest{
			FileId: file.GetFileId(),
		},
		opts...,
	)
	if err != nil {
		return nil, err
	}
	return &userspb.UpdateAvatarResponse{}, nil
}
