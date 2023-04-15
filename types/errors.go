package types

const (
	MissingParams ErrorCode = iota + 1001
	InvalidParams
	ErrorGetNFToken
	ErrorGetNFTokenByTxID
	ErrorRegisterNFToken
	ErrorGetNFTProtocol
	ErrorGetBlock
	ErrorRegisterNFTProtocol
	InvalidCrownAddress
)

const (
	InvalidWaitTime ErrorCode = iota + 2001
)

var errMap = map[ErrorCode]string{
	MissingParams:            "MissingParams",
	InvalidParams:            "InvalidParams",
	ErrorGetNFToken:          "ErrorGetNFToken",
	ErrorGetNFTokenByTxID:    "ErrorGetNFTokenByTxID",
	ErrorRegisterNFToken:     "ErrorRegisterNFToken",
	ErrorGetNFTProtocol:      "ErrorGetNFTProtocol",
	ErrorGetBlock:            "ErrorGetBlock",
	ErrorRegisterNFTProtocol: "ErrorRegisterNFTProtocol",
	InvalidWaitTime:          "InvalidWaitTime",
	InvalidCrownAddress:      "InvalidCrownAddress",
}

func GetError(errCode ErrorCode, err error) Error {
	return Error{
		Code: errCode,
		Name: errMap[errCode],
		Err:  err,
		Skip: false,
	}
}

func GetSkipError(errCode ErrorCode, err error) Error {
	return Error{
		Code: errCode,
		Name: errMap[errCode],
		Err:  err,
		Skip: true,
	}
}
