package iamserver

import (
	"crypto/rand"
	"crypto/subtle"
	"database/sql"
	"encoding/base64"
	"errors"
	"fmt"
	"strings"

	"github.com/jmoiron/sqlx"
	"golang.org/x/crypto/argon2"

	"github.com/kadisoka/kadisoka-framework/iam/pkg/iam"
)

var (
	ErrPasswordHashFormatInvalid       = errors.New("hash format invalid")
	ErrPasswordHashVersionIncompatible = errors.New("hash version incompatible")
)

var argon2PasswordHashParamsEncoding = base64.RawStdEncoding

type argon2PasswordHashingParams struct {
	Memory      uint32
	Iterations  uint32
	Parallelism uint8
	SaltLength  uint32
	KeyLength   uint32
}

//TODO:make this configurable
var argon2PasswordHashingParamsDefault = argon2PasswordHashingParams{
	Memory:      64 * 1024,
	Iterations:  3,
	Parallelism: 2,
	SaltLength:  16,
	KeyLength:   32,
}

const userPasswordTableName = "user_password_dt"

func (core *Core) SetUserPassword(
	callCtx iam.CallContext,
	userRef iam.UserRefKey,
	clearTextPassword string,
) error {
	authCtx := callCtx.Authorization()
	if !authCtx.IsUserContext() || !userRef.EqualsUserRefKey(authCtx.UserRef()) {
		return errors.New("forbidden")
	}

	ctxTime := callCtx.RequestInfo().ReceiveTime

	passwordHash, err := core.hashPassword(clearTextPassword)
	if err != nil {
		return err
	}

	return doTx(core.db, func(tx *sqlx.Tx) error {
		_, txErr := core.db.Exec(
			`UPDATE `+userPasswordTableName+` SET `+
				`d_ts = $1, d_uid = $2, d_tid = $3 `+
				`WHERE user_id = $4 AND d_ts IS NULL`,
			ctxTime, authCtx.UserID().PrimitiveValue(), authCtx.TerminalID().PrimitiveValue(),
			userRef.ID().PrimitiveValue())
		if txErr != nil {
			return txErr
		}
		_, txErr = core.db.Exec(
			`INSERT INTO `+userPasswordTableName+` `+
				`(user_id, password, c_ts, c_uid, c_tid) `+
				`VALUES ($1, $2, $3, $4, $5) `,
			userRef.ID().PrimitiveValue(), passwordHash,
			ctxTime, authCtx.UserID().PrimitiveValue(), authCtx.TerminalID().PrimitiveValue())
		return nil
	})
}

func (core *Core) MatchUserPassword(
	userRef iam.UserRefKey,
	clearTextPassword string,
) (ok bool, err error) {
	passwordHash, err := core.getUserPasswordHash(userRef.ID())
	if err != nil {
		return false, err
	}
	if passwordHash == "" && clearTextPassword == passwordHash {
		return true, err
	}
	return core.matchPasswordAndPasswordHash(clearTextPassword, passwordHash)
}

func (core *Core) getUserPasswordHash(
	userID iam.UserID,
) (hashedPassword string, err error) {
	err = core.db.
		QueryRow(
			`SELECT password `+
				`FROM `+userPasswordTableName+` `+
				`WHERE user_id = $1 AND d_ts IS NULL`,
			userID).
		Scan(&hashedPassword)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", nil
		}
		return "", err
	}

	return hashedPassword, nil
}

func (core *Core) hashPassword(
	password string,
) (encodedPasswordHash string, err error) {
	params := argon2PasswordHashingParamsDefault

	// generate a chryptographically secure random salt
	salt, err := core.generatePasswordSalt(params.SaltLength)
	if err != nil {
		return "", err
	}

	hash := argon2.IDKey(
		[]byte(password),
		salt,
		params.Iterations,
		params.Memory,
		params.Parallelism,
		params.KeyLength,
	)

	// Base64 encode the salt and hashed password.
	saltB64 := argon2PasswordHashParamsEncoding.EncodeToString(salt)
	hashB64 := argon2PasswordHashParamsEncoding.EncodeToString(hash)

	// Return string using the standard encoded hash representation.
	encodedPasswordHash = fmt.Sprintf(
		"$argon2id$v=%d$m=%d,t=%d,p=%d$%s$%s",
		argon2.Version, params.Memory, params.Iterations, params.Parallelism,
		saltB64, hashB64)

	return encodedPasswordHash, nil
}

func (core *Core) generatePasswordSalt(n uint32) ([]byte, error) {
	b := make([]byte, n)
	_, err := rand.Read(b)
	if err != nil {
		return nil, err
	}
	return b, nil
}

func (core *Core) matchPasswordAndPasswordHash(
	clearTextPassword, encodedPasswordHash string,
) (match bool, err error) {
	// Extract the parameters, salt and derived key from the encoded password
	// hash
	if encodedPasswordHash == "" {
		return false, nil
	}

	params, salt, hash, err := core.
		decodePasswordHash(encodedPasswordHash)
	if err != nil {
		return false, err
	}

	// Derive the key from the other password using the same parameters
	otherHash := argon2.IDKey([]byte(clearTextPassword), salt,
		params.Iterations, params.Memory, params.Parallelism, params.KeyLength)

	if subtle.ConstantTimeCompare(hash, otherHash) != 1 {
		return false, nil
	}

	return true, nil
}

func (core *Core) decodePasswordHash(
	encodedPasswordHash string,
) (params *argon2PasswordHashingParams, salt, hash []byte, err error) {
	vals := strings.Split(encodedPasswordHash, "$")

	if len(vals) != 6 {
		return nil, nil, nil, ErrPasswordHashFormatInvalid
	}

	var version int
	_, err = fmt.Sscanf(vals[2], "v=%d", &version)
	if err != nil {
		return nil, nil, nil, err
	}

	if version != argon2.Version {
		return nil, nil, nil, ErrPasswordHashVersionIncompatible
	}

	params = &argon2PasswordHashingParams{}
	_, err = fmt.Sscanf(vals[3],
		"m=%d,t=%d,p=%d",
		&params.Memory, &params.Iterations, &params.Parallelism)
	if err != nil {
		return nil, nil, nil, err
	}

	salt, err = argon2PasswordHashParamsEncoding.DecodeString(vals[4])
	if err != nil {
		return nil, nil, nil, err
	}

	params.SaltLength = uint32(len(salt))

	hash, err = argon2PasswordHashParamsEncoding.DecodeString(vals[5])
	if err != nil {
		return nil, nil, nil, err
	}

	params.KeyLength = uint32(len(hash))

	return params, salt, hash, nil
}
