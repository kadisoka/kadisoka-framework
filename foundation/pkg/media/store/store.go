package store

import (
	"context"
	"encoding/base64"
	"io"
	"strconv"
	"strings"

	"github.com/alloyzeus/go-azfl/errors"
	mediapb "github.com/alloyzeus/go-azgrpc/azgrpc/media/v1"
	"github.com/rez-go/crock32"
	"golang.org/x/crypto/blake2b"

	"github.com/kadisoka/kadisoka-framework/foundation/pkg/app"
)

// Store contains business logic of a file service.
type Store struct {
	config        Config
	serviceClient Service
}

// New instantiates a file service.
func New(config Config, appApp app.App) (*Store, error) {
	if len(config.Modules) == 0 {
		return nil, errors.ArgMsg("config.Modules", "empty")
	}

	if config.StoreService == "" {
		return nil, errors.ArgMsg("config.StoreService", "empty")
	}

	modCfg := config.Modules[config.StoreService]
	if modCfg == nil {
		return nil, errors.ArgMsg("config.StoreService",
			config.StoreService+" not configured")
	}

	serviceClient, err := NewServiceClient(
		config.StoreService, modCfg, appApp)
	if err != nil {
		return nil, errors.ArgWrap("config.StoreService",
			config.StoreService+" initialization failed", err)
	}

	return &Store{
		config:        config,
		serviceClient: serviceClient,
	}, nil
}

// Upload uploads a file
func (mediaStore *Store) Upload(
	ctx context.Context,
	mediaName string,
	contentSource io.Reader,
	mediaType mediapb.MediaType,
) (publicURL string, err error) {
	objectURL, err := mediaStore.serviceClient.
		PutObject(ctx, mediaName, contentSource)
	if err != nil {
		return "", errors.Wrap("object putting", err)
	}

	publicURL = objectURL

	return publicURL, nil
}

// Hash length of 16 is quite short but we are combining the hash with
// the length of the data so we could reduce the chance for collisions.
const nameGenHashLength = 16

// The default key for name generation.
const nameGenKeyDefault = "N0kY"

// GenerateName is used to generate a name, for file name or other identifier,
// based on the content. It utilizes hash so the result could be used to
// prevent duplicates when storing the media object.
func (mediaStore *Store) GenerateName(stream io.Reader) string {
	var keyBytes []byte
	if mediaStore.config.NameGenerationKey != "" {
		key := strings.TrimRight(mediaStore.config.NameGenerationKey, "=")
		var err error
		keyBytes, err = base64.RawStdEncoding.DecodeString(key)
		if err != nil {
			panic(err)
		}
	}
	if len(keyBytes) == 0 {
		keyBytes = []byte(nameGenKeyDefault)
	}

	hasher, err := blake2b.New(nameGenHashLength, keyBytes)
	if err != nil {
		panic(err)
	}

	buf := make([]byte, 1*1024*1024)
	dataSize := 0
	for {
		n, err := stream.Read(buf)
		if err != nil && err != io.EOF {
			panic(err)
		}
		dataSize += n
		hasher.Write(buf)
		if err == io.EOF || n == 0 {
			break
		}
	}

	hashBytes := hasher.Sum(nil)
	encodedHash := strings.ToLower(crock32.Encode(hashBytes)) +
		"K" + strings.ToLower(crock32.Encode(keyBytes[:4])) +
		"N" + strconv.FormatInt(int64(dataSize), 16)

	return encodedHash
}

func (mediaStore Store) ImagesBaseURL() string {
	return mediaStore.config.ImagesBaseURL
}

func ContentTypeInList(contentType string, contentTypeList []string) bool {
	for _, ct := range contentTypeList {
		if ct == contentType {
			return true
		}
	}
	return false
}

func ConfigSkeleton() Config {
	return Config{
		Modules: ModuleConfigSkeletons(),
	}
}
