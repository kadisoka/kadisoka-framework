package iamserver

import (
	"github.com/alloyzeus/go-azfl/errors"

	"github.com/kadisoka/kadisoka-framework/iam/pkg/iam"
)

func (core *Core) contextUserOrNewInstance(
	inputCtx iam.CallInputContext,
) (userID iam.UserID, newInstance bool, err error) {
	if inputCtx == nil {
		return iam.UserIDZero(), false, errors.ArgMsg("inputCtx", "missing")
	}
	ctxAuth := inputCtx.Authorization()
	if ctxAuth.IsUserSubject() {
		userID = ctxAuth.UserID()
		if !core.UserService.IsUserIDRegistered(userID) {
			return iam.UserIDZero(), false, errors.ArgMsg("inputCtx.Authorization", "invalid")
		}
		return userID, false, nil
	}

	userID, err = core.UserService.createUserInstanceInsecure(inputCtx)
	if err != nil {
		return iam.UserIDZero(), false, err
	}

	return userID, true, nil
}
