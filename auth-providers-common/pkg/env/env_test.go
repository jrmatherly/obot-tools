package env

import (
	"os"
	"strings"
	"testing"
	"time"
)

func TestLoadEnvForStruct_Int(t *testing.T) {
	tests := []struct {
		name          string
		envVars       map[string]string
		setup         func(*testing.T)
		cleanup       func(*testing.T)
		expectError   bool
		errorContains string
		validate      func(*testing.T, *testIntStruct)
	}{
		{
			name: "positive int",
			setup: func(t *testing.T) {
				os.Setenv("TEST_INT", "42")
			},
			cleanup: func(t *testing.T) {
				os.Unsetenv("TEST_INT")
			},
			expectError: false,
			validate: func(t *testing.T, s *testIntStruct) {
				if s.Value != 42 {
					t.Errorf("Value = %d, want 42", s.Value)
				}
			},
		},
		{
			name: "negative int",
			setup: func(t *testing.T) {
				os.Setenv("TEST_INT", "-123")
			},
			cleanup: func(t *testing.T) {
				os.Unsetenv("TEST_INT")
			},
			expectError: false,
			validate: func(t *testing.T, s *testIntStruct) {
				if s.Value != -123 {
					t.Errorf("Value = %d, want -123", s.Value)
				}
			},
		},
		{
			name: "zero int",
			setup: func(t *testing.T) {
				os.Setenv("TEST_INT", "0")
			},
			cleanup: func(t *testing.T) {
				os.Unsetenv("TEST_INT")
			},
			expectError: false,
			validate: func(t *testing.T, s *testIntStruct) {
				if s.Value != 0 {
					t.Errorf("Value = %d, want 0", s.Value)
				}
			},
		},
		{
			name: "invalid int",
			setup: func(t *testing.T) {
				os.Setenv("TEST_INT", "not-a-number")
			},
			cleanup: func(t *testing.T) {
				os.Unsetenv("TEST_INT")
			},
			expectError:   true,
			errorContains: "invalid integer value",
		},
		{
			name: "int with default value",
			setup: func(t *testing.T) {
				// Don't set TEST_INT_DEFAULT
			},
			cleanup:     func(t *testing.T) {},
			expectError: false,
			validate: func(t *testing.T, s *testIntStruct) {
				if s.DefaultValue != 99 {
					t.Errorf("DefaultValue = %d, want 99", s.DefaultValue)
				}
			},
		},
		{
			name: "optional int not set",
			setup: func(t *testing.T) {
				// Don't set TEST_INT_OPT
			},
			cleanup:     func(t *testing.T) {},
			expectError: false,
			validate: func(t *testing.T, s *testIntStruct) {
				if s.OptionalValue != 0 {
					t.Errorf("OptionalValue = %d, want 0", s.OptionalValue)
				}
			},
		},
		{
			name: "int from second env var",
			setup: func(t *testing.T) {
				os.Setenv("TEST_INT_ALT", "777")
			},
			cleanup: func(t *testing.T) {
				os.Unsetenv("TEST_INT_ALT")
			},
			expectError: false,
			validate: func(t *testing.T, s *testIntStruct) {
				if s.MultiVar != 777 {
					t.Errorf("MultiVar = %d, want 777", s.MultiVar)
				}
			},
		},
		{
			name: "int from first env var takes precedence",
			setup: func(t *testing.T) {
				os.Setenv("TEST_INT_PRIMARY", "111")
				os.Setenv("TEST_INT_ALT", "222")
			},
			cleanup: func(t *testing.T) {
				os.Unsetenv("TEST_INT_PRIMARY")
				os.Unsetenv("TEST_INT_ALT")
			},
			expectError: false,
			validate: func(t *testing.T, s *testIntStruct) {
				if s.MultiVar != 111 {
					t.Errorf("MultiVar = %d, want 111 (first env var should take precedence)", s.MultiVar)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.setup != nil {
				tt.setup(t)
			}
			if tt.cleanup != nil {
				defer tt.cleanup(t)
			}

			s := &testIntStruct{}
			err := LoadEnvForStruct(s)

			if tt.expectError {
				if err == nil {
					t.Errorf("LoadEnvForStruct() expected error, got nil")
				} else if tt.errorContains != "" && !strings.Contains(err.Error(), tt.errorContains) {
					t.Errorf("LoadEnvForStruct() error = %v, want error containing %q", err, tt.errorContains)
				}
			} else {
				if err != nil {
					t.Errorf("LoadEnvForStruct() unexpected error = %v", err)
				} else if tt.validate != nil {
					tt.validate(t, s)
				}
			}
		})
	}
}

func TestLoadEnvForStruct_IntTypes(t *testing.T) {
	tests := []struct {
		name        string
		envValue    string
		expectError bool
	}{
		{
			name:        "int8 max value",
			envValue:    "127",
			expectError: false,
		},
		{
			name:        "int8 min value",
			envValue:    "-128",
			expectError: false,
		},
		{
			name:        "int16 large value",
			envValue:    "32000",
			expectError: false,
		},
		{
			name:        "int32 large value",
			envValue:    "2000000000",
			expectError: false,
		},
		{
			name:        "int64 very large value",
			envValue:    "9000000000000000000",
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			os.Setenv("TEST_INT8", tt.envValue)
			os.Setenv("TEST_INT16", tt.envValue)
			os.Setenv("TEST_INT32", tt.envValue)
			os.Setenv("TEST_INT64", tt.envValue)
			defer func() {
				os.Unsetenv("TEST_INT8")
				os.Unsetenv("TEST_INT16")
				os.Unsetenv("TEST_INT32")
				os.Unsetenv("TEST_INT64")
			}()

			s := &testIntTypesStruct{}
			err := LoadEnvForStruct(s)

			if tt.expectError {
				if err == nil {
					t.Errorf("LoadEnvForStruct() expected error, got nil")
				}
			} else {
				if err != nil {
					t.Errorf("LoadEnvForStruct() unexpected error = %v", err)
				}
			}
		})
	}
}

