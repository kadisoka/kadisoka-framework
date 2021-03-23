package iamserver

import (
	"github.com/kadisoka/kadisoka-framework/iam/pkg/iam"
)

const userContactPhoneNumberDBTableName = "user_contact_phone_number_dt"

func (core *Core) GetUserContactUserIDs(
	callCtx iam.CallContext,
	userRef iam.UserRefKey,
) ([]iam.UserRefKey, error) {
	userIDRows, err := core.db.
		Query(
			`SELECT DISTINCT `+
				`ph.user_id `+
				`FROM `+userContactPhoneNumberDBTableName+` AS cp `+
				`JOIN `+userKeyPhoneNumberDBTableName+` AS ph ON `+
				`  ph.country_code = cp.contact_country_code `+
				`  AND ph.national_number = cp.contact_national_number `+
				`  AND ph.d_ts IS NULL `+
				`  AND ph.verification_ts IS NOT NULL `+
				`JOIN `+userDBTableName+` AS usr ON `+
				`  usr.id = ph.user_id `+
				`  AND usr.d_ts IS NULL `+
				`WHERE `+
				`  cp.user_id = $1 `+
				`ORDER BY ph.user_id ASC`,
			userRef.IDNum().PrimitiveValue())
	if err != nil {
		return nil, err
	}
	defer userIDRows.Close()

	var userRefs []iam.UserRefKey
	for userIDRows.Next() {
		userIDNum := iam.UserIDNumZero
		err = userIDRows.Scan(&userIDNum)
		if err != nil {
			panic(err)
		}
		userRefs = append(userRefs, iam.NewUserRefKey(userIDNum))
	}
	if err = userIDRows.Err(); err != nil {
		return nil, err
	}

	return userRefs, nil
}
