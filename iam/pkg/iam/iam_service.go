package iam

type IAMService interface {
	ServiceClient

	UserServiceInternal

	TerminalService

	// This below is reserverd for S2S services.
	TerminalFCMRegistrationTokenService
}
