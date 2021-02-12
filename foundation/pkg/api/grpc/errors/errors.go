package errors

import (
	grpccodes "google.golang.org/grpc/codes"
	grpcstatus "google.golang.org/grpc/status"

	"github.com/kadisoka/kadisoka-framework/foundation/pkg/errors"
	accesserrs "github.com/kadisoka/kadisoka-framework/foundation/pkg/errors/access"
)

// Error translates err into gRPC error.
//
// This function attempts to use custom interface to obtain an
// error's description.
func Error(err error) error {
	if err == nil {
		return nil
	}

	if status, ok := grpcstatus.FromError(err); ok {
		return status.Err()
	}

	if code, ok := statusCode(err); ok {
		if code == grpccodes.OK {
			return nil
		}
		return grpcstatus.Error(code, statusDesc(err, ""))
	}

	if err == errors.ErrUnimplemented {
		return grpcstatus.Error(grpccodes.Unimplemented,
			statusDesc(err, ""))
	}

	// Ensure to sort the cases from specialized cases to generic cases
	switch typedErr := err.(type) {
	case accesserrs.Error:
		return grpcstatus.Error(grpccodes.PermissionDenied,
			statusDesc(err, ""))
	case errors.CallError:
		return requestError(typedErr)
	}

	return grpcstatus.Error(grpccodes.Internal,
		statusDesc(err, ""))
}

func statusCode(err error) (code grpccodes.Code, ok bool) {
	if err == nil {
		return grpccodes.OK, true
	}
	if x, ok := err.(interface{ GRPCStatusCode() grpccodes.Code }); ok && x != nil {
		code := x.GRPCStatusCode()
		return code, true
	}
	return grpccodes.Unknown, false
}

func statusDesc(err error, defDesc string) string {
	if d, ok := err.(interface{ GRPCStatusDescription() string }); ok && d != nil {
		return d.GRPCStatusDescription()
	}
	return defDesc
}

func requestError(err errors.CallError) error {
	if err == nil {
		return nil
	}

	// switch err.(type) {
	// case
	// }

	return grpcstatus.Error(grpccodes.InvalidArgument,
		statusDesc(err, ""))
}
