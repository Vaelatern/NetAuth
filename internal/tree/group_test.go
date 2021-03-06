package tree

import (
	"testing"

	"github.com/golang/protobuf/proto"

	"github.com/NetAuth/NetAuth/internal/db"

	pb "github.com/NetAuth/Protocol"
)

func TestNewGroup(t *testing.T) {
	em := getNewEntityManager(t)

	s := []struct {
		name        string
		displayName string
		number      int32
		wantErr     error
	}{
		{"fooGroup", "", 1, nil},
		{"fooGroup", "", 1, ErrDuplicateGroupName},
		{"barGroup", "", -1, nil},
	}
	for _, c := range s {
		if err := em.NewGroup(c.name, c.displayName, "", c.number); err != c.wantErr {
			t.Errorf("Wrong Error: want '%v' got '%v'", c.wantErr, err)
		}
	}
}

func TestListGroups(t *testing.T) {
	em := getNewEntityManager(t)

	names := []string{"aaa", "aab", "aac", "aad", "aae"}
	for _, n := range names {
		if err := em.NewGroup(n, "", "", -1); err != nil {
			t.Fatal(err)
		}
	}

	grps, err := em.ListGroups()
	if err != nil {
		t.Fatal(err)
	}
	if len(grps) != len(names) {
		t.Fatal("Wrong number of groups")
	}
}

func TestDeleteGroup(t *testing.T) {
	em := getNewEntityManager(t)

	if err := em.NewGroup("foo", "", "", -1); err != nil {
		t.Error(err)
	}

	if _, err := em.GetGroupByName("foo"); err != nil {
		t.Error(err)
	}

	if err := em.DeleteGroup("foo"); err != nil {
		t.Error(err)
	}

	if _, err := em.GetGroupByName("foo"); err != db.ErrUnknownGroup {
		t.Error(err)
	}
}

func TestUpdateGroupMetaInternal(t *testing.T) {
	em := getNewEntityManager(t)

	if err := em.NewGroup("foo", "foo", "", -1); err != nil {
		t.Error(err)
	}

	update := &pb.Group{DisplayName: proto.String("Foo Group")}

	if err := em.UpdateGroupMeta("foo", update); err != nil {
		t.Error(err)
	}

	g, err := em.GetGroupByName("foo")
	if err != nil {
		t.Error(err)
	}

	if g.GetDisplayName() != "Foo Group" {
		t.Error("Meta update failed!")
	}
}

func TestManageUntypedGroupMeta(t *testing.T) {
	em := getNewEntityManager(t)

	if err := em.NewGroup("g1","g1", "", -1); err != nil {
		t.Fatal(err)
	}

	cases := []struct {
		groupName string
		mode      string
		key       string
		value     string
		wantKV    []string
		wantErr   error
	}{
		{"g1", "upsert", "k1", "v1", nil, nil},
		{"g1", "read", "*", "", []string{"k1:v1"}, nil},
		{"unknown", "read", "*", "", nil, db.ErrUnknownGroup},
	}

	for i, c := range cases {
		kv, err := em.ManageUntypedGroupMeta(c.groupName, c.mode, c.key, c.value)
		if err != c.wantErr {
			t.Fatalf("%d: Got: %v; Want: %v", i, err, c.wantErr)
		}

		// This function is defined in helpers_test.go
		if !slicesAreEqual(kv, c.wantKV) {
			t.Fatalf("%d: Got: %v; Want: %v", i, kv, c.wantKV)
		}
	}
}

func TestSetSameGroupCapabilityTwice(t *testing.T) {
	em := getNewEntityManager(t)

	// Add an entity
	if err := em.NewGroup("foo", "", "", -1); err != nil {
		t.Error(err)
	}

	e, err := em.db.LoadGroup("foo")
	if err != nil {
		t.Error(err)
	}

	// Set one capability
	em.setGroupCapability(e, "GLOBAL_ROOT")
	if len(e.Capabilities) != 1 {
		t.Error("Wrong number of capabilities set!")
	}

	// Set it again and make sure its still only listed once.
	em.setGroupCapability(e, "GLOBAL_ROOT")
	if len(e.Capabilities) != 1 {
		t.Error("Wrong number of capabilities set!")
	}
}

func TestSetGroupCapabilityBogusGroup(t *testing.T) {
	em := getNewEntityManager(t)

	if err := em.SetGroupCapabilityByName("foo", "GLOBAL_ROOT"); err != db.ErrUnknownGroup {
		t.Error(err)
	}
}

func TestSetGroupCapabilityNoCap(t *testing.T) {
	em := getNewEntityManager(t)

	if err := em.NewGroup("foo", "", "", -1); err != nil {
		t.Fatal(err)
	}

	if err := em.SetGroupCapabilityByName("foo", ""); err != ErrUnknownCapability {
		t.Error(err)
	}
}

func TestRemoveGroupCapability(t *testing.T) {
	em := getNewEntityManager(t)

	// Add an entity
	if err := em.NewGroup("foo", "", "", -1); err != nil {
		t.Error(err)
	}

	e, err := em.db.LoadGroup("foo")
	if err != nil {
		t.Error(err)
	}

	// Set one capability
	em.setGroupCapability(e, "GLOBAL_ROOT")
	if len(e.Capabilities) != 1 {
		t.Error("Wrong number of capabilities set!")
	}
	// Set another capability
	em.setGroupCapability(e, "MODIFY_ENTITY_META")
	if len(e.Capabilities) != 2 {
		t.Error("Wrong number of capabilities set!")
	}

	// Remove it and make sure its gone
	em.removeGroupCapability(e, "GLOBAL_ROOT")
	if len(e.Capabilities) != 1 {
		t.Error("Wrong number of capabilities set!")
	}
}

func TestRemoveGroupCapabilityBogusGroup(t *testing.T) {
	em := getNewEntityManager(t)

	if err := em.RemoveGroupCapabilityByName("foo", "GLOBAL_ROOT"); err != db.ErrUnknownGroup {
		t.Error(err)
	}
}

func TestRemoveGroupCapabilityNoCap(t *testing.T) {
	em := getNewEntityManager(t)

	if err := em.NewGroup("foo", "", "", -1); err != nil {
		t.Fatal(err)
	}

	if err := em.RemoveGroupCapabilityByName("foo", ""); err != ErrUnknownCapability {
		t.Error(err)
	}
}
