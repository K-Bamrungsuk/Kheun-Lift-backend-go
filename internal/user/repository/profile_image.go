package repository

import (
	"bytes"
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

// ProfileImageStore keeps one profile image per user in a Cloudflare R2 bucket (S3 API).
type ProfileImageStore struct {
	client    *s3.Client
	bucket    string
	publicURL string
}

func NewProfileImageStore(accountID, accessKeyID, secretAccessKey, bucket, publicURL string) *ProfileImageStore {
	client := s3.New(s3.Options{
		Region:       "auto",
		BaseEndpoint: aws.String("https://" + accountID + ".r2.cloudflarestorage.com"),
		Credentials:  credentials.NewStaticCredentialsProvider(accessKeyID, secretAccessKey, ""),
	})
	return &ProfileImageStore{client: client, bucket: bucket, publicURL: strings.TrimRight(publicURL, "/")}
}

// One fixed key per user: a new upload overwrites the old image, so nothing is left behind.
func profileImageKey(userID int) string {
	return "profile-images/" + strconv.Itoa(userID)
}

// Put uploads the image and returns its public URL.
func (s *ProfileImageStore) Put(ctx context.Context, userID int, contentType string, data []byte) (string, error) {
	key := profileImageKey(userID)
	_, err := s.client.PutObject(ctx, &s3.PutObjectInput{
		Bucket:      &s.bucket,
		Key:         &key,
		Body:        bytes.NewReader(data),
		ContentType: &contentType,
	})
	if err != nil {
		return "", err
	}
	// The key never changes, so the version query makes browsers load the new image.
	return fmt.Sprintf("%s/%s?v=%d", s.publicURL, key, time.Now().Unix()), nil
}

func (s *ProfileImageStore) Delete(ctx context.Context, userID int) error {
	key := profileImageKey(userID)
	_, err := s.client.DeleteObject(ctx, &s3.DeleteObjectInput{Bucket: &s.bucket, Key: &key})
	return err
}
