package auth_test

import (
	"testing"

	"github.com/vtno/zypher/internal/auth"
	"github.com/vtno/zypher/internal/store"
	"go.uber.org/mock/gomock"
)

func TestAuthGuard_Guard(t *testing.T) {
	ctrl := gomock.NewController(t)

	type test struct {
		name           string
		initMocks      func() store.Store
		args           []string
		expectedResult bool
	}

	tests := []test{
		{
			name: "should return false when authHeader is empty",
			initMocks: func() store.Store {
				return store.NewMockStore(ctrl)
			},
			args:           []string{""},
			expectedResult: false,
		},
		{
			name: "should return false when token is not found in store",
			initMocks: func() store.Store {
				mStore := store.NewMockStore(ctrl)
				mStore.EXPECT().GetByBucket("auth", "token").Return("", nil)
				return mStore
			},
			args:           []string{"Bearer token"},
			expectedResult: false,
		},
		{
			name: "should return true when token is found in store",
			initMocks: func() store.Store {
				mStore := store.NewMockStore(ctrl)
				mStore.EXPECT().GetByBucket("auth", "token").Return("key", nil)
				return mStore
			},
			args:           []string{"Bearer token"},
			expectedResult: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mStore := tt.initMocks()
			authGuard := auth.NewAuthGuard(mStore)
			result := authGuard.Guard(tt.args[0])
			if result != tt.expectedResult {
				t.Errorf("Expected %v, got %v", tt.expectedResult, result)
			}
		})
	}
}
