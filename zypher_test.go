package zypher_test

import (
	"testing"

	"github.com/vtno/zypher"
)

func TestCipher_Encrypt(t *testing.T) {
	type test struct {
		ci          *zypher.Cipher
		name        string
		input       []byte
		expectError bool
	}

	tests := []test{
		{
			ci:          zypher.NewCipher("1234"),
			name:        "errors when key size is not 16, 24, 32",
			input:       []byte("somelongkey"),
			expectError: true,
		},
		{
			ci:          zypher.NewCipher("1234567890123456"),
			name:        "encrypt correctly when key size is 16",
			input:       []byte("somelongkey"),
			expectError: false,
		},
		{
			ci:          zypher.NewCipher("123456789012345678901234"),
			name:        "encrypt correctly when key size is 24",
			input:       []byte("somelongkey"),
			expectError: false,
		},
		{
			ci:          zypher.NewCipher("12345678901234567890123456789012"),
			name:        "encrypt correctly when key size is 32",
			input:       []byte("somelongkey"),
			expectError: false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, err := test.ci.Encrypt(test.input)
			if test.expectError && err == nil {
				t.Errorf("expected error, got nil")
			}
			if !test.expectError && err != nil {
				t.Errorf("expected no error, got %v", err)
			}
			if string(test.input) == string(got) {
				t.Errorf("expected encrypted value to be different than the input: %s, %s", test.input, got)
			}
		})
	}
}

func TestCipher_Decrypt(t *testing.T) {
	encryptor := zypher.NewCipher("1234567890123456")
	decryptor := zypher.NewCipher("1134567890123456")

	type test struct {
		encryptor   *zypher.Cipher
		decryptor   *zypher.Cipher
		name        string
		input       []byte
		expectError bool
	}

	tests := []test{
		{
			encryptor:   encryptor,
			decryptor:   decryptor,
			name:        "fail to decrypt when key is invalid",
			input:       []byte("somelongkey"),
			expectError: true,
		},
		{
			encryptor:   decryptor,
			decryptor:   decryptor,
			name:        "decrypt correctly when key is valid",
			input:       []byte("somelongkey"),
			expectError: false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			encryptedText, err := test.encryptor.Encrypt(test.input)
			if err != nil {
				t.Errorf("unexpected error on Encrypt: %v", err)
			}
			got, err := test.decryptor.Decrypt(encryptedText)
			if test.expectError && err == nil {
				t.Errorf("expected error, got nil")
			}
			if !test.expectError {
				if err != nil {
					t.Errorf("expected no error, got %v", err)
				}
				if string(test.input) != string(got) {
					t.Errorf("expected decrypted value to be equal to the input: %s, %s", test.input, got)
				}
			}
		})
	}
}
