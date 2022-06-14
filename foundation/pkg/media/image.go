package media

import (
	mediapb "github.com/alloyzeus/go-azgrpc/azgrpc/media/v1"
)

//TODO:
// - differentiate between formats with animation support
// - mapping of input formats with output formats. e.g., we will allow a TIFF
//   but we will serve it as PNG.
// - differentiate between 'allowed' and 'supported'.
var imageAllowedContentTypes = []string{"image/jpg", "image/jpeg", "image/png"}

type imageMediaTypeInfo struct {
	mediaType mediapb.MediaType
	storePath string
}

var _ MediaTypeInfo = &imageMediaTypeInfo{}

func (typeInfo *imageMediaTypeInfo) MediaType() mediapb.MediaType {
	if typeInfo.mediaType == mediapb.MediaType_MEDIA_TYPE_UNSPECIFIED {
		return mediapb.MediaType_IMAGE
	}
	return typeInfo.mediaType
}

func (typeInfo *imageMediaTypeInfo) StorePath() string {
	return typeInfo.storePath
}

func (typeInfo *imageMediaTypeInfo) IsContentTypeAllowed(contentType string) bool {
	for _, ct := range imageAllowedContentTypes {
		if ct == contentType {
			return true
		}
	}
	return false
}
