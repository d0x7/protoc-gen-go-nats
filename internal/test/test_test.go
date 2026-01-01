package test

import (
	"context"
	"errors"
	"fmt"
	"regexp"
	"slices"
	"testing"
	"time"

	"github.com/nats-io/nats.go"
	"xiam.li/protonats/go/protonats"
)

func TestInfo(t *testing.T) {
	t.Parallel()
	instance := newNATS(t)
	t.Cleanup(instance.Stop)
	var ids []string
	for range 10 {
		id := NewTestServiceNATSServer(instance.Conn, new(testImplementation)).Info().ID
		ids = append(ids, id)
	}
	cli := NewTestServiceNATSClient(instance.Conn)

	t.Run("WithFinisher", func(t *testing.T) {
		now := time.Now()
		info, err := cli.Info(t.Context())
		dur := time.Since(now)
		if err != nil {
			t.Fatalf("Error calling method: %v", err)
		}
		if dur > 300*time.Millisecond {
			t.Fatalf("Info did not return within 500ms (%v)", dur)
		}
		if len(info) != 10 {
			t.Fatalf("Expected 10 responses, got %d", len(info))
		}
	})

	t.Run("WithoutFinisher", func(t *testing.T) {
		now := time.Now()
		info, err := cli.Info(t.Context(), protonats.WithoutFinisher())
		dur := time.Since(now)
		if err != nil {
			t.Fatalf("Error calling method: %v", err)
		}
		if dur <= 5000*time.Millisecond {
			t.Fatalf("Info returned under 4.5s (%v)", dur)
		}
		if len(info) != 10 {
			t.Fatalf("Expected 10 responses, got %d", len(info))
		}
	})
}

func TestNormal(t *testing.T) {
	t.Parallel()
	instance := newNATS(t)
	t.Cleanup(instance.Stop)
	var ids []string
	for range 3 {
		id := NewTestServiceNATSServer(instance.Conn, new(testImplementation)).Info().ID
		ids = append(ids, id)
	}
	cli := NewTestServiceNATSClient(instance.Conn)

	t.Run("TestTest", func(t *testing.T) {
		t.Parallel()
		resp, err := cli.NormalTestTest(t.Context(), &Test{Test: "Test Client"})
		if err != nil {
			t.Fatalf("Error calling method: %v", err)
		}
		re := regexp.MustCompile(`^server replying to Test Client from ([a-zA-Z0-9]+)$`)
		matches := re.FindStringSubmatch(resp.Test)
		if len(matches) != 2 {
			t.Fatalf("Response format doesn't match: %v", resp.Test)
		}
		if !slices.Contains(ids, matches[1]) {
			t.Fatalf("Server ID not found in the list of IDs: %v", matches[1])
		}
	})

	t.Run("EmptyTest", func(t *testing.T) {
		t.Parallel()
		resp, err := cli.NormalEmptyTest(t.Context())
		if err != nil {
			t.Fatalf("Error calling method: %v", err)
		}
		re := regexp.MustCompile(`^server replying to empty from ([a-zA-Z0-9]+)$`)
		matches := re.FindStringSubmatch(resp.Test)
		if len(matches) != 2 {
			t.Fatalf("Response format doesn't match: %v", resp.Test)
		}
		if !slices.Contains(ids, matches[1]) {
			t.Fatalf("Server ID not found in the list of IDs: %v", matches[1])
		}
	})

	t.Run("TestEmpty", func(t *testing.T) {
		t.Parallel()
		err := cli.NormalTestEmpty(t.Context(), &Test{Test: "Test Client"})
		if err != nil {
			t.Fatalf("Error calling method: %v", err)
		}
	})

	t.Run("EmptyEmpty", func(t *testing.T) {
		t.Parallel()
		err := cli.NormalEmptyEmpty(t.Context())
		if err != nil {
			t.Fatalf("Error calling method: %v", err)
		}
	})
}

