package iam

func (idNum UserIDNum) IsNormalAccount() bool {
	return idNum.IsValid() && !idNum.HasBotBits()
}
