package account

import (
	"testing"
)

func BenchmarkArgon2idHash(b *testing.B) {
	password := []byte("correct horse battery staple")
	salt, err := randomSalt()
	if err != nil {
		b.Fatal(err)
	}

	b.ReportAllocs()
	b.ResetTimer()
	for b.Loop() {
		_ = HashPasswordArgon2id(password, salt)
	}
}

func BenchmarkArgon2idHashParallel(b *testing.B) {
	password := []byte("correct horse battery staple")
	salt, err := randomSalt()
	if err != nil {
		b.Fatal(err)
	}

	b.ReportAllocs()
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_ = HashPasswordArgon2id(password, salt)
		}
	})
}
