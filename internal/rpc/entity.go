package rpc

import (
	"context"
	"fmt"
	"log"

	"github.com/golang/protobuf/proto"

	pb "github.com/NetAuth/NetAuth/pkg/proto"
)

func (s *NetAuthServer) NewEntity(ctx context.Context, r *pb.ModEntityRequest) (*pb.SimpleResult, error) {
	client := r.GetInfo()
	e := r.GetEntity()
	t := r.GetAuthToken()

	c, err := s.Token.Validate(t)
	if err != nil {
		return &pb.SimpleResult{Msg: proto.String("Authentication Failure")}, nil
	}

	// Verify the correct capability is present in the token.
	if !c.HasCapability("CREATE_ENTITY") {
		return &pb.SimpleResult{Msg: proto.String("Requestor not qualified"), Success: proto.Bool(false)}, nil
	}

	if err := s.Tree.NewEntity(e.GetID(), e.GetUidNumber(), e.GetSecret()); err != nil {
		return &pb.SimpleResult{Success: proto.Bool(false), Msg: proto.String(fmt.Sprintf("%s", err))}, nil
	}

	log.Printf("New entity '%s' created by '%s' (%s@%s)",
		e.GetID(),
		c.EntityID,
		client.GetService(),
		client.GetID())

	return &pb.SimpleResult{
		Msg:     proto.String("New entity created successfully"),
		Success: proto.Bool(true),
	}, nil
}

func (s *NetAuthServer) RemoveEntity(ctx context.Context, r *pb.ModEntityRequest) (*pb.SimpleResult, error) {
	client := r.GetInfo()
	e := r.GetEntity()
	t := r.GetAuthToken()

	c, err := s.Token.Validate(t)
	if err != nil {
		return &pb.SimpleResult{Msg: proto.String("Authentication Failure")}, nil
	}

	// Verify the correct capability is present in the token.
	if !c.HasCapability("DESTROY_ENTITY") {
		return &pb.SimpleResult{Msg: proto.String("Requestor not qualified"), Success: proto.Bool(false)}, nil
	}

	if err := s.Tree.DeleteEntityByID(e.GetID()); err != nil {
		return &pb.SimpleResult{Success: proto.Bool(false), Msg: proto.String(fmt.Sprintf("%s", err))}, nil
	}

	log.Printf("Entity '%s' removed by '%s' (%s@%s)",
		e.GetID(),
		c.EntityID,
		client.GetService(),
		client.GetID())

	return &pb.SimpleResult{
		Msg:     proto.String("Entity removed successfully"),
		Success: proto.Bool(true),
	}, nil
}

func (s *NetAuthServer) EntityInfo(ctx context.Context, r *pb.NetAuthRequest) (*pb.Entity, error) {
	client := r.GetInfo()
	e := r.GetEntity()

	log.Printf("Info requested on '%s' (%s@%s)",
		e.GetID(),
		client.GetService(),
		client.GetID())

	return s.Tree.GetEntity(e.GetID())
}

func (s *NetAuthServer) ModifyEntityMeta(ctx context.Context, r *pb.ModEntityRequest) (*pb.SimpleResult, error) {
	client := r.GetInfo()
	e := r.GetEntity()
	t := r.GetAuthToken()

	c, err := s.Token.Validate(t)
	if err != nil {
		return &pb.SimpleResult{Msg: proto.String("Authentication Failure")}, nil
	}

	// Verify the correct capability is present in the token.
	if !c.HasCapability("MODIFY_ENTITY_META") {
		return &pb.SimpleResult{Msg: proto.String("Requestor not qualified"), Success: proto.Bool(false)}, nil
	}

	if err := s.Tree.UpdateEntityMeta(e.GetID(), e.GetMeta()); err != nil {
		log.Printf("Metadata update error: %s", err)
		return nil, err
	}

	log.Printf("Metadata for '%s' by '%s' completed (%s@%s)",
		e.GetID(),
		c.EntityID,
		client.GetService(),
		client.GetID())

	return &pb.SimpleResult{Success: proto.Bool(true), Msg: proto.String("Metadata Updated")}, nil
}