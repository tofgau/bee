package main

import (
	"net/http"
)

func returnCode200(w http.ResponseWriter, r *http.Request, text string) {
	// see http://golang.org/pkg/net/http/#pkg-constants
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
	w.Write([]byte(text))
}

func returnCode500(w http.ResponseWriter, r *http.Request, text string) {
	// see http://golang.org/pkg/net/http/#pkg-constants
	w.WriteHeader(http.StatusInternalServerError)
	w.Write([]byte(text))
}

func returnCode400(w http.ResponseWriter, r *http.Request, text string) {
	// see http://golang.org/pkg/net/http/#pkg-constants
	w.WriteHeader(http.StatusBadRequest)
	w.Write([]byte(text))
}

func returnCode401(w http.ResponseWriter, r *http.Request, text string) {
	// see http://golang.org/pkg/net/http/#pkg-constants
	w.WriteHeader(http.StatusUnauthorized)
	w.Write([]byte(text))
}

/*
 StatusContinue           = 100
     StatusSwitchingProtocols = 101

     StatusOK                   = 200
     StatusCreated              = 201
     StatusAccepted             = 202
     StatusNonAuthoritativeInfo = 203
     StatusNoContent            = 204
     StatusResetContent         = 205
     StatusPartialContent       = 206

     StatusMultipleChoices   = 300
     StatusMovedPermanently  = 301
     StatusFound             = 302
     StatusSeeOther          = 303
     StatusNotModified       = 304
     StatusUseProxy          = 305
     StatusTemporaryRedirect = 307

     StatusBadRequest                   = 400
     StatusUnauthorized                 = 401
     StatusPaymentRequired              = 402
     StatusForbidden                    = 403
     StatusNotFound                     = 404
     StatusMethodNotAllowed             = 405
     StatusNotAcceptable                = 406
     StatusProxyAuthRequired            = 407
     StatusRequestTimeout               = 408
     StatusConflict                     = 409
     StatusGone                         = 410
     StatusLengthRequired               = 411
     StatusPreconditionFailed           = 412
     StatusRequestEntityTooLarge        = 413
     StatusRequestURITooLong            = 414
     StatusUnsupportedMediaType         = 415
     StatusRequestedRangeNotSatisfiable = 416
     StatusExpectationFailed            = 417
     StatusTeapot                       = 418

     StatusInternalServerError     = 500
     StatusNotImplemented          = 501
     StatusBadGateway              = 502
     StatusServiceUnavailable      = 503
     StatusGatewayTimeout          = 504
     StatusHTTPVersionNotSupported = 505

*/