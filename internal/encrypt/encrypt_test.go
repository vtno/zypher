package encrypt_test

import (
	// "os"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/vtno/zypher/internal/encrypt"
	// "golang.org/x/exp/maps"
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
	if msg != encrypt.SynopsisMsg {
		t.Errorf("Expected correct synopsis message, got %s", msg)
	}
}

func TestEncrypt_Run(t *testing.T) {
	ctrl := gomock.NewController(t)

	type test struct {
		name            string
		args            []string
		expectedErrCode int
		// envs            []map[string]string
		initMocks func() encrypt.CipherFactory
	}

	tests := []test{
		{
			name:            "run successfully with input from args",
			args:            []string{"-k", "key", "sometext"},
			expectedErrCode: 0,
			initMocks: func() encrypt.CipherFactory {
				mockCipherFactory := encrypt.NewMockCipherFactory(ctrl)
				mockCipher := encrypt.NewMockCipher(ctrl)
				mockCipher.EXPECT().Encrypt([]byte("sometext")).Times(1)
				mockCipherFactory.EXPECT().NewCipher(gomock.Any()).Return(mockCipher).Times(1)
				return mockCipherFactory
			},
		},
		{
			name:            "run successfully with input from file",
			args:            []string{"-k", "key", "-f", "input.txt", "-o", "input.enc"},
			expectedErrCode: 0,
			initMocks: func() encrypt.CipherFactory {
				mockCipherFactory := encrypt.NewMockCipherFactory(ctrl)
				mockCipher := encrypt.NewMockCipher(ctrl)
				mockCipher.EXPECT().Encrypt([]byte("content")).Times(1)
				mockCipherFactory.EXPECT().NewCipher(gomock.Any()).Return(mockCipher).Times(1)
				return mockCipherFactory
			},
		},
		// {
		// 	name:            "run successfully with key from ZYPHER_KEY env",
		// 	args:            []string{"sometext"},
		// 	envs:            []map[string]string{{"ZYPHER_KEY": "key"}},
		// 	expectedErrCode: 0,
		// },
		{
			name:            "fails when no key provided",
			args:            []string{"sometext"},
			expectedErrCode: 1,
			initMocks: func() encrypt.CipherFactory {
				mockCipherFactory := encrypt.NewMockCipherFactory(ctrl)
				mockCipher := encrypt.NewMockCipher(ctrl)
				mockCipher.EXPECT().Encrypt(gomock.Any()).Times(0)
				mockCipherFactory.EXPECT().NewCipher(gomock.Any()).Return(mockCipher).Times(0)
				return mockCipherFactory
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// os.Clearenv()
			// for _, env := range maps.Keys(tt.envs) {
			// 	for _, k := range maps.Keys(env) {
			// 		os.Setenv(k, env[k])
			// 	}
			// }
			mockCipherFactory := tt.initMocks()
			encryptCmd := encrypt.NewEncryptCmd(
				mockCipherFactory,
				encrypt.WithFileReader(func(s string) ([]byte, error) {
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
