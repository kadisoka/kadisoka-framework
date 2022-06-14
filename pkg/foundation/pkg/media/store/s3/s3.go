package s3

import (
	"context"
	"io"

	"github.com/alloyzeus/go-azfl/errors"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"

	"github.com/kadisoka/kadisoka-framework/pkg/foundation/pkg/app"
	mediastore "github.com/kadisoka/kadisoka-framework/pkg/foundation/pkg/media/store"
)

type Config struct {
	Region          string `env:"REGION,required"`
	BucketName      string `env:"BUCKET_NAME,required"`
	AccessKeyID     string `env:"ACCESS_KEY_ID"`
	SecretAccessKey string `env:"SECRET_ACCESS_KEY"`
}

const ServiceName = "s3"

func init() {
	mediastore.RegisterModule(
		ServiceName,
		mediastore.Module{
			NewService: NewService,
			ServiceConfigSkeleton: func() mediastore.ServiceConfig {
				cfg := ConfigSkeleton()
				return &cfg
			},
		})
}

func ConfigSkeleton() Config { return Config{} }

func NewService(
	config mediastore.ServiceConfig,
	_ app.App,
) (mediastore.Service, error) {
	if config == nil {
		return nil, errors.ArgMsg("config", "missing")
	}

	conf, ok := config.(*Config)
	if !ok {
		return nil, errors.ArgMsg("config", "type invalid")
	}
	if conf == nil || conf.Region == "" || conf.BucketName == "" {
		return nil, errors.ArgMsg("config", "fields invalid")
	}

	var creds *credentials.Credentials
	if conf.AccessKeyID != "" {
		creds = credentials.NewStaticCredentials(
			conf.AccessKeyID,
			conf.SecretAccessKey,
			"",
		)
	}
	sess, err := session.NewSession(&aws.Config{
		Region:      aws.String(conf.Region),
		Credentials: creds,
	})
	if err != nil {
		return nil, errors.Wrap("AWS session", err)
	}

	const uploadPartSize = 10 * 1024 * 1024 // 10MiB

	return &Service{
		bucketName: conf.BucketName,
		uploader: s3manager.NewUploader(sess, func(u *s3manager.Uploader) {
			u.PartSize = uploadPartSize
		}),
	}, nil
}

type Service struct {
	bucketName string
	uploader   *s3manager.Uploader
}

var _ mediastore.Service = &Service{}

func (svc *Service) PutObject(
	ctx context.Context,
	targetKey string,
	contentSource io.Reader,
) (finalURL string, err error) {
	result, err := svc.uploader.
		Upload(&s3manager.UploadInput{
			Bucket: aws.String(svc.bucketName),
			Key:    aws.String(targetKey),
			Body:   contentSource,
		})
	if err != nil {
		return "", errors.Wrap("upload", err)
	}

	return result.Location, nil
}
