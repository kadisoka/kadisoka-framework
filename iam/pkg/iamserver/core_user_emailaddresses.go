package iamserver

import (
	"database/sql"
	"time"

	"github.com/kadisoka/foundation/pkg/errors"
	"github.com/lib/pq"

	"github.com/kadisoka/iam/pkg/iam"
	"github.com/kadisoka/iam/pkg/iamserver/eav10n"
)

//TODO(exa): there should be getters for different purpose (e.g.,
// for login / primary, for display / contact, for actual mailing, for recovery, etc)
func (core *Core) GetUserPrimaryEmailAddress(
	callCtx iam.CallContext,
	userID iam.UserID,
) (*iam.EmailAddress, error) {
	var rawInput string
	err := core.db.
		QueryRow(
			`SELECT raw_input `+
				`FROM user_email_addresses `+
				`WHERE user_id=$1 AND is_primary IS TRUE `+
				`AND deletion_time IS NULL AND verification_time IS NOT NULL`,
			userID).
		Scan(&rawInput)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	emailAddress, err := iam.EmailAddressFromString(rawInput)
	if err != nil {
		panic(err)
	}
	return &emailAddress, nil
}

// The ID of the user which provided email address is their verified primary.
func (core *Core) getUserIDByPrimaryEmailAddress(
	emailAddress iam.EmailAddress,
) (ownerUserID iam.UserID, err error) {
	queryStr :=
		`SELECT user_id ` +
			`FROM user_email_addresses ` +
			`WHERE local_part = $1 AND domain_part = $2 ` +
			`AND is_primary IS TRUE AND deletion_time IS NULL ` +
			`AND verification_time IS NOT NULL`
	err = core.db.
		QueryRow(queryStr,
			emailAddress.LocalPart(),
			emailAddress.DomainPart()).
		Scan(&ownerUserID)
	if err != nil {
		if err == sql.ErrNoRows {
			return iam.UserIDZero, nil
		}
		return iam.UserIDZero, err
	}

	return
}

// The ID of the user which provided email address is their primary,
// verified or not.
func (core *Core) getUserIDByPrimaryEmailAddressAllowUnverified(
	emailAddress iam.EmailAddress,
) (ownerUserID iam.UserID, verified bool, err error) {
	queryStr :=
		`SELECT user_id, CASE WHEN verification_time IS NULL THEN false ELSE true END AS verified ` +
			`FROM user_email_addresses ` +
			`WHERE local_part = $1 AND domain_part = $2 ` +
			`AND is_primary IS TRUE AND deletion_time IS NULL ` +
			`ORDER BY creation_time DESC LIMIT 1`
	err = core.db.
		QueryRow(queryStr,
			emailAddress.LocalPart(),
			emailAddress.DomainPart()).
		Scan(&ownerUserID, &verified)
	if err != nil {
		if err == sql.ErrNoRows {
			return iam.UserIDZero, false, nil
		}
		return iam.UserIDZero, false, err
	}

	return
}

func (core *Core) SetUserPrimaryEmailAddress(
	callCtx iam.CallContext,
	userID iam.UserID,
	emailAddress iam.EmailAddress,
	verificationMethods []eav10n.VerificationMethod,
) (verificationID int64, codeExpiry *time.Time, err error) {
	authCtx := callCtx.Authorization()
	if !authCtx.IsUserContext() {
		return 0, nil, iam.ErrUserContextRequired
	}
	// Don't allow changing other user's for now
	if userID != authCtx.UserID {
		return 0, nil, iam.ErrContextUserNotAllowedToPerformActionOnResource
	}

	existingOwnerUserID, err := core.
		getUserIDByPrimaryEmailAddress(emailAddress)
	if err != nil {
		return 0, nil, errors.Wrap("getUserIDByPrimaryEmailAddress", err)
	}
	if existingOwnerUserID.IsValid() {
		if existingOwnerUserID != authCtx.UserID {
			return 0, nil, errors.ArgMsg("emailAddress", "conflict")
		}
		return 0, nil, nil
	}

	alreadyVerified, err := core.setUserPrimaryEmailAddress(
		authCtx.Actor(), authCtx.UserID, emailAddress)
	if err != nil {
		panic(err)
	}
	if alreadyVerified {
		return 0, nil, nil
	}

	//TODO: user-set has higher priority over terminal's
	userLanguages, err := core.getTerminalAcceptLanguages(authCtx.TerminalID())

	verificationID, codeExpiry, err = core.eaVerifier.
		StartVerification(callCtx, emailAddress,
			0, userLanguages, verificationMethods)
	if err != nil {
		switch err.(type) {
		case eav10n.InvalidEmailAddressError:
			return 0, nil, errors.Arg("emailAddress", err)
		}
		return 0, nil, errors.Wrap("eaVerifier.StartVerification", err)
	}

	return
}

