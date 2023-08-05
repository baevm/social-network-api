package media

import (
	"context"
	"mime/multipart"
	"social-network-api/internal/db/models"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
)

type UploaderResult struct {
	SecureLink  string `json:"secure_link"`
	PublicLink  string `json:"public_link"`
	ExternalRef string `json:"external_ref"`
}

type Repo struct {
	client *cloudinary.Cloudinary
}

func NewRepo(cloudName, apiKey, apiSecret string) *Repo {
	cld, _ := cloudinary.NewFromParams(cloudName, apiKey, apiSecret)

	return &Repo{
		client: cld,
	}
}

func (s *Repo) Upload(ctx context.Context, file multipart.File, folder string) (*UploaderResult, error) {
	asset, err := s.client.Upload.Upload(ctx, file, uploader.UploadParams{
		Folder: folder,
	})
	if err != nil {
		return nil, models.ErrUploadFailed
	}

	return &UploaderResult{
		PublicLink:  asset.URL,
		SecureLink:  asset.SecureURL,
		ExternalRef: asset.PublicID,
	}, nil
}

func (s *Repo) Delete(ctx context.Context, publicID string) error {
	_, err := s.client.Upload.Destroy(ctx, uploader.DestroyParams{
		PublicID: publicID,
	})

	if err != nil {
		return models.ErrGetFileFailed
	}

	return nil
}
