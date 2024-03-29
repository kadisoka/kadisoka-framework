package grpc

import (
	"context"

	"github.com/alloyzeus/go-azfl/errors"
	iampb "github.com/alloyzeus/go-azgrpc/azgrpc/iam/v1"
	pbtypes "github.com/gogo/protobuf/types"
	"google.golang.org/grpc"
	grpccodes "google.golang.org/grpc/codes"
	grpcstatus "google.golang.org/grpc/status"

	grpcerrs "github.com/kadisoka/kadisoka-framework/pkg/foundation/pkg/api/grpc/errors"
	"github.com/kadisoka/kadisoka-framework/pkg/iam/pkg/iam"
	"github.com/kadisoka/kadisoka-framework/pkg/iam/pkg/iamserver"
	"github.com/kadisoka/kadisoka-framework/pkg/volib/pkg/telephony"
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
	inputCtx context.Context,
	reqProto *iampb.InitiateUserTerminalAuthorizationByPhoneNumberRequest,
) (*iampb.InitiateUserTerminalAuthorizationByPhoneNumberResponse, error) {
	reqCtx, err := authServer.iamServerCore.GRPCCallInputContext(inputCtx)
	if err != nil {
		panic(err) //TODO: translate and return the error
	}
	ctxAuth := reqCtx.Authorization()
	if ctxAuth.IsStaticallyValid() {
		logCtx(reqCtx).
			Warn().Msgf("Authorization context must not be valid: %#v", reqCtx)
		return nil, grpcstatus.Error(grpccodes.Unauthenticated, "")
	}

	appID, err := iam.ApplicationIDFromAZIDText(reqProto.ClientCredentials.ClientId)
	if err != nil {
		panic(err)
	}

	app, err := authServer.iamServerCore.
		AuthenticatedApplication(appID, reqProto.ClientCredentials.ClientSecret)
	if err != nil {
		panic(err)
	}

	if app == nil {
		logCtx(reqCtx).
			Warn().Msgf("Client authentication failed")
		return nil, grpcstatus.Error(grpccodes.Unauthenticated, "")
	}

	phoneNumber, err := telephony.PhoneNumberFromString(reqProto.PhoneNumber)
	if err != nil {
		logCtx(reqCtx).
			Warn().Err(err).Str("phone_number", reqProto.PhoneNumber).
			Msg("Phone number format")
		return nil, grpcstatus.Error(grpccodes.InvalidArgument, "")
	}

	authStartOutCtx, authStartOutData := authServer.iamServerCore.
		StartTerminalAuthorizationByPhoneNumber(
			reqCtx,
			iamserver.TerminalAuthorizationByPhoneNumberStartInputData{
				PhoneNumber:         phoneNumber,
				VerificationMethods: nil,
				TerminalAuthorizationStartInputBaseData: iamserver.TerminalAuthorizationStartInputBaseData{
					ApplicationID: appID,
					DisplayName:   reqProto.TerminalInfo.DisplayName,
				},
			})
	if err = authStartOutCtx.Err; err != nil {
		switch err.(type) {
		case errors.CallError:
			logCtx(reqCtx).
				Warn().Err(err).
				Msgf("StartTerminalAuthorizationByPhoneNumber %v",
					phoneNumber)
			return nil, grpcstatus.Error(grpccodes.InvalidArgument, "")
		}
		logCtx(reqCtx).
			Error().Err(err).
			Msgf("StartTerminalAuthorizationByPhoneNumber %v",
				phoneNumber)
		return nil, grpcerrs.Error(err)
	}

	var codeExpiryProto *pbtypes.Timestamp
	if codeExpiry := authStartOutData.VerificationCodeExpiryTime; codeExpiry != nil {
		codeExpiryProto, err = pbtypes.TimestampProto(*codeExpiry)
		if err != nil {
			panic(err)
		}
	}
	resp := iampb.InitiateUserTerminalAuthorizationByPhoneNumberResponse{
		TerminalId:                 authStartOutData.TerminalID.AZIDText(),
		VerificationCodeExpiryTime: codeExpiryProto,
	}
	return &resp, nil
}

