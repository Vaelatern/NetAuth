package memdb

import (
	"testing"

	"github.com/golang/protobuf/proto"

	"github.com/NetAuth/NetAuth/internal/db"

	pb "github.com/NetAuth/Protocol"
)

func TestDiscoverEntities(t *testing.T) {
	x, err := New()
	if err != nil {
		t.Fatal(err)
	}

	l, err := x.DiscoverEntityIDs()
	if err != nil {
		t.Error(err)
	}

	// At this point there are no entities, so the length should
	// be 0.
	if len(l) != 0 {
		t.Error("DiscoverEntityIDs made up an entity!")
	}

	// We'll save an entity that just has the ID set.  This isn't
	// very realistic, but its the minimum data needed to put a
	// file on disk.
	if err := x.SaveEntity(&pb.Entity{ID: proto.String("foo")}); err != nil {
		t.Error(err)
	}

	// Rerun discovery.
	l, err = x.DiscoverEntityIDs()
	if err != nil {
		t.Error(err)
	}

	// Now there should be one file on disk, and the ID that we've
	// discovered should be 'foo'
	if len(l) != 1 {
		t.Error("DiscoverEntityIDs failed to discover any entities!")
	}
	if l[0] != "foo" {
		t.Errorf("DiscoverEntityIDs discovered the wrong name: '%s'", l[0])
	}
}

func TestSaveLoadDeleteEntity(t *testing.T) {
	x, err := New()
	if err != nil {
		t.Fatal(err)
	}

	e := &pb.Entity{ID: proto.String("foo")}

	// Write an entity to disk
	if err := x.SaveEntity(e); err != nil {
		t.Error(err)
	}

	// Load  it back, it  should still be  the same, but  we check
	// this to be sure.
	ne, err := x.LoadEntity("foo")
	if err != nil {
		t.Error(err)
	}
	if !proto.Equal(e, ne) {
		t.Errorf("Loaded entity and original are not equivalent! '%v', '%v'", e, ne)
	}

	// Delete the entity and confirm that loading it returns an
	// error.
	if err := x.DeleteEntity("foo"); err != nil {
		t.Error(err)
	}
	if _, err := x.LoadEntity("foo"); err != db.ErrUnknownEntity {
		t.Error(err)
	}
}

func TestDiscoverGroups(t *testing.T) {
	x, err := New()
	if err != nil {
		t.Fatal(err)
	}

	l, err := x.DiscoverGroupNames()
	if err != nil {
		t.Error(err)
	}

	// At this point there are no groups, so the length should
	// be 0.
	if len(l) != 0 {
		t.Error("DiscoverGroupNames made up an entity!")
	}

	// We'll save an entity that just has the ID set.  This isn't
	// very realistic, but its the minimum data needed to put a
	// file on disk.
	if err := x.SaveGroup(&pb.Group{Name: proto.String("foo")}); err != nil {
		t.Error(err)
	}

	// Rerun discovery.
	l, err = x.DiscoverGroupNames()
	if err != nil {
		t.Error(err)
	}

	// Now there should be one file on disk, and the ID that we've
	// discovered should be 'foo'
	if len(l) != 1 {
		t.Error("DiscoverGroupNames failed to discover any groups!")
	}
	if l[0] != "foo" {
		t.Errorf("DiscoverGroupNames discovered the wrong name: '%s'", l[0])
	}
}

func TestGroupSaveLoadDelete(t *testing.T) {
	x, err := New()
	if err != nil {
		t.Fatal(err)
	}

	g := &pb.Group{Name: proto.String("foo")}

	// Write an entity to disk
	if err := x.SaveGroup(g); err != nil {
		t.Error(err)
	}

	// Load  it back, it  should still be  the same, but  we check
	// this to be sure.
	ng, err := x.LoadGroup("foo")
	if err != nil {
		t.Error(err)
	}
	if !proto.Equal(g, ng) {
		t.Errorf("Loaded group and original are not equivalent! '%v', '%v'", g, ng)
	}

	// Delete the group and confirm that loading it returns an
	// error.
	if err := x.DeleteGroup("foo"); err != nil {
		t.Error(err)
	}
	if _, err := x.LoadGroup("foo"); err != db.ErrUnknownGroup {
		t.Error(err)
	}
}

func TestDeleteEntityUnknown(t *testing.T) {
	x, err := New()
	if err != nil {
		t.Fatal(err)
	}

	if err := x.DeleteEntity("unknown-entity"); err != db.ErrUnknownEntity {
		t.Error(err)
	}
}

func TestDeleteGroupUnknown(t *testing.T) {
	x, err := New()
	if err != nil {
		t.Fatal(err)
	}

	if err := x.DeleteGroup("unknown-group"); err != db.ErrUnknownGroup {
		t.Error(err)
	}
}

func TestHealthCheck(t *testing.T) {
	x, err := New()
	if err != nil {
		t.Fatal(err)
	}

	// Fish out the concrete type to call non-interface methods.
	rx, ok := x.(*MemDB)
	if !ok {
		t.Fatal("Type assertion failed, bad type!")
	}

	if r := rx.healthCheck(); r.OK != true {
		t.Error("hard coded health check somehow changed")
	}
}
