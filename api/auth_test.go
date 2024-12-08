package api

import (
	"testing"
)

func TestPassValidation(t *testing.T) {
	type wanted struct {
		valid   bool
		message string
	}
	tests := []struct {
		name  string
		input string
		want  wanted
	}{
		{"Pass shorter than 8 letters should be invalid", "pass", wanted{false, "Password is too short"}},
		{"Pass without any numbers should be invalid", "Pa._ssWord", wanted{false, "Password has to contain at least 1 number"}},
		{"Pass without uppercase letter should be invalid", "2a._ss4ord", wanted{false, "Password has to contain at least 1 uppercase letter"}},
		{"Pass without lowercase letter should be invalid", "PASS_W0RD", wanted{false, "Password has to contain at least 1 lowercase letter"}},
		{"Pass without symbol should be invalid", "PassW0rd", wanted{false, "Password has to contain at least 1 symbol"}},
		{"Pass should be valid", "Pass_W0rd", wanted{true, ""}},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			validity, message := IsValidPassword(&test.input)
			if validity != test.want.valid || message != test.want.message {
				t.Errorf("Got %+v, want %+v", wanted{validity, message}, test.want)
			}
		})
	}
}
