package media

import (
	"strings"

	"github.com/gabriel-vasile/mimetype"
	mediapb "github.com/rez-go/crux-apis/crux/media/v1"
)

func DetectType(buf []byte) string {
	// Detect always returns valid MIME.
	return mimetype.Detect(buf).String()
}

type MediaTypeInfo interface {
	// MediaType returns the type of media for this info.
	MediaType() mediapb.MediaType

	// StorePath returns a string which usually used to construct
	// path for storing the media files.
	StorePath() string

	// IsContentTypeAllowed returns true if the provided content type string
	// is allowed for the media type.
	IsContentTypeAllowed(contentType string) bool
}

var mediaTypeRegistry = map[mediapb.MediaType]MediaTypeInfo{
	mediapb.MediaType_IMAGE: &imageMediaTypeInfo{
		mediaType: mediapb.MediaType_IMAGE,
		storePath: "images"},
}

func GetMediaTypeInfoByTypeName(mediaTypeName string) MediaTypeInfo {
	if mediaTypeName == "" {
		return nil
	}
	mediaTypeName = strings.ToUpper(mediaTypeName)
	if v, ok := mediapb.MediaType_value[mediaTypeName]; ok {
		mediaType := mediapb.MediaType(v)
		return GetMediaTypeInfo(mediaType)
	}
	return nil
}

func GetMediaTypeInfo(mediaType mediapb.MediaType) MediaTypeInfo {
	return mediaTypeRegistry[mediaType]
}
