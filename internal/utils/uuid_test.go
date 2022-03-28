package utils

import (
	"testing"
)

func Test_GetUUID1(t *testing.T) {
	uuid, err := GetUUID1()
	if err != nil {
		t.Error(err)
	}
	t.Log(uuid)
}

func Test_GetUUID4(t *testing.T) {
	uuid, err := GetUUID4()
	if err != nil {
		t.Error(err)
	}
	t.Log(uuid)
}

func Test_GenerateRequestID(t *testing.T) {
	uuid := GenerateRequestID()
	t.Log(uuid)
}

func Benchmark_GenerateRequestIDBatch(b *testing.B) {
	b.StopTimer()
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		uuid := GenerateRequestID()
		b.Log(uuid)
	}
}
