package entity_manager

import "testing"

func TestBasicCapabilities(t *testing.T) {
	s := []struct {
		ID         string
		uidNumber  int32
		secret     string
		capability string
		test       string
		err        error
	}{
		{"foo", -1, "a", "GLOBAL_ROOT", "GLOBAL_ROOT", nil},
		{"bar", -1, "a", "CREATE_ENTITY", "CREATE_ENTITY", nil},
		{"baz", -1, "a", "GLOBAL_ROOT", "CREATE_ENTITY", nil},
		{"bam", -1, "a", "CREATE_ENTITY", "GLOBAL_ROOT", E_ENTITY_UNQUALIFIED},
	}

	for _, c := range s {
		if err := NewEntity(c.ID, c.uidNumber, c.secret); err != nil {
			t.Error(err)
		}

		e, err := GetEntityByID(c.ID)
		if err != nil {
			t.Error(err)
		}
		
		SetCapability(e, c.capability)

		if err = checkCapability(e, c.test); err != c.err {
			t.Error(err)
		}
	}
}

func TestDuplicateCapabilityAdd(t *testing.T) {

}