func TestLoadEnvForStruct_Bool(t *testing.T) {
	tests := []struct {
		name          string
		envValue      string
		expectError   bool
		errorContains string
		expectedValue bool
	}{
		{
			name:          "true lowercase",
			envValue:      "true",
			expectError:   false,
			expectedValue: true,
		},
		{
			name:          "True capitalized",
			envValue:      "True",
			expectError:   false,
			expectedValue: true,
		},
		{
			name:          "TRUE uppercase",
			envValue:      "TRUE",
			expectError:   false,
			expectedValue: true,
		},
		{
			name:          "t single letter",
			envValue:      "t",
			expectError:   false,
			expectedValue: true,
		},
		{
			name:          "T uppercase single letter",
			envValue:      "T",
			expectError:   false,
			expectedValue: true,
		},
		{
			name:          "1 numeric",
			envValue:      "1",
			expectError:   false,
			expectedValue: true,
		},
		{
			name:          "false lowercase",
			envValue:      "false",
			expectError:   false,
			expectedValue: false,
		},
		{
			name:          "False capitalized",
			envValue:      "False",
			expectError:   false,
			expectedValue: false,
		},
		{
			name:          "FALSE uppercase",
			envValue:      "FALSE",
			expectError:   false,
			expectedValue: false,
		},
		{
			name:          "f single letter",
			envValue:      "f",
			expectError:   false,
			expectedValue: false,
		},
		{
			name:          "F uppercase single letter",
			envValue:      "F",
			expectError:   false,
			expectedValue: false,
		},
		{
			name:          "0 numeric",
			envValue:      "0",
			expectError:   false,
			expectedValue: false,
		},
		{
			name:          "invalid bool value",
			envValue:      "not-a-bool",
			expectError:   true,
			errorContains: "invalid boolean value",
		},
		{
			name:          "yes not accepted",
			envValue:      "yes",
			expectError:   true,
			errorContains: "invalid boolean value",
		},
		{
			name:          "no not accepted",
			envValue:      "no",
			expectError:   true,
			errorContains: "invalid boolean value",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			os.Setenv("TEST_BOOL", tt.envValue)
			defer os.Unsetenv("TEST_BOOL")

			s := &testBoolStruct{}
			err := LoadEnvForStruct(s)

			if tt.expectError {
				if err == nil {
					t.Errorf("LoadEnvForStruct() expected error, got nil")
				} else if !strings.Contains(err.Error(), tt.errorContains) {
					t.Errorf("LoadEnvForStruct() error = %v, want error containing %q", err, tt.errorContains)
				}
			} else {
				if err != nil {
					t.Errorf("LoadEnvForStruct() unexpected error = %v", err)
				} else if s.Value != tt.expectedValue {
					t.Errorf("Value = %v, want %v", s.Value, tt.expectedValue)
				}
			}
		})
	}
}

