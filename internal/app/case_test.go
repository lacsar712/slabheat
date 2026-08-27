package app_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/lacsar712/slabheat/internal/app"
	"github.com/lacsar712/slabheat/internal/clock"
	"github.com/lacsar712/slabheat/internal/config"
	"github.com/lacsar712/slabheat/internal/model"
)

func TestCase(t *testing.T) {
	clk := clock.NewManual(time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC))
	cfg := config.Default("BD-1")
	a, err := app.BootstrapWithClock(cfg, clk)
	if err != nil {
		t.Fatal(err)
	}
	if err := a.Trip(context.Background(), "test"); err != nil {
		t.Fatal(err)
	}
	err = a.DescalcAfterShutdown(context.Background(), 100)
	if err == nil {
		t.Fatal("expected blowdown limit error")
	}
	if !errors.Is(err, model.ErrDescalcLimit) {
		t.Fatalf("expected ErrDescalcLimit, got %v", err)
	}
}