func (core *Core) setUserPrimaryEmailAddress(
	actor iam.Actor,
	userID iam.UserID,
	emailAddress iam.EmailAddress,
) (alreadyVerified bool, err error) {
	tNow := time.Now().UTC()

	xres, err := core.db.Exec(
		`INSERT INTO user_email_addresses (`+
			`user_id, local_part, domain_part, raw_input, is_primary, `+
			`creation_time, creation_user_id, creation_terminal_id `+
			`) VALUES (`+
			`$1, $2, $3, $4, $5, $6, $7, $8`+
			`) `+
			`ON CONFLICT (user_id, local_part, domain_part) WHERE deletion_time IS NULL `+
			`DO NOTHING`,
		userID,
		emailAddress.LocalPart(),
		emailAddress.DomainPart(),
		emailAddress.RawInput(),
		true,
		tNow,
		actor.UserID,
		actor.TerminalID)
	if err != nil {
		return false, err
	}

	n, err := xres.RowsAffected()
	if err != nil {
		return false, err
	}
	if n == 1 {
		return false, nil
	}

	err = core.db.QueryRow(
		`SELECT CASE WHEN verification_time IS NULL THEN false ELSE true END AS verified `+
			`FROM user_email_addresses `+
			`WHERE user_id = $1 AND local_part = $2 AND domain_part = $3 AND is_primary IS TRUE`,
		userID, emailAddress.LocalPart(), emailAddress.DomainPart()).
		Scan(&alreadyVerified)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}
		return false, err
	}

	return
}

func (core *Core) ConfirmUserEmailAddressVerification(
	callCtx iam.CallContext,
	verificationID int64,
	code string,
) (updated bool, err error) {
	authCtx := callCtx.Authorization()
	err = core.eaVerifier.ConfirmVerification(
		callCtx, verificationID, code)
	if err != nil {
		switch err {
		case eav10n.ErrVerificationCodeMismatch:
			return false, errors.ArgMsg("code", "mismatch")
		case eav10n.ErrVerificationCodeExpired:
			return false, errors.ArgMsg("code", "expired")
		}
		return false, errors.Wrap("eaVerifier.ConfirmVerification", err)
	}

	tNow := time.Now().UTC()
	emailAddress, err := core.eaVerifier.
		GetEmailAddressByVerificationID(verificationID)
	// An unexpected condition which could cause bad state
	if err != nil {
		panic(err)
	}

	updated, err = core.
		ensureUserEmailAddressVerifiedFlag(
			authCtx.UserID, *emailAddress,
			&tNow, verificationID)
	if err != nil {
		panic(err)
	}

	return updated, nil
}

func (core *Core) ensureUserEmailAddressVerifiedFlag(
	userID iam.UserID,
	emailAddress iam.EmailAddress,
	verificationTime *time.Time,
	verificationID int64,
) (bool, error) {
	var err error
	var xres sql.Result

	xres, err = core.db.Exec(
		`UPDATE user_email_addresses SET (`+
			`verification_time, verification_id`+
			`) = ( `+
			`$1, $2`+
			`) WHERE user_id = $3 `+
			`AND local_part = $4 AND domain_part = $5 `+
			`AND deletion_time IS NULL AND verification_time IS NULL`,
		verificationTime,
		verificationID,
		userID,
		emailAddress.LocalPart(),
		emailAddress.DomainPart())
	if err != nil {
		pqErr, _ := err.(*pq.Error)
		if pqErr != nil &&
			pqErr.Code == "23505" &&
			pqErr.Constraint == "user_email_addresses_local_part_domain_part_uidx" {
			return false, errors.ArgMsg("emailAddress", "conflict")
		}
		return false, err
	}

	var n int64
	n, err = xres.RowsAffected()
	if err != nil {
		return false, err
	}
	return n == 1, nil
}
