package test

import (
	"context"
	"fmt"
	"time"

	"xiam.li/protonats/go/protonats"
)

type testImplementation struct {
	id string
}

func (t *testImplementation) SetTestServiceId(id string) {
	t.id = id
}

func (t *testImplementation) NormalTestTest(_ context.Context, req *Test) (*Test, error) {
	return &Test{Test: fmt.Sprintf("server replying to %s from %s", req.Test, t.id)}, nil
}

func (t *testImplementation) NormalEmptyTest(context.Context) (*Test, error) {
	return &Test{Test: fmt.Sprintf("server replying to empty from %s", t.id)}, nil
}

func (t *testImplementation) NormalTestEmpty(context.Context, *Test) error {
	return nil
}

func (t *testImplementation) NormalEmptyEmpty(context.Context) error {
	return nil
}

func (t *testImplementation) ErrServiceError(context.Context, *Test) (*Test, error) {
	return nil, protonats.ServiceError{Code: "1337", Description: "This is a service error"}
}

func (t *testImplementation) ErrServerError(context.Context, *Test) (*Test, error) {
	return nil, protonats.NewServerErr("1337", "This is a server error")
}

func (t *testImplementation) ErrServiceErrorBroadcast(context.Context, *Test) (*Test, error) {
	return nil, protonats.ServiceError{Code: "1337", Description: "This is a service error from " + t.id}
}

func (t *testImplementation) ErrServerErrorBroadcast(context.Context, *Test) (*Test, error) {
	return nil, protonats.NewServerErr("1337", "This is a server error from "+t.id)
}

func (t *testImplementation) NormalBroadcastTestTest(_ context.Context, req *Test) (*Test, error) {
	return &Test{Test: fmt.Sprintf("server replying to %s from %s", req.Test, t.id)}, nil
}

func (t *testImplementation) NormalBroadcastEmptyTest(context.Context) (*Test, error) {
	return &Test{Test: fmt.Sprintf("server replying to empty from %s", t.id)}, nil
}

func (t *testImplementation) NormalBroadcastTestEmpty(context.Context, *Test) error {
	return nil
}

func (t *testImplementation) NormalBroadcastEmptyEmpty(context.Context) error {
	return nil
}

func (t *testImplementation) LeaderOnlyTestTest(_ context.Context, req *Test) (*Test, error) {
	return &Test{Test: fmt.Sprintf("leader replying to %s from %s", req.Test, t.id)}, nil
}

func (t *testImplementation) LeaderOnlyEmptyTest(context.Context) (*Test, error) {
	return &Test{Test: fmt.Sprintf("leader replying to empty from %s", t.id)}, nil
}

func (t *testImplementation) LeaderOnlyTestEmpty(context.Context, *Test) error {
	return nil
}

func (t *testImplementation) LeaderOnlyEmptyEmpty(context.Context) error {
	return nil
}

func (t *testImplementation) LeaderOnlyBroadcastTestTest(_ context.Context, req *Test) (*Test, error) {
	return &Test{Test: fmt.Sprintf("leader replying to %s from %s", req.Test, t.id)}, nil
}

func (t *testImplementation) LeaderOnlyBroadcastEmptyTest(context.Context) (*Test, error) {
	return &Test{Test: fmt.Sprintf("leader replying to empty from %s", t.id)}, nil
}

func (t *testImplementation) LeaderOnlyBroadcastTestEmpty(context.Context, *Test) error {
	return nil
}

func (t *testImplementation) LeaderOnlyBroadcastEmptyEmpty(context.Context) error {
	return nil
}

func (t *testImplementation) FollowerOnlyTestTest(_ context.Context, req *Test) (*Test, error) {
	return &Test{Test: fmt.Sprintf("follower replying to %s from %s", req.Test, t.id)}, nil
}

func (t *testImplementation) FollowerOnlyEmptyTest(context.Context) (*Test, error) {
	return &Test{Test: fmt.Sprintf("follower replying to empty from %s", t.id)}, nil
}

func (t *testImplementation) FollowerOnlyTestEmpty(context.Context, *Test) error {
	return nil
}

func (t *testImplementation) FollowerOnlyEmptyEmpty(context.Context) error {
	return nil
}

func (t *testImplementation) FollowerOnlyBroadcastTestTest(_ context.Context, req *Test) (*Test, error) {
	return &Test{Test: fmt.Sprintf("follower replying to %s from %s", req.Test, t.id)}, nil
}

func (t *testImplementation) FollowerOnlyBroadcastEmptyTest(context.Context) (*Test, error) {
	return &Test{Test: fmt.Sprintf("follower replying to empty from %s", t.id)}, nil
}

func (t *testImplementation) FollowerOnlyBroadcastTestEmpty(context.Context, *Test) error {
	return nil
}

func (t *testImplementation) FollowerOnlyBroadcastEmptyEmpty(context.Context) error {
	return nil
}

func (t *testImplementation) ThreeSecondDelay(context.Context) error {
	time.Sleep(3 * time.Second)
	return nil
}

// Interface guard
var _ TestServiceNATSServer = (*testImplementation)(nil)
