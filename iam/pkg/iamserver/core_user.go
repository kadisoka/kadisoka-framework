package iamserver

import (
	"github.com/alloyzeus/go-azfl/azfl/errors"
	iampb "github.com/rez-go/crux-apis/crux/iam/v1"

	"github.com/kadisoka/kadisoka-framework/iam/pkg/iam"
)

func (core *Core) GetUserContactInformation(
	callCtx iam.OpInputContext,
	userRef iam.UserRefKey,
) (*iampb.UserContactInfoData, error) {
	//TODO: access control
	return core.getUserContactInformationNoAC(callCtx, userRef)
}

func (core *Core) getUserContactInformationNoAC(
	callCtx iam.OpInputContext,
	userRef iam.UserRefKey,
) (*iampb.UserContactInfoData, error) {
	userPhoneNumber, err := core.
		GetUserKeyPhoneNumber(callCtx, userRef)
	if err != nil {
		return nil, errors.Wrap("get user key phone number", err)
	}
	if userPhoneNumber == nil {
		return nil, nil
	}
	return &iampb.UserContactInfoData{
		PhoneNumber: userPhoneNumber.String(),
	}, nil
}
