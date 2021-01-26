package api

import "testing"

func Test_validateIPv4(t *testing.T) {
	tests := []struct {
		name string
		v    IPv4
		want bool
	}{
		{"test_0", "192.168.1.7", true},
		{"test_1", "192.168.1.7.i", false},
		{"test_2", "192.AS.1.7.i", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := validateIPv4(tt.v); got != tt.want {
				t.Errorf("validateIPv4() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_validateIPv6(t *testing.T) {
	tests := []struct {
		name string
		v    IPv6
		want bool
	}{
		{"test_0", "fe80::4ca0:8fff:fe70:52d", true},
		{"test_1", "fe80::4ca0:8fff:fe70:52d29", false},
		{"test_2", "fe80::4ca0:8fff:fe70:52z", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := validateIPv6(tt.v); got != tt.want {
				t.Errorf("validateIPv4() = %v, want %v", got, tt.want)
			}
		})
	}
}
