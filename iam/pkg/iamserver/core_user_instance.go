package iamserver

import (
	"github.com/alloyzeus/go-azfl/azfl/errors"

	"github.com/kadisoka/kadisoka-framework/iam/pkg/iam"
)

func (core *Core) contextUserOrNewInstance(
	callCtx iam.OpInputContext,
) (userRef iam.UserRefKey, newInstance bool, err error) {
	if callCtx == nil {
		return iam.UserRefKeyZero(), false, errors.ArgMsg("callCtx", "missing")
	}
	ctxAuth := callCtx.Authorization()
	if ctxAuth.IsUserContext() {
		userRef = ctxAuth.UserRef()
		if !core.UserService.IsUserRefKeyRegistered(userRef) {
			return iam.UserRefKeyZero(), false, errors.ArgMsg("callCtx.Authorization", "invalid")
		}
		return userRef, false, nil
	}

	userRef, err = core.UserService.createUserInstanceNoAC(callCtx)
	if err != nil {
		return iam.UserRefKeyZero(), false, err
	}

	return userRef, true, nil
}
