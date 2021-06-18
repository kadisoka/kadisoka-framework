package grpc

import (
	"context"
	"time"

	"github.com/alloyzeus/go-azfl/azfl/errors"
	pbtypes "github.com/gogo/protobuf/types"
	iampb "github.com/rez-go/crux-apis/crux/iam/v1"
	"google.golang.org/grpc"
	grpccodes "google.golang.org/grpc/codes"
	grpcstatus "google.golang.org/grpc/status"

	grpcerrs "github.com/kadisoka/kadisoka-framework/foundation/pkg/api/grpc/errors"
	"github.com/kadisoka/kadisoka-framework/iam/pkg/iam"
	"github.com/kadisoka/kadisoka-framework/iam/pkg/iamserver"
	"github.com/kadisoka/kadisoka-framework/volib/pkg/telephony"
)

type TerminalAuthorizationServiceServer struct {
	iamServerCore *iamserver.Core
}

func NewTerminalAuthorizationServiceServer(
	iamServerCore *iamserver.Core,
	grpcServer *grpc.Server,
) *TerminalAuthorizationServiceServer {
	authServer := &TerminalAuthorizationServiceServer{
		iamServerCore,
	}
	iampb.RegisterTerminalAuthorizationServiceServer(grpcServer, authServer)
	return authServer
}

//TODO: verification methods
func (authServer *TerminalAuthorizationServiceServer) InitiateUserTerminalAuthorizationByPhoneNumber(
	callCtx context.Context,
	reqProto *iampb.InitiateUserTerminalAuthorizationByPhoneNumberRequest,
) (*iampb.InitiateUserTerminalAuthorizationByPhoneNumberResponse, error) {
	reqCtx, err := authServer.iamServerCore.GRPCOpInputContext(callCtx)
	if err != nil {
		panic(err) //TODO: translate and return the error
	}
	ctxAuth := reqCtx.Authorization()
	if ctxAuth.IsValid() {
		logCtx(reqCtx).
			Warn().Msgf("Authorization context must not be valid: %#v", reqCtx)
		return nil, grpcstatus.Error(grpccodes.Unauthenticated, "")
	}

	appRef, err := iam.ApplicationRefKeyFromAZIDText(reqProto.ClientCredentials.ClientId)
	if err != nil {
		panic(err)
	}

	phoneNumber, err := telephony.PhoneNumberFromString(reqProto.PhoneNumber)
	if err != nil {
		logCtx(reqCtx).
			Warn().Err(err).Str("phone_number", reqProto.PhoneNumber).
			Msg("Phone number format")
		return nil, grpcstatus.Error(grpccodes.InvalidArgument, "")
	}

	authStartOutput := authServer.iamServerCore.
		StartTerminalAuthorizationByPhoneNumber(
			iamserver.TerminalAuthorizationByPhoneNumberStartInput{
				Context:        reqCtx,
				ApplicationRef: appRef,
				Data: iamserver.TerminalAuthorizationByPhoneNumberStartInputData{
					PhoneNumber:         phoneNumber,
					VerificationMethods: nil,
					TerminalAuthorizationStartInputBaseData: iamserver.TerminalAuthorizationStartInputBaseData{
						DisplayName: reqProto.TerminalInfo.DisplayName,
					},
				},
			})
	if err = authStartOutput.Context.Err; err != nil {
		switch err.(type) {
		case errors.CallError:
			logCtx(reqCtx).
				Warn().Err(err).
				Msgf("StartTerminalAuthorizationByPhoneNumber with %v failed",
					phoneNumber)
			return nil, grpcstatus.Error(grpccodes.InvalidArgument, "")
		}
		logCtx(reqCtx).
			Error().Err(err).
			Msgf("StartTerminalAuthorizationByPhoneNumber with %v failed",
				phoneNumber)
		return nil, grpcerrs.Error(err)
	}

	var codeExpiryProto *pbtypes.Timestamp
	if codeExpiry := authStartOutput.Data.VerificationCodeExpiryTime; codeExpiry != nil {
		codeExpiryProto, err = pbtypes.TimestampProto(*codeExpiry)
		if err != nil {
			panic(err)
		}
	}
	resp := iampb.InitiateUserTerminalAuthorizationByPhoneNumberResponse{
		TerminalId:                 authStartOutput.Data.TerminalRef.AZIDText(),
		VerificationCodeExpiryTime: codeExpiryProto,
	}
	return &resp, nil
}

