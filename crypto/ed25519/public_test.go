package ed25519

import (
	"reflect"
	"testing"

	"github.com/mailchain/mailchain/crypto"
	"github.com/stretchr/testify/assert"
)

func TestPublicKey_Bytes(t *testing.T) {
	tests := []struct {
		name string
		pk   PublicKey
		want []byte
	}{
		{
			"sofia",
			sofiaPublicKey,
			sofiaPublicKeyBytes,
		},
		{
			"charlotte",
			charlottePublicKey,
			charlottePublicKeyBytes,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.pk.Bytes(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("PublicKey.Bytes() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPublicKeyFromBytes(t *testing.T) {
	type args struct {
		keyBytes []byte
	}
	tests := []struct {
		name    string
		args    args
		want    crypto.PublicKey
		wantErr bool
	}{
		{
			"sofia",
			args{
				sofiaPublicKeyBytes,
			},
			&sofiaPublicKey,
			false,
		},
		{
			"err-too-short",
			args{
				[]byte{0x72, 0x3c, 0xaa, 0x23},
			},
			nil,
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := PublicKeyFromBytes(tt.args.keyBytes)
			if (err != nil) != tt.wantErr {
				t.Errorf("PublicKeyFromBytes() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !assert.Equal(t, tt.want, got) {
				t.Errorf("PublicKeyFromBytes() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPublicKey_Kind(t *testing.T) {
	tests := []struct {
		name string
		pk   PublicKey
		want string
	}{
		{
			"charlotte",
			charlottePublicKey,
			crypto.KindED25519,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.pk.Kind(); !assert.Equal(t, tt.want, got) {
				t.Errorf("PublicKey.Kind() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPublicKey_Verify(t *testing.T) {
	tests := []struct {
		name    string
		pk      PublicKey
		message []byte
		sig     []byte
		want    bool
	}{
		{
			"success-charlotte",
			charlottePublicKey,
			[]byte("message"),
			[]byte{0x7d, 0x51, 0xea, 0xfa, 0x52, 0x78, 0x31, 0x69, 0xd0, 0xa9, 0x4a, 0xc, 0x9f, 0x2b, 0xca, 0xd5, 0xe0, 0x3d, 0x29, 0x17, 0x33, 0x0, 0x93, 0xf, 0xf3, 0xc7, 0xd6, 0x3b, 0xfd, 0x64, 0x17, 0xae, 0x1b, 0xc8, 0x1f, 0xef, 0x51, 0xba, 0x14, 0x9a, 0xe8, 0xa1, 0xe1, 0xda, 0xe0, 0x5f, 0xdc, 0xa5, 0x7, 0x8b, 0x14, 0xba, 0xc4, 0xcf, 0x26, 0xcc, 0xc6, 0x1, 0x1e, 0x5e, 0xab, 0x77, 0x3, 0xc},
			true,
		},
		{
			"success-sofia",
			sofiaPublicKey,
			[]byte("egassem"),
			[]byte{0xde, 0x6c, 0x88, 0xe6, 0x9c, 0x9f, 0x93, 0xb, 0x59, 0xdd, 0xf4, 0x80, 0xc2, 0x9a, 0x55, 0x79, 0xec, 0x89, 0x5c, 0xa9, 0x7a, 0x36, 0xf6, 0x69, 0x74, 0xc1, 0xf0, 0x15, 0x5c, 0xc0, 0x66, 0x75, 0x2e, 0xcd, 0x9a, 0x9b, 0x41, 0x35, 0xd2, 0x72, 0x32, 0xe0, 0x54, 0x80, 0xbc, 0x98, 0x58, 0x1, 0xa9, 0xfd, 0xe4, 0x27, 0xc7, 0xef, 0xa5, 0x42, 0x5f, 0xf, 0x46, 0x49, 0xb8, 0xad, 0xbd, 0x5},
			true,
		},
		{
			"err-invalid-signature-charlotte",
			charlottePublicKey,
			[]byte("message"),
			[]byte{0x2e, 0x11, 0x79, 0x88, 0xb7, 0xc, 0x44, 0xac, 0xdb, 0xe0, 0x27, 0x2a, 0x30, 0x7b, 0x42, 0xf, 0x21, 0xe, 0x5a, 0x79, 0x37, 0x2b, 0x5a, 0xf4, 0x3d, 0x6a, 0x5, 0xf9, 0xab, 0xa, 0x83, 0x4, 0x99, 0x95, 0x5c, 0xc8, 0x98, 0x4, 0xeb, 0x21, 0x4, 0x14, 0x95, 0x1b, 0x79, 0xbc, 0x67, 0xa6, 0x4, 0x66, 0xc9, 0xa4, 0xec, 0xc0, 0xc2, 0x42, 0x51, 0x38, 0xee, 0x29, 0xe9, 0x54, 0x2b, 0x3},
			false,
		},
		{
			"err-invalid-signature-sofia",
			sofiaPublicKey,
			[]byte("egassem"),
			[]byte{0x3, 0x5e, 0x29, 0xa9, 0xd8, 0xb3, 0x6c, 0xb8, 0x92, 0x8, 0xc9, 0x45, 0x15, 0xf9, 0x7c, 0x4d, 0x56, 0x7c, 0x36, 0x84, 0x34, 0x2, 0xcc, 0x51, 0xa2, 0x45, 0x5a, 0x39, 0xce, 0x11, 0xf2, 0x7, 0x33, 0xa7, 0xe2, 0x37, 0x9c, 0x8, 0xa4, 0x60, 0x1a, 0x23, 0x7e, 0x24, 0xb5, 0x49, 0x67, 0xf0, 0x8b, 0x78, 0xea, 0x28, 0x97, 0xe, 0x2, 0xb1, 0xaa, 0x76, 0x8b, 0x6b, 0xfc, 0xa2, 0x60, 0xc},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.pk.Verify(tt.message, tt.sig)
			if !assert.Equal(t, tt.want, got) {
				t.Errorf("PublicKey.Verify() = %v, want %v", got, tt.want)
			}
		})
	}
}
