package crypto_test

import (
	"os"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/vtno/zypher/internal/crypto"
)

func TestEncrypt_Help(t *testing.T) {
	ctrl := gomock.NewController(t)
	encryptCmd := crypto.NewEncryptCmd(crypto.NewMockCipherFactory(ctrl))
	msg := encryptCmd.Help()
	if msg != crypto.HelpMsg {
		t.Errorf("Expected correct help message, got %s", msg)
	}
}

func TestEncrypt_Synopsis(t *testing.T) {
	ctrl := gomock.NewController(t)
	encryptCmd := crypto.NewEncryptCmd(crypto.NewMockCipherFactory(ctrl))
	msg := encryptCmd.Synopsis()
	if msg != crypto.SynopsisMsg {
		t.Errorf("Expected correct synopsis message, got %s", msg)
	}
}

func TestEncrypt_Run(t *testing.T) {
	ctrl := gomock.NewController(t)

	type test struct {
		name            string
		args            []string
		expectedErrCode int
		envs            map[string]string
		initMocks       func() crypto.CipherFactory
	}

	tests := []test{
		{
			name:            "run successfully with input from args",
			args:            []string{"-k", "key", "sometext"},
			expectedErrCode: 0,
			initMocks: func() crypto.CipherFactory {
				mockCipherFactory := crypto.NewMockCipherFactory(ctrl)
				mockCipher := crypto.NewMockCipher(ctrl)
				mockCipher.EXPECT().Encrypt([]byte("sometext")).Times(1)
				mockCipherFactory.EXPECT().NewCipher(gomock.Any()).Return(mockCipher).Times(1)
				return mockCipherFactory
			},
		},
		{
			name:            "run successfully with input from file",
			args:            []string{"-k", "key", "-f", "input.txt", "-o", "input.enc"},
			expectedErrCode: 0,
			initMocks: func() crypto.CipherFactory {
				mockCipherFactory := crypto.NewMockCipherFactory(ctrl)
				mockCipher := crypto.NewMockCipher(ctrl)
				mockCipher.EXPECT().Encrypt([]byte("content")).Times(1)
				mockCipherFactory.EXPECT().NewCipher(gomock.Any()).Return(mockCipher).Times(1)
				return mockCipherFactory
			},
		},
		{
			name:            "run successfully with key from ZYPHER_KEY env",
			args:            []string{"sometext"},
			envs:            map[string]string{"ZYPHER_KEY": "key"},
			expectedErrCode: 0,
			initMocks: func() crypto.CipherFactory {
				mockCipherFactory := crypto.NewMockCipherFactory(ctrl)
				mockCipher := crypto.NewMockCipher(ctrl)
				mockCipher.EXPECT().Encrypt([]byte("sometext")).Times(1)
				mockCipherFactory.EXPECT().NewCipher(gomock.Any()).Return(mockCipher).Times(1)
				return mockCipherFactory
			},
		},
		{
			name:            "fails when no key provided",
			args:            []string{"sometext"},
			expectedErrCode: 1,
			initMocks: func() crypto.CipherFactory {
				mockCipherFactory := crypto.NewMockCipherFactory(ctrl)
				mockCipher := crypto.NewMockCipher(ctrl)
				mockCipher.EXPECT().Encrypt(gomock.Any()).Times(0)
				mockCipherFactory.EXPECT().NewCipher(gomock.Any()).Return(mockCipher).Times(0)
				return mockCipherFactory
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			os.Clearenv()
			for k, v := range tt.envs {
				os.Setenv(k, v)
			}
			mockCipherFactory := tt.initMocks()
			encryptCmd := crypto.NewEncryptCmd(
				mockCipherFactory,
				crypto.WithFileReader(func(s string) ([]byte, error) {
					return []byte("content"), nil
				}),
			)
			errCode := encryptCmd.Run(tt.args)
			if errCode != tt.expectedErrCode {
				t.Errorf("Expected code %d, got %d", tt.expectedErrCode, errCode)
			}
		})
	}
}