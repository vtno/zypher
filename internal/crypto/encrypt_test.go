package crypto_test

import (
	"errors"
	"io/fs"
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
		initMocks       func() (crypto.CipherFactory, crypto.FileReaderWriter)
	}

	tests := []test{
		{
			name:            "run successfully with input from args",
			args:            []string{"-k", "key", "sometext"},
			expectedErrCode: 0,
			initMocks: func() (crypto.CipherFactory, crypto.FileReaderWriter) {
				mockCipherFactory := crypto.NewMockCipherFactory(ctrl)
				mockCipher := crypto.NewMockCipher(ctrl)
				mockCipher.EXPECT().Encrypt([]byte("sometext")).Times(1)
				mockCipherFactory.EXPECT().NewCipher(gomock.Any()).Return(mockCipher).Times(1)
				return mockCipherFactory, nil
			},
		},
		{
			name:            "run successfully with input from file and output to a file",
			args:            []string{"-k", "key", "-f", "input.txt", "-o", "input.enc"},
			expectedErrCode: 0,
			initMocks: func() (crypto.CipherFactory, crypto.FileReaderWriter) {
				mockCipherFactory := crypto.NewMockCipherFactory(ctrl)
				mockCipher := crypto.NewMockCipher(ctrl)
				mockCipher.EXPECT().Encrypt([]byte("content")).Times(1)
				mockCipherFactory.EXPECT().NewCipher(gomock.Any()).Return(mockCipher).Times(1)
				mockFileReaderWriter := crypto.NewMockFileReaderWriter(ctrl)
				mockFileReaderWriter.EXPECT().ReadFile("input.txt").Return([]byte("content"), nil).Times(1)
				mockFileReaderWriter.EXPECT().WriteFile("input.enc", gomock.Any(), fs.FileMode(0600)).Times(1)
				return mockCipherFactory, mockFileReaderWriter
			},
		},
		{
			name:            "run successfully with key from zypher.key file",
			args:            []string{"-f", "input.txt", "-o", "input.enc"},
			expectedErrCode: 0,
			initMocks: func() (crypto.CipherFactory, crypto.FileReaderWriter) {
				mockCipherFactory := crypto.NewMockCipherFactory(ctrl)
				mockCipher := crypto.NewMockCipher(ctrl)
				mockCipher.EXPECT().Encrypt([]byte("content")).Times(1)
				mockCipherFactory.EXPECT().NewCipher(gomock.Any()).Return(mockCipher).Times(1)
				mockFileReaderWriter := crypto.NewMockFileReaderWriter(ctrl)
				gomock.InOrder(
					mockFileReaderWriter.EXPECT().ReadFile("zypher.key").Return([]byte("key"), nil).Times(1),
					mockFileReaderWriter.EXPECT().ReadFile("input.txt").Return([]byte("content"), nil).Times(1),
				)
				mockFileReaderWriter.EXPECT().WriteFile("input.enc", gomock.Any(), fs.FileMode(0600)).Times(1)
				return mockCipherFactory, mockFileReaderWriter
			},
		},
		{
			name:            "run successfully with key from overridden another.key file",
			args:            []string{"-f", "input.txt", "-o", "input.enc", "-kf", "another.key"},
			expectedErrCode: 0,
			initMocks: func() (crypto.CipherFactory, crypto.FileReaderWriter) {
				mockCipherFactory := crypto.NewMockCipherFactory(ctrl)
				mockCipher := crypto.NewMockCipher(ctrl)
				mockCipher.EXPECT().Encrypt([]byte("content")).Times(1)
				mockCipherFactory.EXPECT().NewCipher(gomock.Any()).Return(mockCipher).Times(1)
				mockFileReaderWriter := crypto.NewMockFileReaderWriter(ctrl)
				gomock.InOrder(
					mockFileReaderWriter.EXPECT().ReadFile("another.key").Return([]byte("anotherkey"), nil).Times(1),
					mockFileReaderWriter.EXPECT().ReadFile("input.txt").Return([]byte("content"), nil).Times(1),
				)
				mockFileReaderWriter.EXPECT().WriteFile("input.enc", gomock.Any(), fs.FileMode(0600)).Times(1)
				return mockCipherFactory, mockFileReaderWriter
			},
		},
		{
			name:            "run successfully with key from ZYPHER_KEY env",
			args:            []string{"sometext"},
			envs:            map[string]string{"ZYPHER_KEY": "key"},
			expectedErrCode: 0,
			initMocks: func() (crypto.CipherFactory, crypto.FileReaderWriter) {
				mockCipherFactory := crypto.NewMockCipherFactory(ctrl)
				mockCipher := crypto.NewMockCipher(ctrl)
				mockCipher.EXPECT().Encrypt([]byte("sometext")).Times(1)
				mockFileReaderWriter := crypto.NewMockFileReaderWriter(ctrl)
				mockFileReaderWriter.EXPECT().ReadFile("zypher.key").Return(nil, errors.New("file not exist")).Times(1)
				mockCipherFactory.EXPECT().NewCipher(gomock.Any()).Return(mockCipher).Times(1)
				return mockCipherFactory, mockFileReaderWriter
			},
		},
		{
			name:            "fails when no key provided",
			args:            []string{"sometext"},
			expectedErrCode: 1,
			initMocks: func() (crypto.CipherFactory, crypto.FileReaderWriter) {
				mockCipherFactory := crypto.NewMockCipherFactory(ctrl)
				mockCipher := crypto.NewMockCipher(ctrl)
				mockCipher.EXPECT().Encrypt(gomock.Any()).Times(0)
				mockFileReaderWriter := crypto.NewMockFileReaderWriter(ctrl)
				mockFileReaderWriter.EXPECT().ReadFile("zypher.key").Return(nil, errors.New("file not exist")).Times(1)
				mockCipherFactory.EXPECT().NewCipher(gomock.Any()).Return(mockCipher).Times(0)
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
			encryptCmd := crypto.NewEncryptCmd(
				mockCipherFactory,
				crypto.WithFileReaderWriter(mockFileReaderWriter),
			)
			errCode := encryptCmd.Run(tt.args)
			if errCode != tt.expectedErrCode {
				t.Errorf("Expected code %d, got %d", tt.expectedErrCode, errCode)
			}
		})
	}
}
