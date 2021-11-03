package netem

import (
	"strings"
	"testing"
	"time"
)

var iface = "wlp5s0"

func TestNetemDelay(t *testing.T) {
	opt := Option{
		NetworkIface: iface,
	}

	netem, err := New(opt)
	if err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		Name          string
		Delay         time.Duration
		ExpectedOuput string
	}{
		{
			Name:          "Test 100 milliseconds delay",
			Delay:         100 * time.Millisecond,
			ExpectedOuput: "delay 100.0ms",
		},
		{
			Name:          "Test 2000 microseconds delay",
			Delay:         2000 * time.Microsecond,
			ExpectedOuput: "delay 2.0ms",
		},
		{
			Name:          "Test 2 seconds delay",
			Delay:         2 * time.Second,
			ExpectedOuput: "delay 2.0s",
		},
		{
			Name:          "Test 5 seconds delay inputted with ms",
			Delay:         5000 * time.Millisecond,
			ExpectedOuput: "delay 5.0s",
		},
		{
			Name:          "Test 1000 millisecond delay inputted with S",
			Delay:         1 * time.Second,
			ExpectedOuput: "delay 1.0s",
		},
	}

	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			if err := netem.AddDelay(tt.Delay); err != nil {
				t.Fatal(err)
			}

			// check output from tc
			out, err := netem.Show()
			if err != nil {
				t.Fatal(err)
			}

			output := strings.Join(out, "")

			if !strings.Contains(output, tt.ExpectedOuput) {
				t.Errorf("output from tc not equal with %v \n tc output: %v", tt.ExpectedOuput, out)
			}

			// cleanup
			if err := netem.DeleteDelay(tt.Delay); err != nil {
				t.Fatal(err)
			}
		})
	}
}

func TestNetemDuplicateAdd(t *testing.T) {
	opt := Option{
		NetworkIface: iface,
	}

	netem, err := New(opt)
	if err != nil {
		t.Fatal(err)
	}

	expectedError := "Exclusivity flag on, cannot modify"

	if err := netem.AddDelay(1 * time.Second); err != nil {
		t.Fatal(err)
	}

	if err := netem.AddDelay(1 * time.Second); err != nil {
		if !strings.Contains(err.Error(), expectedError) {
			t.Errorf("got error %v, want error %v", err.Error(), expectedError)
		}
	}

	// cleanup
	if err := netem.DeleteDelay(1 * time.Second); err != nil {
		t.Fatal(err)
	}
}

func TestNetemChangeRules(t *testing.T) {
	opt := Option{
		NetworkIface: iface,
	}

	netem, err := New(opt)
	if err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		Name          string
		Delay         time.Duration
		ModifyDelay   time.Duration
		ExpectedOuput string
	}{
		{
			Name:          "Test Modify 100ms delay to 200ms",
			Delay:         100 * time.Millisecond,
			ModifyDelay:   200 * time.Millisecond,
			ExpectedOuput: "delay 200.0ms",
		},
		{
			Name:          "Test Modify 2000ms delay to 2s",
			Delay:         2000 * time.Millisecond,
			ModifyDelay:   2 * time.Second,
			ExpectedOuput: "delay 2.0s",
		},
		{
			Name:          "Test Modify 2s delay to 500ms",
			Delay:         2 * time.Second,
			ModifyDelay:   500 * time.Millisecond,
			ExpectedOuput: "delay 500.0ms",
		},
		{
			Name:          "Test modify 5s delay to 1s",
			Delay:         5000 * time.Millisecond,
			ModifyDelay:   1000 * time.Millisecond,
			ExpectedOuput: "delay 1.0s",
		},
	}

	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			if err := netem.AddDelay(tt.Delay); err != nil {
				t.Fatal(err)
			}

			// modify delay
			if err := netem.ChangeDelay(tt.ModifyDelay); err != nil {
				t.Fatal(err)
			}

			// check output from tc
			out, err := netem.Show()
			if err != nil {
				t.Fatal(err)
			}

			output := strings.Join(out, "")

			if !strings.Contains(output, tt.ExpectedOuput) {
				t.Errorf("output from tc not equal with %v \n tc output: %v", tt.ExpectedOuput, out)
			}

			// cleanup
			if err := netem.DeleteDelay(tt.Delay); err != nil {
				t.Fatal(err)
			}
		})
	}
}
