package encrypt_test

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/vtno/zypher/internal/encrypt"
)

func TestEncrypt_Help(t *testing.T) {
	ctrl := gomock.NewController(t)
	encryptCmd := encrypt.NewEncryptCmd(encrypt.NewMockCipherFactory(ctrl))
	msg := encryptCmd.Help()
	if msg != encrypt.HelpMsg {
		t.Errorf("Expected correct help message, got %s", msg)
	}
}

func TestEncrypt_Synopsis(t *testing.T) {
	ctrl := gomock.NewController(t)
	encryptCmd := encrypt.NewEncryptCmd(encrypt.NewMockCipherFactory(ctrl))
	msg := encryptCmd.Synopsis()
	if msg != "encrypts input value or file with the provided key and prints the encrypted value to stdout or create a file" {
		t.Errorf("Expected correct synopsis message, got %s", msg)
	}
}

func TestEncrypt_Run(t *testing.T) {
	ctrl := gomock.NewController(t)
	encryptCmd := encrypt.NewEncryptCmd(encrypt.NewMockCipherFactory(ctrl))
	errCode := encryptCmd.Run([]string{"-k", "key", "-f", "input.txt", "-o", "input.enc"})
	if errCode != 0 {
		t.Errorf("Expected code 0, got %d", errCode)
	}
}
