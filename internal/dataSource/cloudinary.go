package datasource

import (
	"betapa-antik-service/configs"
	"context"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
)

type UploadResult struct {
	URL      string
	PublicID string
}

type CloudinaryService interface {
	UploadFromMultipart(
		ctx context.Context,
		file *multipart.FileHeader,
		folder, filename string,
	) (*UploadResult, error)

	UploadFromReader(
		ctx context.Context,
		reader io.Reader,
		folder, filename string,
	) (*UploadResult, error)

	// DestroyImage(ctx context.Context, publicID string) error
	ExtractPublicIDFromURL(url string) (string, error)
	DeleteImageByURL(ctx context.Context, url string) error
}

func NewCloudinaryService(cfg *configs.CloudinaryConfig) (CloudinaryService, error) {
	cld, err := cloudinary.NewFromParams(
		cfg.CloudName,
		cfg.APIKey,
		cfg.APISecret,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to init cloudinary: %w", err)
	}

	return &cloudinaryServiceImpl{cld: cld}, nil
}

type cloudinaryServiceImpl struct {
	cld *cloudinary.Cloudinary
}

func (c *cloudinaryServiceImpl) UploadFromMultipart(
	ctx context.Context,
	file *multipart.FileHeader,
	folder, filename string,
) (*UploadResult, error) {

	f, err := file.Open()
	if err != nil {
		return nil, err
	}
	defer f.Close()

	return c.upload(ctx, f, folder, filename)
}

func (c *cloudinaryServiceImpl) UploadFromReader(
	ctx context.Context,
	reader io.Reader,
	folder, filename string,
) (*UploadResult, error) {

	return c.upload(ctx, reader, folder, filename)
}

func (c *cloudinaryServiceImpl) upload(
	ctx context.Context,
	reader io.Reader,
	folder, filename string,
) (*UploadResult, error) {

	publicID := fmt.Sprintf("%s/%s", folder, filename)

	resp, err := c.cld.Upload.Upload(ctx, reader, uploader.UploadParams{
		PublicID:     publicID,
		Folder:       folder,
		Overwrite:    boolPtr(true),
		ResourceType: "image",
	})
	if err != nil {
		return nil, err
	}

	return &UploadResult{
		URL:      resp.SecureURL,
		PublicID: resp.PublicID,
	}, nil
}

// func (c *cloudinaryServiceImpl) DestroyImage(
// 	ctx context.Context,
// 	publicID string,
// ) error {

// 	resp, err := c.cld.Upload.Destroy(ctx, uploader.DestroyParams{
// 		PublicID:   publicID,
// 		Invalidate: boolPtr(true),
// 	})
// 	if err != nil {
// 		return err
// 	}

// 	if resp.Result != "ok" {
// 		return fmt.Errorf("delete failed: %s", resp.Result)
// 	}
// 	return nil
// }

func (c *cloudinaryServiceImpl) ExtractPublicIDFromURL(url string) (string, error) {
	// contoh:
	// https://res.cloudinary.com/demo/image/upload/v1700000000/folder/name.jpg
	// => folder/name

	if url == "" {
		return "", errors.New("empty image url")
	}

	parts := strings.Split(url, "/upload/")
	if len(parts) < 2 {
		return "", fmt.Errorf("invalid cloudinary url format: %s", url)
	}

	publicIDWithExt := parts[1]

	// buang version (v123)
	publicIDWithExt = regexp.MustCompile(`^v\d+/`).ReplaceAllString(publicIDWithExt, "")

	// buang ekstensi
	publicID := strings.TrimSuffix(
		publicIDWithExt,
		filepath.Ext(publicIDWithExt),
	)

	return publicID, nil
}

func (c *cloudinaryServiceImpl) DeleteImageByURL(ctx context.Context, url string) error {
	publicID, err := c.ExtractPublicIDFromURL(url)
	if err != nil {
		return err
	}
	if publicID == "" {
		return nil
	}

	_, err = c.cld.Upload.Destroy(ctx, uploader.DestroyParams{
		PublicID: publicID,
	})
	return err
}

// func (c *cloudinaryServiceImpl) UploadImageBytes(ctx context.Context, file io.Reader, folder, filename string) (*UploadResult, error) {
// 	publicID := fmt.Sprintf("%s/%s", folder, filename)
// 	resp, err := c.cld.Upload.Upload(ctx, file, uploader.UploadParams{
// 		PublicID:     publicID,
// 		Folder:       folder,
// 		Overwrite:    boolPtr(true),
// 		ResourceType: "image",
// 	})
// 	if err != nil {
// 		return nil, err
// 	}
// 	return &UploadResult{URL: resp.SecureURL, PublicID: resp.PublicID}, nil
// }

func boolPtr(b bool) *bool { return &b }
