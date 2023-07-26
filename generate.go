package zypher

//go:generate mockgen -source=internal/crypto/base.go -destination=internal/crypto/mock.go -package=crypto
//go:generate mockgen -source=internal/server/server.go -destination=internal/server/mock.go -package=server
//go:generate mockgen -source=internal/server/store/store.go -destination=internal/server/store/mock.go -package=store
//go:generate mockgen -source=internal/server/auth/authentication.go -destination=internal/server/auth/mock.go -package=auth
