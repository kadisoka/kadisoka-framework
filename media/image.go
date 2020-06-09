package media

import (
	"github.com/rez-go/crux-apis/crux/media/v1"
)

var imageAllowedContentTypes = []string{"image/jpg", "image/jpeg", "image/png"}

type imageMediaTypeInfo struct {
	mediaType     mediapb.MediaType
	directoryName string
}

func (typeInfo *imageMediaTypeInfo) MediaType() mediapb.MediaType {
	if typeInfo.mediaType == mediapb.MediaType_MEDIA_TYPE_UNSPECIFIED {
		return mediapb.MediaType_IMAGE
	}
	return typeInfo.mediaType
}
func (typeInfo *imageMediaTypeInfo) DirectoryName() string {
	if typeInfo.directoryName == "" {
		panic("directory name is unspecified")
	}
	return typeInfo.directoryName
}

func (typeInfo *imageMediaTypeInfo) IsContentTypeAllowed(contentType string) bool {
	for _, ct := range imageAllowedContentTypes {
		if ct == contentType {
			return true
		}
	}
	return false
}
