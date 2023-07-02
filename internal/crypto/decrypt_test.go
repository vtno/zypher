package crypto_test

import (
	"encoding/base64"
	"errors"
	"io/fs"
	"os"
	"testing"

	"github.com/vtno/zypher/internal/crypto"
	"go.uber.org/mock/gomock"
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
		envs            map[string]string
		initMocks       func() (crypto.CipherFactory, crypto.FileReaderWriter)
	}
	base64Content := base64.StdEncoding.EncodeToString([]byte("encryptedcontent"))
	tests := []test{
		{
			name:            "run successfully with input from args",
			args:            []string{"-k", "key", base64Content},
			expectedErrCode: 0,
			initMocks: func() (crypto.CipherFactory, crypto.FileReaderWriter) {
				mockCipherFactory := crypto.NewMockCipherFactory(ctrl)
				mockCipher := crypto.NewMockCipher(ctrl)
				mockCipher.EXPECT().Decrypt([]byte("encryptedcontent")).Times(1)
				mockCipherFactory.EXPECT().NewCipher(gomock.Any()).Return(mockCipher).Times(1)

				return mockCipherFactory, nil
			},
		},
		{
			name:            "run successfully with input from file and output into a file",
			args:            []string{"-k", "key", "-f", "input.enc", "-o", "input.txt"},
			expectedErrCode: 0,
			initMocks: func() (crypto.CipherFactory, crypto.FileReaderWriter) {
				mockCipherFactory := crypto.NewMockCipherFactory(ctrl)
				mockCipher := crypto.NewMockCipher(ctrl)
				mockCipher.EXPECT().Decrypt([]byte("encryptedcontent")).Return([]byte("encryptedcontent"), nil).Times(1)
				mockCipherFactory.EXPECT().NewCipher(gomock.Any()).Return(mockCipher).Times(1)
				mockFileReaderWriter := crypto.NewMockFileReaderWriter(ctrl)
				mockFileReaderWriter.EXPECT().ReadFile("input.enc").Return([]byte(base64Content), nil).Times(1)
				mockFileReaderWriter.EXPECT().WriteFile("input.txt", []byte("encryptedcontent"), fs.FileMode(0600)).Times(1)
				return mockCipherFactory, mockFileReaderWriter
			},
		},
		{
			name:            "run successfully with key from ZYPHER_KEY env",
			args:            []string{"-f", "input.enc"},
			expectedErrCode: 0,
			envs:            map[string]string{"ZYPHER_KEY": "key"},
			initMocks: func() (crypto.CipherFactory, crypto.FileReaderWriter) {
				mockCipherFactory := crypto.NewMockCipherFactory(ctrl)
				mockCipher := crypto.NewMockCipher(ctrl)
				mockCipher.EXPECT().Decrypt([]byte("encryptedcontent")).Times(1)
				mockCipherFactory.EXPECT().NewCipher(gomock.Any()).Return(mockCipher).Times(1)
				mockFileReaderWriter := crypto.NewMockFileReaderWriter(ctrl)
				mockFileReaderWriter.EXPECT().ReadFile("zypher.key").Return(nil, errors.New("file not exist")).Times(1)
				mockFileReaderWriter.EXPECT().ReadFile("input.enc").Return([]byte(base64Content), nil).Times(1)
				return mockCipherFactory, mockFileReaderWriter
			},
		},
		{
			name:            "run successfully with key from zypher.key file",
			args:            []string{"-f", "input.enc"},
			expectedErrCode: 0,
			initMocks: func() (crypto.CipherFactory, crypto.FileReaderWriter) {
				mockCipherFactory := crypto.NewMockCipherFactory(ctrl)
				mockCipher := crypto.NewMockCipher(ctrl)
				mockCipher.EXPECT().Decrypt([]byte("encryptedcontent")).Times(1)
				mockCipherFactory.EXPECT().NewCipher(gomock.Any()).Return(mockCipher).Times(1)
				mockFileReaderWriter := crypto.NewMockFileReaderWriter(ctrl)
				mockFileReaderWriter.EXPECT().ReadFile("zypher.key").Return([]byte("key"), nil).Times(1)
				mockFileReaderWriter.EXPECT().ReadFile("input.enc").Return([]byte(base64Content), nil).Times(1)
				return mockCipherFactory, mockFileReaderWriter
			},
		},
		{
			name:            "run successfully with key from overridden another.key file",
			args:            []string{"-f", "input.enc", "-kf", "another.key"},
			expectedErrCode: 0,
			initMocks: func() (crypto.CipherFactory, crypto.FileReaderWriter) {
				mockCipherFactory := crypto.NewMockCipherFactory(ctrl)
				mockCipher := crypto.NewMockCipher(ctrl)
				mockCipher.EXPECT().Decrypt([]byte("encryptedcontent")).Times(1)
				mockCipherFactory.EXPECT().NewCipher(gomock.Any()).Return(mockCipher).Times(1)
				mockFileReaderWriter := crypto.NewMockFileReaderWriter(ctrl)
				mockFileReaderWriter.EXPECT().ReadFile("another.key").Return([]byte("anotherkey"), nil).Times(1)
				mockFileReaderWriter.EXPECT().ReadFile("input.enc").Return([]byte(base64Content), nil).Times(1)
				return mockCipherFactory, mockFileReaderWriter
			},
		},
		{
			name:            "fails when key not provided",
			args:            []string{"encryptedtext"},
			expectedErrCode: 1,
			initMocks: func() (crypto.CipherFactory, crypto.FileReaderWriter) {
				mockCipherFactory := crypto.NewMockCipherFactory(ctrl)
				mockCipher := crypto.NewMockCipher(ctrl)
				mockCipher.EXPECT().Decrypt(gomock.Any()).Times(0)
				mockFileReaderWriter := crypto.NewMockFileReaderWriter(ctrl)
				mockFileReaderWriter.EXPECT().ReadFile("zypher.key").Return(nil, errors.New("file not exist")).Times(1)
				mockCipherFactory.EXPECT().NewCipher(gomock.Any()).Times(0)
				return mockCipherFactory, mockFileReaderWriter
			},
		},
		{
			name:            "fails when file not exist",
			args:            []string{"-k", "key", "-f", "not-exist.txt"},
			expectedErrCode: 1,
			initMocks: func() (crypto.CipherFactory, crypto.FileReaderWriter) {
				mockCipherFactory := crypto.NewMockCipherFactory(ctrl)
				mockCipher := crypto.NewMockCipher(ctrl)
				mockCipher.EXPECT().Decrypt(gomock.Any()).Times(0)
				mockCipherFactory.EXPECT().NewCipher(gomock.Any()).Return(mockCipher).Times(1)
				mockFileReaderWriter := crypto.NewMockFileReaderWriter(ctrl)
				mockFileReaderWriter.EXPECT().ReadFile("not-exist.txt").Return(nil, errors.New("file not exist")).Times(1)
				return mockCipherFactory, mockFileReaderWriter
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			os.Clearenv()
			for k, v := range tt.envs {
				os.Setenv(k, v)
			}
			mockCipherFactory, mockFileReaderWriter := tt.initMocks()
			decryptCmd := crypto.NewDecryptCmd(
				mockCipherFactory,
				crypto.WithFileReaderWriter(mockFileReaderWriter),
			)
			errCode := decryptCmd.Run(tt.args)
			if errCode != tt.expectedErrCode {
				t.Errorf("Expected code %d, got %d", tt.expectedErrCode, errCode)
			}
		})
	}
}