func TestLoadEnvForStruct_BoolOptionalAndDefault(t *testing.T) {
	tests := []struct {
		name          string
		setup         func(*testing.T)
		cleanup       func(*testing.T)
		expectError   bool
		validateValue func(*testing.T, *testBoolStruct)
	}{
		{
			name: "bool with default true",
			setup: func(t *testing.T) {
				// Don't set TEST_BOOL_DEFAULT
			},
			cleanup:     func(t *testing.T) {},
			expectError: false,
			validateValue: func(t *testing.T, s *testBoolStruct) {
				if !s.DefaultValue {
					t.Errorf("DefaultValue = %v, want true", s.DefaultValue)
				}
			},
		},
		{
			name: "bool with default false",
			setup: func(t *testing.T) {
				// Don't set TEST_BOOL_DEFAULT_FALSE
			},
			cleanup:     func(t *testing.T) {},
			expectError: false,
			validateValue: func(t *testing.T, s *testBoolStruct) {
				if s.DefaultFalse {
					t.Errorf("DefaultFalse = %v, want false", s.DefaultFalse)
				}
			},
		},
		{
			name: "optional bool not set",
			setup: func(t *testing.T) {
				// Don't set TEST_BOOL_OPT
			},
			cleanup:     func(t *testing.T) {},
			expectError: false,
			validateValue: func(t *testing.T, s *testBoolStruct) {
				if s.OptionalValue {
					t.Errorf("OptionalValue = %v, want false", s.OptionalValue)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.setup != nil {
				tt.setup(t)
			}
			if tt.cleanup != nil {
				defer tt.cleanup(t)
			}

			s := &testBoolStruct{}
			err := LoadEnvForStruct(s)

			if tt.expectError {
				if err == nil {
					t.Errorf("LoadEnvForStruct() expected error, got nil")
				}
			} else {
				if err != nil {
					t.Errorf("LoadEnvForStruct() unexpected error = %v", err)
				} else if tt.validateValue != nil {
					tt.validateValue(t, s)
				}
			}
		})
	}
}

func TestLoadEnvForStruct_Duration(t *testing.T) {
	tests := []struct {
		name             string
		envValue         string
		expectError      bool
		errorContains    string
		expectedDuration time.Duration
	}{
		{
			name:             "milliseconds",
			envValue:         "300ms",
			expectError:      false,
			expectedDuration: 300 * time.Millisecond,
		},
		{
			name:             "seconds",
			envValue:         "30s",
			expectError:      false,
			expectedDuration: 30 * time.Second,
		},
		{
			name:             "minutes",
			envValue:         "5m",
			expectError:      false,
			expectedDuration: 5 * time.Minute,
		},
		{
			name:             "hours",
			envValue:         "2h",
			expectError:      false,
			expectedDuration: 2 * time.Hour,
		},
		{
			name:             "combined hours and minutes",
			envValue:         "2h30m",
			expectError:      false,
			expectedDuration: 2*time.Hour + 30*time.Minute,
		},
		{
			name:             "combined hours minutes seconds",
			envValue:         "1h30m45s",
			expectError:      false,
			expectedDuration: 1*time.Hour + 30*time.Minute + 45*time.Second,
		},
		{
			name:             "fractional hours",
			envValue:         "1.5h",
			expectError:      false,
			expectedDuration: 90 * time.Minute,
		},
		{
			name:             "fractional seconds",
			envValue:         "2.5s",
			expectError:      false,
			expectedDuration: 2500 * time.Millisecond,
		},
		{
			name:             "microseconds",
			envValue:         "500us",
			expectError:      false,
			expectedDuration: 500 * time.Microsecond,
		},
		{
			name:             "nanoseconds",
			envValue:         "1000ns",
			expectError:      false,
			expectedDuration: 1000 * time.Nanosecond,
		},
		{
			name:             "zero duration",
			envValue:         "0s",
			expectError:      false,
			expectedDuration: 0,
		},
		{
			name:          "invalid duration - no unit",
			envValue:      "300",
			expectError:   true,
			errorContains: "invalid duration value",
		},
		{
			name:          "invalid duration - bad format",
			envValue:      "not-a-duration",
			expectError:   true,
			errorContains: "invalid duration value",
		},
		{
			name:             "zero duration with 0 value",
			envValue:         "0",
			expectError:      false,
			expectedDuration: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			os.Setenv("TEST_DURATION", tt.envValue)
			defer os.Unsetenv("TEST_DURATION")

			s := &testDurationStruct{}
			err := LoadEnvForStruct(s)

			if tt.expectError {
				if err == nil {
					t.Errorf("LoadEnvForStruct() expected error, got nil")
				} else if tt.errorContains != "" && !strings.Contains(err.Error(), tt.errorContains) {
					t.Errorf("LoadEnvForStruct() error = %v, want error containing %q", err, tt.errorContains)
				}
			} else {
				if err != nil {
					t.Errorf("LoadEnvForStruct() unexpected error = %v", err)
				} else if s.Value != tt.expectedDuration {
					t.Errorf("Value = %v, want %v", s.Value, tt.expectedDuration)
				}
			}
		})
	}
}