func TestNormalBroadcast(t *testing.T) {
	t.Parallel()
	instance := newNATS(t)
	t.Cleanup(instance.Stop)
	var ids []string
	for range 3 {
		id := NewTestServiceNATSServer(instance.Conn, new(testImplementation)).Info().ID
		ids = append(ids, id)
	}
	cli := NewTestServiceNATSClient(instance.Conn)

	t.Run("TestTest", func(t *testing.T) {
		t.Parallel()
		resp, srvErrs, err := cli.NormalBroadcastTestTest(t.Context(), &Test{Test: "Test Client"})
		if err != nil {
			t.Fatalf("Error calling method: %v", err)
		}
		if len(srvErrs) != 0 {
			t.Fatalf("Unexpected server errors: %v", srvErrs)
		}
		if len(resp) != 3 {
			t.Fatalf("Expected 3 responses, got %d: %v", len(resp), resp)
		}
		for _, r := range resp {
			re := regexp.MustCompile(`^server replying to Test Client from ([a-zA-Z0-9]+)$`)
			matches := re.FindStringSubmatch(r.Test)
			if len(matches) != 2 {
				t.Fatalf("Response format doesn't match: %v", r.Test)
			}
			if !slices.Contains(ids, matches[1]) {
				t.Fatalf("Server ID not found in the list of IDs: %v", matches[1])
			}
		}
	})

	t.Run("EmptyTest", func(t *testing.T) {
		t.Parallel()
		resp, srvErrs, err := cli.NormalBroadcastEmptyTest(t.Context())
		if err != nil {
			t.Fatalf("Error calling method: %v", err)
		}
		if len(srvErrs) != 0 {
			t.Fatalf("Unexpected server errors: %v", srvErrs)
		}
		if len(resp) != 3 {
			t.Fatalf("Expected 3 responses, got %d: %v", len(resp), resp)
		}
		for _, r := range resp {
			re := regexp.MustCompile(`^server replying to empty from ([a-zA-Z0-9]+)$`)
			matches := re.FindStringSubmatch(r.Test)
			if len(matches) != 2 {
				t.Fatalf("Response format doesn't match: %v", r.Test)
			}
			if !slices.Contains(ids, matches[1]) {
				t.Fatalf("Server ID not found in the list of IDs: %v", matches[1])
			}
		}
	})

	t.Run("TestEmpty", func(t *testing.T) {
		t.Parallel()
		srvErrs, err := cli.NormalBroadcastTestEmpty(t.Context(), &Test{Test: "Test Client"})
		if err != nil {
			t.Fatalf("Error calling method: %v", err)
		}
		if len(srvErrs) != 0 {
			t.Fatalf("Unexpected server errors: %v", srvErrs)
		}
	})

	t.Run("EmptyEmpty", func(t *testing.T) {
		t.Parallel()
		srvErrs, err := cli.NormalBroadcastEmptyEmpty(t.Context())
		if err != nil {
			t.Fatalf("Error calling method: %v", err)
		}
		if len(srvErrs) != 0 {
			t.Fatalf("Unexpected server errors: %v", srvErrs)
		}
	})
}

func TestErr(t *testing.T) {
	t.Parallel()
	instance := newNATS(t)
	t.Cleanup(instance.Stop)
	NewTestServiceNATSServer(instance.Conn, new(testImplementation))
	cli := NewTestServiceNATSClient(instance.Conn)

	t.Run("ServiceError", func(t *testing.T) {
		t.Parallel()
		resp, err := cli.ErrServiceError(t.Context(), &Test{Test: "Test Client"})
		if err == nil {
			t.Fatalf("Expected error, got nil")
		}
		if resp != nil {
			t.Fatalf("Expected nil response, got: %v", resp)
		}
		if !protonats.IsServiceError(err) {
			t.Fatalf("Expected service error, got: %v", err)
		}
	})

	t.Run("ServerError", func(t *testing.T) {
		t.Parallel()
		resp, err := cli.ErrServerError(t.Context(), &Test{Test: "Test Client"})
		if err == nil {
			t.Fatalf("Expected error, got nil")
		}
		if !protonats.IsServiceError(err) {
			t.Fatalf("Expected service error, got: %v", err)
		}
		if resp != nil {
			t.Fatalf("Expected nil response, got: %v", resp)
		}
	})
}

func TestErrBroadcast(t *testing.T) {
	t.Parallel()
	instance := newNATS(t)
	t.Cleanup(instance.Stop)
	var ids []string
	for range 3 {
		id := NewTestServiceNATSServer(instance.Conn, new(testImplementation)).Info().ID
		ids = append(ids, id)
	}
	cli := NewTestServiceNATSClient(instance.Conn)

	t.Run("ServiceError", func(t *testing.T) {
		t.Parallel()
		resp, srvErrs, err := cli.ErrServiceErrorBroadcast(t.Context(), &Test{Test: "Test Client"})
		if err != nil {
			t.Fatalf("Error calling method: %v", err)
		}
		if len(resp) != 0 {
			t.Fatalf("Unexpected responses: %v", resp)
		}
		if len(srvErrs) != 3 {
			t.Fatalf("Expected 3 service errors, got %d: %v", len(srvErrs), srvErrs)
		}
		for _, e := range srvErrs {
			if !protonats.IsServiceError(e) {
				t.Fatalf("Expected service error, got: %v", e)
			}
		}
	})

	t.Run("ServerError", func(t *testing.T) {
		t.Parallel()
		resp, srvErrs, err := cli.ErrServerErrorBroadcast(t.Context(), &Test{Test: "Test Client"})
		if err != nil {
			t.Fatalf("Error calling method: %v", err)
		}
		if len(resp) != 0 {
			t.Fatalf("Unexpected responses: %v", resp)
		}
		if len(srvErrs) != 3 {
			t.Fatalf("Expected 3 service errors, got %d: %v", len(srvErrs), srvErrs)
		}
		for _, e := range srvErrs {
			if !protonats.IsServiceError(e) {
				t.Fatalf("Expected service error, got: %v", e)
			}
		}
	})
}

