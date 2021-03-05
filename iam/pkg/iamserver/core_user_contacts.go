package iamserver

import (
	"github.com/kadisoka/kadisoka-framework/iam/pkg/iam"
)

func (core *Core) GetUserContactUserIDs(
	callCtx iam.CallContext,
	userRef iam.UserRefKey,
) ([]iam.UserRefKey, error) {
	userIDRows, err := core.db.
		Query(
			`SELECT DISTINCT `+
				`ph.user_id `+
				`FROM user_contact_phone_numbers AS cp `+
				`JOIN `+userIdentifierPhoneNumberTableName+` AS ph ON `+
				`  ph.country_code = cp.contact_country_code `+
				`  AND ph.national_number = cp.contact_national_number `+
				`  AND ph.deletion_time IS NULL `+
				`  AND ph.verification_time IS NOT NULL `+
				`JOIN users AS usr ON `+
				`  usr.id = ph.user_id `+
				`  AND usr.deletion_time IS NULL `+
				`WHERE `+
				`  cp.user_id = $1 `+
				`ORDER BY ph.user_id ASC`,
			userRef.ID().PrimitiveValue())
	if err != nil {
		return nil, err
	}
	defer userIDRows.Close()

	var userRefs []iam.UserRefKey
	for userIDRows.Next() {
		uid := iam.UserIDZero
		err = userIDRows.Scan(&uid)
		if err != nil {
			panic(err)
		}
		userRefs = append(userRefs, iam.NewUserRefKey(uid))
	}
	if err = userIDRows.Err(); err != nil {
		return nil, err
	}

	return userRefs, nil
}
