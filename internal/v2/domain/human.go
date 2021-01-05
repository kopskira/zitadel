package domain

import (
	"github.com/caos/zitadel/internal/crypto"
	caos_errors "github.com/caos/zitadel/internal/errors"
	es_models "github.com/caos/zitadel/internal/eventstore/models"
	"strings"
	"time"
)

type Human struct {
	es_models.ObjectRoot

	*Password
	*Profile
	*Email
	*Phone
	*Address
	ExternalIDPs       []*ExternalIDP
	InitCode           *InitUserCode
	EmailCode          *EmailCode
	PhoneCode          *PhoneCode
	PasswordCode       *PasswordCode
	OTP                *OTP
	U2FTokens          []*WebAuthNToken
	PasswordlessTokens []*WebAuthNToken
	U2FLogins          []*WebAuthNLogin
	PasswordlessLogins []*WebAuthNLogin
}

type InitUserCode struct {
	es_models.ObjectRoot

	Code   *crypto.CryptoValue
	Expiry time.Duration
}

type Gender int32

const (
	GenderUnspecified Gender = iota
	GenderFemale
	GenderMale
	GenderDiverse

	genderCount
)

func (f Gender) Valid() bool {
	return f >= 0 && f < genderCount
}

func (u *Human) IsValid() bool {
	return u.Profile != nil && u.FirstName != "" && u.LastName != "" && u.Email != nil && u.Email.IsValid() && u.Phone == nil || (u.Phone != nil && u.Phone.PhoneNumber != "" && u.Phone.IsValid())
}

func (u *Human) CheckOrgIAMPolicy(userName string, policy *OrgIAMPolicy) error {
	if policy == nil {
		return caos_errors.ThrowPreconditionFailed(nil, "DOMAIN-zSH7j", "Errors.Users.OrgIamPolicyNil")
	}
	if policy.UserLoginMustBeDomain && strings.Contains(userName, "@") {
		return caos_errors.ThrowPreconditionFailed(nil, "DOMAIN-se4sJ", "Errors.User.EmailAsUsernameNotAllowed")
	}
	if !policy.UserLoginMustBeDomain && u.Profile != nil && userName == "" && u.Email != nil {
		userName = u.EmailAddress
	}
	return nil
}

func (u *Human) SetNamesAsDisplayname() {
	if u.Profile != nil && u.DisplayName == "" && u.FirstName != "" && u.LastName != "" {
		u.DisplayName = u.FirstName + " " + u.LastName
	}
}

func (u *Human) HashPasswordIfExisting(policy *PasswordComplexityPolicy, passwordAlg crypto.HashAlgorithm, onetime bool) error {
	if u.Password != nil {
		return u.Password.HashPasswordIfExisting(policy, passwordAlg, onetime)
	}
	return nil
}