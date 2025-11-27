package test

import (
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

func (t *testImplementation) NormalTestTest(req *Test) (*Test, error) {
	return &Test{Test: fmt.Sprintf("server replying to %s from %s", req.Test, t.id)}, nil
}

func (t *testImplementation) NormalEmptyTest() (*Test, error) {
	return &Test{Test: fmt.Sprintf("server replying to empty from %s", t.id)}, nil
}

func (t *testImplementation) NormalTestEmpty(*Test) error {
	return nil
}

func (t *testImplementation) NormalEmptyEmpty() error {
	return nil
}

func (t *testImplementation) ErrServiceError(*Test) (*Test, error) {
	return nil, protonats.ServiceError{Code: "1337", Description: "This is a service error"}
}

func (t *testImplementation) ErrServerError(*Test) (*Test, error) {
	return nil, protonats.NewServerErr("1337", "This is a server error")
}

func (t *testImplementation) ErrServiceErrorBroadcast(*Test) (*Test, error) {
	return nil, protonats.ServiceError{Code: "1337", Description: "This is a service error from " + t.id}
}

func (t *testImplementation) ErrServerErrorBroadcast(*Test) (*Test, error) {
	return nil, protonats.NewServerErr("1337", "This is a server error from "+t.id)
}

func (t *testImplementation) NormalBroadcastTestTest(req *Test) (*Test, error) {
	return &Test{Test: fmt.Sprintf("server replying to %s from %s", req.Test, t.id)}, nil
}

func (t *testImplementation) NormalBroadcastEmptyTest() (*Test, error) {
	return &Test{Test: fmt.Sprintf("server replying to empty from %s", t.id)}, nil
}

func (t *testImplementation) NormalBroadcastTestEmpty(*Test) error {
	return nil
}

func (t *testImplementation) NormalBroadcastEmptyEmpty() error {
	return nil
}

func (t *testImplementation) LeaderOnlyTestTest(req *Test) (*Test, error) {
	return &Test{Test: fmt.Sprintf("leader replying to %s from %s", req.Test, t.id)}, nil
}

func (t *testImplementation) LeaderOnlyEmptyTest() (*Test, error) {
	return &Test{Test: fmt.Sprintf("leader replying to empty from %s", t.id)}, nil
}

func (t *testImplementation) LeaderOnlyTestEmpty(*Test) error {
	return nil
}

func (t *testImplementation) LeaderOnlyEmptyEmpty() error {
	return nil
}

func (t *testImplementation) LeaderOnlyBroadcastTestTest(req *Test) (*Test, error) {
	return &Test{Test: fmt.Sprintf("leader replying to %s from %s", req.Test, t.id)}, nil
}

func (t *testImplementation) LeaderOnlyBroadcastEmptyTest() (*Test, error) {
	return &Test{Test: fmt.Sprintf("leader replying to empty from %s", t.id)}, nil
}

func (t *testImplementation) LeaderOnlyBroadcastTestEmpty(*Test) error {
	return nil
}

func (t *testImplementation) LeaderOnlyBroadcastEmptyEmpty() error {
	return nil
}

func (t *testImplementation) FollowerOnlyTestTest(req *Test) (*Test, error) {
	return &Test{Test: fmt.Sprintf("follower replying to %s from %s", req.Test, t.id)}, nil
}

func (t *testImplementation) FollowerOnlyEmptyTest() (*Test, error) {
	return &Test{Test: fmt.Sprintf("follower replying to empty from %s", t.id)}, nil
}

func (t *testImplementation) FollowerOnlyTestEmpty(*Test) error {
	return nil
}

func (t *testImplementation) FollowerOnlyEmptyEmpty() error {
	return nil
}

func (t *testImplementation) FollowerOnlyBroadcastTestTest(req *Test) (*Test, error) {
	return &Test{Test: fmt.Sprintf("follower replying to %s from %s", req.Test, t.id)}, nil
}

func (t *testImplementation) FollowerOnlyBroadcastEmptyTest() (*Test, error) {
	return &Test{Test: fmt.Sprintf("follower replying to empty from %s", t.id)}, nil
}

func (t *testImplementation) FollowerOnlyBroadcastTestEmpty(*Test) error {
	return nil
}

func (t *testImplementation) FollowerOnlyBroadcastEmptyEmpty() error {
	return nil
}

func (t *testImplementation) ThreeSecondDelay() error {
	time.Sleep(3 * time.Second)
	return nil
}

// Interface guard
var _ TestServiceNATSServer = (*testImplementation)(nil)
