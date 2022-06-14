package iam

func (idNum UserIDNum) IsNormalAccount() bool {
	return idNum.IsStaticallyValid() && !idNum.HasBotBits()
}
