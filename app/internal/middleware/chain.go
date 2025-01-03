package middleware

import (
	"net/http"

	"github.com/justinas/alice"
)

var (
	ClearChain  = alice.New(LoggerMiddleWare)
	secureChain = alice.New(LoggerMiddleWare, JWtMiddleWare)
)

func CreateCleanChain(h http.Handler) http.Handler {
	return alice.New(LoggerMiddleWare).Then(h)
}

func CreateSecureChain(h http.Handler) http.Handler {
	return alice.New(LoggerMiddleWare, JWtMiddleWare).Then(h)
}

func CreateSecureChainForFunc(h http.HandlerFunc) http.Handler {
	return alice.New(LoggerMiddleWare, JWtMiddleWare).Then(h)
}

func CreateCleanChainForFunc(h http.HandlerFunc) http.Handler {
	return alice.New(LoggerMiddleWare).Then(h)
}
