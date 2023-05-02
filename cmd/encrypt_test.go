package cmd_test

import (
	"testing"

	gomock "github.com/golang/mock/gomock"
	"github.com/vtno/zypher/cmd"
)

func TestEncrypt_Help(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	m := cmd.NewMockCipher(ctrl)
	encrypt := cmd.NewEncryptCmd(m)
	msg := encrypt.Help()
	if msg != cmd.HelpMsg {
		t.Errorf("Expected correct help message, got %s", msg)
	}
}

func TestEncrypt_Synopsis(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	m := cmd.NewMockCipher(ctrl)
	encrypt := cmd.NewEncryptCmd(m)
	msg := encrypt.Synopsis()
	if msg != "encrypts input value or file with the provided key and prints the encrypted value to stdout or create a file" {
		t.Errorf("Expected correct synopsis message, got %s", msg)
	}
}