func (authServer *TerminalAuthorizationServiceServer) ConfirmTerminalAuthorization(
	callCtx context.Context, reqProto *iampb.ConfirmTerminalAuthorizationRequest,
) (*iampb.ConfirmTerminalAuthorizationResponse, error) {
	reqCtx, err := authServer.iamServerCore.GRPCOpInputContext(callCtx)
	if err != nil {
		panic(err) //TODO: translate and return the error
	}
	ctxAuth := reqCtx.Authorization()
	if ctxAuth.IsValid() {
		logCtx(reqCtx).
			Warn().Msgf("Authorization context must not be valid: %#v", ctxAuth)
		return nil, grpcstatus.Error(grpccodes.Unauthenticated, "")
	}

	termRef, err := iam.TerminalRefKeyFromAZIDText(reqProto.TerminalId)
	if err != nil {
		logCtx(reqCtx).
			Warn().Err(err).
			Msgf("Unable to parse terminal ID %q", reqProto.TerminalId)
		return nil, grpcstatus.Error(grpccodes.InvalidArgument, "")
	}

	termSecret, _, err := authServer.iamServerCore.
		ConfirmTerminalAuthorization(
			reqCtx, termRef, reqProto.VerificationCode)
	if err != nil {
		logCtx(reqCtx).
			Warn().Err(err).
			Msgf("Terminal authorization confirm failed: %v")
		return nil, grpcerrs.Error(err)
	}

	return &iampb.ConfirmTerminalAuthorizationResponse{
		TerminalSecret: termSecret,
	}, nil
}

func (authServer *TerminalAuthorizationServiceServer) GenerateAccessTokenByTerminalCredentials(
	callCtx context.Context, reqProto *iampb.GenerateAccessTokenByTerminalCredentialsRequest,
) (*iampb.GenerateAccessTokenByTerminalCredentialsResponse, error) {
	reqCtx, err := authServer.iamServerCore.GRPCOpInputContext(callCtx)
	if err != nil {
		panic(err) //TODO: translate and return the error
	}
	ctxAuth := reqCtx.Authorization()
	if ctxAuth.IsValid() {
		logCtx(reqCtx).
			Warn().Msgf("Authorization context must not be valid: %#v", ctxAuth)
		return nil, grpcstatus.Error(grpccodes.Unauthenticated, "")
	}

	termRef, err := iam.TerminalRefKeyFromAZIDText(reqProto.TerminalId)
	if err != nil {
		logCtx(reqCtx).
			Warn().Err(err).Str("terminal", reqProto.TerminalId).
			Msg("Terminal ID parsing")
		return nil, grpcstatus.Error(grpccodes.InvalidArgument, "")
	}

	authOK, userRef, err := authServer.iamServerCore.
		AuthenticateTerminal(termRef, reqProto.TerminalSecret)
	if err != nil {
		logCtx(reqCtx).
			Warn().Err(err).Str("terminal", termRef.AZIDText()).
			Msg("Terminal authentication")
		return nil, grpcerrs.Error(err)
	}
	if !authOK {
		logCtx(reqCtx).
			Warn().Str("terminal", termRef.AZIDText()).Msg("Terminal authentication")
		return nil, grpcstatus.Error(grpccodes.InvalidArgument, "")
	}

	if userRef.IsSound() {
		userInstInfo, err := authServer.iamServerCore.UserService.
			GetUserInstanceInfo(reqCtx, userRef)
		if err != nil {
			logCtx(reqCtx).
				Warn().Err(err).Str("terminal", termRef.AZIDText()).
				Msg("Terminal user account state")
			return nil, grpcerrs.Error(err)
		}
		if userInstInfo == nil || !userInstInfo.IsActive() {
			var status string
			if userInstInfo == nil {
				status = "not exist"
			} else {
				status = "deleted"
			}
			logCtx(reqCtx).
				Warn().Str("terminal", termRef.AZIDText()).Str("user", userRef.AZIDText()).
				Msg("Terminal user account " + status)
			return nil, grpcstatus.Error(grpccodes.InvalidArgument, "")
		}
	}

	issueTime := time.Now().UTC()
	tokenString, err := authServer.iamServerCore.
		GenerateAccessTokenJWT(reqCtx, termRef, userRef, issueTime)
	if err != nil {
		panic(err)
	}

	return &iampb.GenerateAccessTokenByTerminalCredentialsResponse{
		AccessToken: tokenString,
		AuthorizationData: &iampb.AuthorizationData{
			SubjectUserId: userRef.AZIDText(),
		},
	}, nil
}