func TestExtraSubject(t *testing.T) {
	t.Parallel()
	instance := newNATS(t)
	t.Cleanup(instance.Stop)
	var ids []string
	for i := range 3 {
		id := fmt.Sprintf("instance%02d", i)
		impl := new(testImplementation)
		_ = NewTestServiceNATSServer(instance.Conn, impl, protonats.WithExtraSubjectSrv(id))
		ids = append(ids, impl.id)
		impl.id = impl.id + " aka " + id
	}
	cli := NewTestServiceNATSClient(instance.Conn)

	for i := range 3 {
		id := fmt.Sprintf("instance%02d", i)
		t.Run("EmptyTest/"+id, func(t *testing.T) {
			t.Parallel()
			resp, srvErrs, err := cli.NormalBroadcastEmptyTest(t.Context(), protonats.WithExtraSubject(id))
			if err != nil {
				t.Fatalf("Error calling method: %v", err)
			}
			if len(srvErrs) != 0 {
				t.Fatalf("Unexpected server errors: %v", srvErrs)
			}
			if len(resp) != 1 {
				t.Fatalf("Expected one response, got %d: %v", len(resp), resp)
			}
			for _, r := range resp {
				re := regexp.MustCompile(`^server replying to empty from ([a-zA-Z0-9]+) aka ([a-zA-Z0-9]+)$`)
				matches := re.FindStringSubmatch(r.Test)
				if len(matches) != 3 {
					t.Fatalf("Response format doesn't match: %v", r.Test)
				}
				if !slices.Contains(ids, matches[1]) {
					t.Fatalf("Server ID not found in the list of follower IDs: %v", matches[1])
				}
				if matches[2] != id {
					t.Fatalf("Extra subject doesn't match: %v", matches[2])
				}
			}
		})
	}
}

func TestContext(t *testing.T) {
	t.Parallel()
	instance := newNATS(t)
	t.Cleanup(instance.Stop)
	_ = NewTestServiceNATSServer(instance.Conn, new(testImplementation))
	cli := NewTestServiceNATSClient(instance.Conn)

	t.Run("WithTimeout", func(t *testing.T) {
		t.Parallel()
		start := time.Now()
		err := cli.ThreeSecondDelay(t.Context(), protonats.WithTimeout(1*time.Second))
		if err == nil {
			t.Fatalf("Expected error, got nil")
		}
		if !errors.Is(err, context.DeadlineExceeded) {
			t.Fatalf("Expected deadline exceeded error, got: %v", err)
		}
		elapsed := time.Since(start)

		if elapsed < time.Second || elapsed > 1100*time.Millisecond {
			t.Fatalf("expected 1s timeout, got %v", elapsed)
		}
	})

	t.Run("ContextCancel", func(t *testing.T) {
		t.Parallel()
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		go func() {
			time.Sleep(1 * time.Second)
			cancel()
		}()
		err := cli.ThreeSecondDelay(ctx)
		if err == nil {
			t.Fatalf("Expected error, got nil")
		}
		if !errors.Is(err, context.Canceled) {
			t.Fatalf("Expected context cancelled error, got: %v", err)
		}
	})

	t.Run("ContextDeadline", func(t *testing.T) {
		t.Parallel()
		start := time.Now()
		ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
		defer cancel()
		err := cli.ThreeSecondDelay(ctx)
		if err == nil {
			t.Fatalf("Expected error, got nil")
		}
		if !errors.Is(err, context.DeadlineExceeded) {
			t.Fatalf("Expected deadline exceeded error, got: %v", err)
		}
		elapsed := time.Since(start)

		if elapsed < time.Second || elapsed > 1100*time.Millisecond {
			t.Fatalf("expected 1s timeout, got %v", elapsed)
		}
	})
}

func TestNoResponder(t *testing.T) {
	t.Parallel()
	instance := newNATS(t)
	t.Cleanup(instance.Stop)
	cli := NewTestServiceNATSClient(instance.Conn)

	t.Run("NoResponder", func(t *testing.T) {
		t.Parallel()
		err := cli.NormalEmptyEmpty(t.Context())
		if err == nil {
			t.Fatalf("Expected error, got nil")
		}
		if !errors.Is(err, nats.ErrNoResponders) {
			t.Fatalf("Expected no responder error, got: %v", err)
		}
	})
}
