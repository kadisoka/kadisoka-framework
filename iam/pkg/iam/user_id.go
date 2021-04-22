package iam

func (idNum UserIDNum) IsNormalAccount() bool {
	return idNum.IsSound() && !idNum.HasBotBits()
}
