package ecom

import (
	"testing"
)

func TestModifierSeparation(t *testing.T) {
	orig := GetCreditCard(ValidVisa)

	origExpiration := orig.Expiration

	newCard := orig.MakeCardExpirationUnique()

	if newCard.Expiration == orig.Expiration {
		t.Errorf("New card has the same expiration: %s=%s", orig.Expiration, newCard.Expiration)
	}

	copy := GetCreditCard(ValidVisa)
	if copy.Expiration != origExpiration {
		t.Errorf("Instance of original has changed: %s!=%s", origExpiration, copy.Expiration)
	}
}

func TestUniqueExpiration(t *testing.T) {
	expirations := make(map[string]struct{}, 0)

	cc := GetCreditCard(ValidVisa)

	expirations[cc.Expiration] = struct{}{}

	n1 := cc.MakeCardExpirationUnique()

	if v,ok := expirations[n1.Expiration] ; ok {
		t.Errorf("Expiration already exists: %s=%s", v, n1.Expiration)
	}

	expirations[n1.Expiration] = struct{}{}

	n2 := cc.MakeCardExpirationUnique()

	if v,ok := expirations[n2.Expiration] ; ok {
		t.Errorf("Expiration already exists: %s=%s", v, n2.Expiration)
	}

	expirations[n2.Expiration] = struct{}{}

	n3 := n1.MakeCardExpirationUnique()

	if v,ok := expirations[n3.Expiration] ; ok {
		t.Errorf("Expiration already exists: %s=%s", v, n3.Expiration)
	}



}