func (authServer *TerminalAuthorizationServiceServer) ConfirmTerminalAuthorization(
	inputCtx context.Context, reqProto *iampb.ConfirmTerminalAuthorizationRequest,
) (*iampb.ConfirmTerminalAuthorizationResponse, error) {
	reqCtx, err := authServer.iamServerCore.GRPCCallInputContext(inputCtx)
	if err != nil {
		panic(err) //TODO: translate and return the error
	}
	ctxAuth := reqCtx.Authorization()
	if ctxAuth.IsStaticallyValid() {
		logCtx(reqCtx).
			Warn().Msgf("Authorization context must not be valid: %#v", ctxAuth)
		return nil, grpcstatus.Error(grpccodes.Unauthenticated, "")
	}

	termID, err := iam.TerminalIDFromAZIDText(reqProto.TerminalId)
	if err != nil {
		logCtx(reqCtx).
			Warn().Err(err).
			Msgf("Unable to parse terminal ID %q", reqProto.TerminalId)
		return nil, grpcstatus.Error(grpccodes.InvalidArgument, "")
	}

	termSecret, _, err := authServer.iamServerCore.
		ConfirmTerminalAuthorization(
			reqCtx, termID, reqProto.VerificationCode)
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
	inputCtx context.Context, reqProto *iampb.GenerateAccessTokenByTerminalCredentialsRequest,
) (*iampb.GenerateAccessTokenByTerminalCredentialsResponse, error) {
	reqCtx, err := authServer.iamServerCore.GRPCCallInputContext(inputCtx)
	if err != nil {
		panic(err) //TODO: translate and return the error
	}
	ctxAuth := reqCtx.Authorization()
	if ctxAuth.IsStaticallyValid() {
		logCtx(reqCtx).
			Warn().Msgf("Authorization context must not be valid: %#v", ctxAuth)
		return nil, grpcstatus.Error(grpccodes.Unauthenticated, "")
	}

	termID, err := iam.TerminalIDFromAZIDText(reqProto.TerminalId)
	if err != nil {
		logCtx(reqCtx).
			Warn().Err(err).Str("terminal", reqProto.TerminalId).
			Msg("Terminal ID parsing")
		return nil, grpcstatus.Error(grpccodes.InvalidArgument, "")
	}

	authOK, userID, err := authServer.iamServerCore.
		AuthenticateTerminal(termID, reqProto.TerminalSecret)
	if err != nil {
		logCtx(reqCtx).
			Warn().Err(err).Str("terminal", termID.AZIDText()).
			Msg("Terminal authentication")
		return nil, grpcerrs.Error(err)
	}
	if !authOK {
		logCtx(reqCtx).
			Warn().Str("terminal", termID.AZIDText()).Msg("Terminal authentication")
		return nil, grpcstatus.Error(grpccodes.InvalidArgument, "")
	}

	if userID.IsStaticallyValid() {
		userInstInfo, err := authServer.iamServerCore.UserService.
			GetUserInstanceInfo(reqCtx, userID)
		if err != nil {
			logCtx(reqCtx).
				Warn().Err(err).Str("terminal", termID.AZIDText()).
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
				Warn().Str("terminal", termID.AZIDText()).Str("user", userID.AZIDText()).
				Msg("Terminal user account " + status)
			return nil, grpcstatus.Error(grpccodes.InvalidArgument, "")
		}
	}

	tokenString, err := authServer.iamServerCore.
		GenerateAccessTokenJWT(reqCtx, termID, userID)
	if err != nil {
		panic(err)
	}

	return &iampb.GenerateAccessTokenByTerminalCredentialsResponse{
		AccessToken: tokenString,
		AuthorizationData: &iampb.AuthorizationData{
			SubjectUserId: userID.AZIDText(),
		},
	}, nil
}
