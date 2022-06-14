package iamserver

import (
	"io"
	"path"
	"strings"

	"github.com/alloyzeus/go-azfl/errors"
	mediapb "github.com/alloyzeus/go-azgrpc/azgrpc/media/v1"
	"github.com/jmoiron/sqlx"

	"github.com/kadisoka/kadisoka-framework/pkg/foundation/pkg/media"
	"github.com/kadisoka/kadisoka-framework/pkg/iam/pkg/iam"
)

type ProfileImageFile interface {
	io.Reader
	io.Seeker
}

func (core *Core) SetUserProfileImageURL(
	inputCtx iam.CallInputContext,
	userID iam.UserID,
	profileImageURL string,
) error {
	ctxAuth := inputCtx.Authorization()
	// Change this if we want to allow service client to update a user's profile
	// (we'll need a better access control for service clients)
	if !ctxAuth.IsUserSubject() {
		return iam.ErrUserContextRequired
	}
	// Don't allow changing other user's for now
	if !ctxAuth.IsUser(userID) {
		return iam.ErrOperationNotAllowed
	}
	if profileImageURL != "" && !core.isUserProfileImageURLAllowed(profileImageURL) {
		return errors.ArgMsg("profileImageURL", "unsupported")
	}

	//TODO: on changes, update caches, emit events only if there's any changes

	return doTx(core.db, func(dbTx *sqlx.Tx) error {
		_, txErr := dbTx.Exec(
			`UPDATE `+userProfileImageKeyDBTableName+` `+
				"SET _md_ts = $1, _md_uid = $2, _md_tid = $3 "+
				"WHERE user_id = $2 AND _md_ts IS NULL",
			inputCtx.CallInputMetadata().ReceiveTime,
			ctxAuth.UserIDNum().PrimitiveValue(),
			ctxAuth.TerminalIDNum().PrimitiveValue())
		if txErr != nil {
			return errors.Wrap("mark current profile image URL as deleted", txErr)
		}
		if profileImageURL != "" {
			_, txErr = dbTx.Exec(
				`INSERT INTO `+userProfileImageKeyDBTableName+` `+
					"(user_id, profile_image_key, _mc_uid, _mc_tid) VALUES "+
					"($1, $2, $3, $4)",
				ctxAuth.UserIDNum().PrimitiveValue(), profileImageURL,
				ctxAuth.UserIDNum().PrimitiveValue(), ctxAuth.TerminalIDNum().PrimitiveValue())
			if txErr != nil {
				return errors.Wrap("insert new profile image URL", txErr)
			}
		}
		return nil
	})
}

func (core *Core) SetUserProfileImageByFile(
	inputCtx iam.CallInputContext,
	userID iam.UserID,
	imageFile ProfileImageFile,
) (imageURL string, err error) {
	//TODO: configurable
	const bucketSubPath = "user_profile_images/"
	const mediaType = mediapb.MediaType_IMAGE

	mediaTypeInfo := media.GetMediaTypeInfo(mediaType)
	if mediaTypeInfo == nil {
		return "", errors.Msg("media type info unavailable") //.fields({mediaType: MediaType_IMAGE})
	}

	detectionBytes := make([]byte, 512)
	_, err = imageFile.Read(detectionBytes)
	if err != nil {
		return "", errors.Wrap("content type detection", err)
	}
	imageFile.Seek(0, io.SeekStart)

	contentType := media.DetectType(detectionBytes)
	if !mediaTypeInfo.IsContentTypeAllowed(contentType) {
		return "", errors.ArgMsg("imageFile", "media type not allowed")
	}

	filename := core.mediaStore.GenerateName(imageFile)
	imageFile.Seek(0, io.SeekStart)

	imageKey, err := core.mediaStore.
		Upload(
			inputCtx,
			path.Join(bucketSubPath, filename),
			imageFile,
			mediaType)
	if err != nil {
		return "", errors.Wrap("file store", err)
	}

	err = core.SetUserProfileImageURL(inputCtx, userID, imageKey)
	if err != nil {
		return "", errors.Wrap("user profile image URL update", err)
	}

	return core.BuildUserProfileImageURL(imageKey), nil
}

func (core *Core) BuildUserProfileImageURL(imageKey string) string {
	if !strings.HasPrefix(imageKey, "https://") && !strings.HasPrefix(imageKey, "http://") {
		var imagesBaseURL string
		if core.mediaStore != nil {
			imagesBaseURL = core.mediaStore.ImagesBaseURL()
		}
		if imagesBaseURL != "" {
			imagesBaseURL = strings.TrimSuffix(imagesBaseURL, "/")
			if strings.HasPrefix(imageKey, "/") {
				return imagesBaseURL + imageKey
			} else {
				return imagesBaseURL + "/" + imageKey
			}
		}
	}
	return imageKey
}
