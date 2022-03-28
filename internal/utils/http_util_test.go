package utils

import (
	"testing"
)

func Test_LocalIP(t *testing.T) {
	ip := LocalIP()
	t.Logf("local ip %v", ip)
}

func Test_LocalIPByTCP(t *testing.T) {
	ip := LocalIPByTCP()
	t.Logf("local ip %v", ip)
}

func Test_LocalIPWithVal(t *testing.T) {
	ip := LocalIPWithVal("localhost")
	t.Logf("local ip %v", ip)
}

func Test_PublicIP(t *testing.T) {
	ip, err := PublicIP()
	if err != nil {
		t.Error(err)
	}
	t.Logf("public ip %v", ip)
}

func Test_FreePort(t *testing.T) {
	port, err := FreePort("tcp")
	if err != nil {
		t.Error(err)
	}
	t.Logf("free port %v", port)
}
