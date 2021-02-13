package minio

import (
	"io"

	"github.com/minio/minio-go/v6"

	"github.com/kadisoka/kadisoka-framework/foundation/pkg/app"
	"github.com/kadisoka/kadisoka-framework/foundation/pkg/errors"
	mediastore "github.com/kadisoka/kadisoka-framework/foundation/pkg/media/store"
)

type Config struct {
	Region          string `env:"REGION"`
	BucketName      string `env:"BUCKET_NAME"`
	AccessKeyID     string `env:"ACCESS_KEY_ID,required"`
	SecretAccessKey string `env:"SECRET_ACCESS_KEY,required"`
	Endpoint        string `env:"ENDPOINT,required"`
	UseSSL          bool   `env:"USE_SSL"`
}

const ServiceName = "minio"

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
	if conf.Endpoint == "" {
		return nil, errors.ArgMsg("config.Endpoint", "empty")
	}
	if conf.AccessKeyID == "" || conf.SecretAccessKey == "" {
		return nil, errors.ArgMsg("config", "access key required")
	}

	// Initialize minio client object.
	minioClient, err := minio.New(
		conf.Endpoint, conf.AccessKeyID, conf.SecretAccessKey, conf.UseSSL)
	if err != nil {
		return nil, errors.Wrap("minio client instantiation", err)
	}

	// Make a new bucket called mymusic.
	bucketName := conf.BucketName
	location := conf.Region

	err = minioClient.MakeBucket(bucketName, location)
	if err != nil {
		// Check to see if we already own this bucket (which happens if you run this twice)
		exists, errBucketExists := minioClient.BucketExists(bucketName)
		if errBucketExists != nil || !exists {
			return nil, errors.Wrap("bucket creation", err)
		}
	}

	return &Service{
		bucketName:  bucketName,
		minioClient: minioClient,
	}, nil
}

type Service struct {
	bucketName  string
	minioClient *minio.Client
}

var _ mediastore.Service = &Service{}

func (svc *Service) PutObject(
	targetKey string, contentSource io.Reader,
) (finalURL string, err error) {
	_, err = svc.minioClient.
		PutObject(svc.bucketName, targetKey, contentSource, -1, minio.PutObjectOptions{})
	if err != nil {
		return "", errors.Wrap("upload", err)
	}

	//TODO: final URL, not target key
	return targetKey, nil
}