func TestLoadEnvForStruct_DurationOptionalAndDefault(t *testing.T) {
	tests := []struct {
		name          string
		setup         func(*testing.T)
		cleanup       func(*testing.T)
		expectError   bool
		validateValue func(*testing.T, *testDurationStruct)
	}{
		{
			name: "duration with default",
			setup: func(t *testing.T) {
				// Don't set TEST_DURATION_DEFAULT
			},
			cleanup:     func(t *testing.T) {},
			expectError: false,
			validateValue: func(t *testing.T, s *testDurationStruct) {
				expected := 10 * time.Minute
				if s.DefaultValue != expected {
					t.Errorf("DefaultValue = %v, want %v", s.DefaultValue, expected)
				}
			},
		},
		{
			name: "optional duration not set",
			setup: func(t *testing.T) {
				// Don't set TEST_DURATION_OPT
			},
			cleanup:     func(t *testing.T) {},
			expectError: false,
			validateValue: func(t *testing.T, s *testDurationStruct) {
				if s.OptionalValue != 0 {
					t.Errorf("OptionalValue = %v, want 0", s.OptionalValue)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.setup != nil {
				tt.setup(t)
			}
			if tt.cleanup != nil {
				defer tt.cleanup(t)
			}

			s := &testDurationStruct{}
			err := LoadEnvForStruct(s)

			if tt.expectError {
				if err == nil {
					t.Errorf("LoadEnvForStruct() expected error, got nil")
				}
			} else {
				if err != nil {
					t.Errorf("LoadEnvForStruct() unexpected error = %v", err)
				} else if tt.validateValue != nil {
					tt.validateValue(t, s)
				}
			}
		})
	}
}

func TestLoadEnvForStruct_String(t *testing.T) {
	// Regression test for existing string functionality
	tests := []struct {
		name        string
		envValue    string
		expectError bool
	}{
		{
			name:        "simple string",
			envValue:    "hello world",
			expectError: false,
		},
		{
			name:        "empty string",
			envValue:    "",
			expectError: true, // Required field
		},
		{
			name:        "string with special characters",
			envValue:    "hello@#$%^&*()",
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.envValue != "" {
				os.Setenv("TEST_STRING", tt.envValue)
				defer os.Unsetenv("TEST_STRING")
			} else {
				os.Unsetenv("TEST_STRING")
			}

			s := &testStringStruct{}
			err := LoadEnvForStruct(s)

			if tt.expectError {
				if err == nil {
					t.Errorf("LoadEnvForStruct() expected error, got nil")
				}
			} else {
				if err != nil {
					t.Errorf("LoadEnvForStruct() unexpected error = %v", err)
				} else if s.Value != tt.envValue {
					t.Errorf("Value = %q, want %q", s.Value, tt.envValue)
				}
			}
		})
	}
}

func TestLoadEnvForStruct_StringSlice(t *testing.T) {
	// Regression test for existing []string functionality
	tests := []struct {
		name          string
		envValue      string
		expectedSlice []string
	}{
		{
			name:          "single value",
			envValue:      "value1",
			expectedSlice: []string{"value1"},
		},
		{
			name:          "multiple values",
			envValue:      "value1,value2,value3",
			expectedSlice: []string{"value1", "value2", "value3"},
		},
		{
			name:          "values with spaces",
			envValue:      "value 1,value 2",
			expectedSlice: []string{"value 1", "value 2"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			os.Setenv("TEST_SLICE", tt.envValue)
			defer os.Unsetenv("TEST_SLICE")

			s := &testSliceStruct{}
			err := LoadEnvForStruct(s)

			if err != nil {
				t.Errorf("LoadEnvForStruct() unexpected error = %v", err)
			} else {
				if len(s.Values) != len(tt.expectedSlice) {
					t.Errorf("Values length = %d, want %d", len(s.Values), len(tt.expectedSlice))
				} else {
					for i := range s.Values {
						if s.Values[i] != tt.expectedSlice[i] {
							t.Errorf("Values[%d] = %q, want %q", i, s.Values[i], tt.expectedSlice[i])
						}
					}
				}
			}
		})
	}
}

