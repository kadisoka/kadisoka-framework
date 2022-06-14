package iamserver

import (
	"github.com/alloyzeus/go-azfl/errors"
	iampb "github.com/alloyzeus/go-azgrpc/azgrpc/iam/v1"

	"github.com/kadisoka/kadisoka-framework/pkg/iam/pkg/iam"
)

func (core *Core) GetUserContactInformation(
	inputCtx iam.CallInputContext,
	userID iam.UserID,
) (*iampb.UserContactInfoData, error) {
	//TODO: access control
	return core.getUserContactInformationInsecure(inputCtx, userID)
}

func (core *Core) getUserContactInformationInsecure(
	inputCtx iam.CallInputContext,
	userID iam.UserID,
) (*iampb.UserContactInfoData, error) {
	userPhoneNumber, err := core.
		GetUserKeyPhoneNumber(inputCtx, userID)
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
