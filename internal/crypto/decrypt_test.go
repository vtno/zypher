package crypto_test

import (
	"testing"

	gomock "github.com/golang/mock/gomock"
	"github.com/vtno/zypher/internal/crypto"
)

func TestDecrypt_Help(t *testing.T) {
	ctrl := gomock.NewController(t)
	decryptCmd := crypto.NewDecryptCmd(crypto.NewMockCipherFactory(ctrl))
	msg := decryptCmd.Help()
	if msg != crypto.DecryptHelpMsg {
		t.Errorf("Expected correct help message, got %s", msg)
	}
}

func TestDecrypt_Synopsis(t *testing.T) {
	ctrl := gomock.NewController(t)
	decryptCmd := crypto.NewDecryptCmd(crypto.NewMockCipherFactory(ctrl))
	msg := decryptCmd.Synopsis()
	if msg != crypto.DecryptSynopsis {
		t.Errorf("Expected correct synopsis message, got %s", msg)
	}
}

func TestDecrypt_Run(t *testing.T) {
	ctrl := gomock.NewController(t)

	type test struct {
		name            string
		args            []string
		expectedErrCode int
		// envs            []map[string]string
		initMocks       func() crypto.CipherFactory
		useOSFileReader bool
	}

	tests := []test{
		{
			name:            "run successfully with input from args",
			args:            []string{"-k", "key", "encryptedcontent"},
			expectedErrCode: 0,
			initMocks: func() crypto.CipherFactory {
				mockCipherFactory := crypto.NewMockCipherFactory(ctrl)
				mockCipher := crypto.NewMockCipher(ctrl)
				mockCipher.EXPECT().Decrypt([]byte("encryptedcontent")).Times(1)
				mockCipherFactory.EXPECT().NewCipher(gomock.Any()).Return(mockCipher).Times(1)
				return mockCipherFactory
			},
		},
		{
			name:            "run successfully with input from file",
			args:            []string{"-k", "key", "-f", "input.enc"},
			expectedErrCode: 0,
			initMocks: func() crypto.CipherFactory {
				mockCipherFactory := crypto.NewMockCipherFactory(ctrl)
				mockCipher := crypto.NewMockCipher(ctrl)
				mockCipher.EXPECT().Decrypt([]byte("content")).Times(1)
				mockCipherFactory.EXPECT().NewCipher(gomock.Any()).Return(mockCipher).Times(1)
				return mockCipherFactory
			},
		},
		{
			name:            "fails when key not provided",
			args:            []string{"encryptedtext"},
			expectedErrCode: 1,
			initMocks: func() crypto.CipherFactory {
				mockCipherFactory := crypto.NewMockCipherFactory(ctrl)
				mockCipher := crypto.NewMockCipher(ctrl)
				mockCipher.EXPECT().Decrypt(gomock.Any()).Times(0)
				mockCipherFactory.EXPECT().NewCipher(gomock.Any()).Times(0)
				return mockCipherFactory
			},
		},
		{
			name:            "fails when file not exist",
			args:            []string{"-k", "key", "-f", "not-exist.txt"},
			expectedErrCode: 1,
			initMocks: func() crypto.CipherFactory {
				mockCipherFactory := crypto.NewMockCipherFactory(ctrl)
				mockCipher := crypto.NewMockCipher(ctrl)
				mockCipher.EXPECT().Decrypt(gomock.Any()).Times(0)
				mockCipherFactory.EXPECT().NewCipher(gomock.Any()).Return(mockCipher).Times(1)
				return mockCipherFactory
			},
			useOSFileReader: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockCipherFactory := tt.initMocks()
			var decryptCmd *crypto.DecryptCmd
			if tt.useOSFileReader {
				decryptCmd = crypto.NewDecryptCmd(mockCipherFactory)
			} else {
				decryptCmd = crypto.NewDecryptCmd(
					mockCipherFactory,
					crypto.WithFileReader(func(s string) ([]byte, error) {
						return []byte("content"), nil
					}),
				)
			}
			errCode := decryptCmd.Run(tt.args)
			if errCode != tt.expectedErrCode {
				t.Errorf("Expected code %d, got %d", tt.expectedErrCode, errCode)
			}
		})
	}
}