func TestLoadEnvForStruct_StringPointer(t *testing.T) {
	// Regression test for existing *string functionality
	t.Run("string pointer set", func(t *testing.T) {
		os.Setenv("TEST_PTR", "pointer value")
		defer os.Unsetenv("TEST_PTR")

		s := &testPointerStruct{}
		err := LoadEnvForStruct(s)

		if err != nil {
			t.Errorf("LoadEnvForStruct() unexpected error = %v", err)
		} else if s.Value == nil {
			t.Errorf("Value is nil, want non-nil")
		} else if *s.Value != "pointer value" {
			t.Errorf("*Value = %q, want %q", *s.Value, "pointer value")
		}
	})
}

func TestLoadEnvForStruct_Mixed(t *testing.T) {
	// Test struct with mixed types
	t.Run("all types together", func(t *testing.T) {
		os.Setenv("MIXED_STRING", "test")
		os.Setenv("MIXED_INT", "42")
		os.Setenv("MIXED_BOOL", "true")
		os.Setenv("MIXED_DURATION", "5m")
		defer func() {
			os.Unsetenv("MIXED_STRING")
			os.Unsetenv("MIXED_INT")
			os.Unsetenv("MIXED_BOOL")
			os.Unsetenv("MIXED_DURATION")
		}()

		s := &testMixedStruct{}
		err := LoadEnvForStruct(s)

		if err != nil {
			t.Errorf("LoadEnvForStruct() unexpected error = %v", err)
		} else {
			if s.StringValue != "test" {
				t.Errorf("StringValue = %q, want %q", s.StringValue, "test")
			}
			if s.IntValue != 42 {
				t.Errorf("IntValue = %d, want 42", s.IntValue)
			}
			if !s.BoolValue {
				t.Errorf("BoolValue = %v, want true", s.BoolValue)
			}
			if s.DurationValue != 5*time.Minute {
				t.Errorf("DurationValue = %v, want %v", s.DurationValue, 5*time.Minute)
			}
		}
	})
}

func TestLoadEnvForStruct_InvalidInput(t *testing.T) {
	t.Run("non-struct", func(t *testing.T) {
		var i int
		err := LoadEnvForStruct(&i)

		if err == nil {
			t.Errorf("LoadEnvForStruct() expected error, got nil")
		} else if !strings.Contains(err.Error(), "expected a struct") {
			t.Errorf("LoadEnvForStruct() error = %v, want error containing %q", err, "expected a struct")
		}
	})
}

// Test structs
type testIntStruct struct {
	Value         int `env:"TEST_INT" optional:"true"`
	DefaultValue  int `env:"TEST_INT_DEFAULT" default:"99"`
	OptionalValue int `env:"TEST_INT_OPT" optional:"true"`
	MultiVar      int `env:"TEST_INT_PRIMARY,TEST_INT_ALT" optional:"true"`
}

type testIntTypesStruct struct {
	Int8Value  int8  `env:"TEST_INT8"`
	Int16Value int16 `env:"TEST_INT16"`
	Int32Value int32 `env:"TEST_INT32"`
	Int64Value int64 `env:"TEST_INT64"`
}

type testBoolStruct struct {
	Value         bool `env:"TEST_BOOL" optional:"true"`
	DefaultValue  bool `env:"TEST_BOOL_DEFAULT" default:"true"`
	DefaultFalse  bool `env:"TEST_BOOL_DEFAULT_FALSE" default:"false"`
	OptionalValue bool `env:"TEST_BOOL_OPT" optional:"true"`
}

type testDurationStruct struct {
	Value         time.Duration `env:"TEST_DURATION" optional:"true"`
	DefaultValue  time.Duration `env:"TEST_DURATION_DEFAULT" default:"10m"`
	OptionalValue time.Duration `env:"TEST_DURATION_OPT" optional:"true"`
}

type testStringStruct struct {
	Value string `env:"TEST_STRING"`
}

type testSliceStruct struct {
	Values []string `env:"TEST_SLICE"`
}

type testPointerStruct struct {
	Value *string `env:"TEST_PTR"`
}

type testMixedStruct struct {
	StringValue   string        `env:"MIXED_STRING"`
	IntValue      int           `env:"MIXED_INT"`
	BoolValue     bool          `env:"MIXED_BOOL"`
	DurationValue time.Duration `env:"MIXED_DURATION"`
}
