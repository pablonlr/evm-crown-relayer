package crown

import (
	"strings"

	rtypes "github.com/pablonlr/poly-crown-relayer/types"
)

const REQUIRED_PROTOCOL_CONFIRMATIONS = 6

const (
	ProtoNotFoundPrefix = "Can't find an NFT protocol"
	TokenNotFoundPrefix = "Can't find an NFT record"
)

func (crw *CrownResolver) GetNFToken(params ...string) rtypes.TaskResult {
	if len(params) < 2 {
		return rtypes.TaskResult{
			Err: rtypes.GetError(rtypes.MissingParams, nil),
		}
	}
	protocol := params[0]
	tokenID := params[1]
	crwToken, err := crw.client.GetNFToken(protocol, tokenID)
	if err != nil {
		if err.Code == -1 && strings.HasPrefix(err.Message, TokenNotFoundPrefix) {
			return rtypes.TaskResult{
				ResultValue: nil,
			}
		}
		return rtypes.TaskResult{
			Err: rtypes.GetError(rtypes.ErrorGetNFToken, err),
		}
	}
	return rtypes.TaskResult{
		ResultValue: crwToken,
	}

}

func (crw *CrownResolver) RegisterNFToken(params ...string) rtypes.TaskResult {
	if len(params) < 4 {
		return rtypes.TaskResult{
			Err: rtypes.GetError(rtypes.MissingParams, nil),
		}
	}
	crw.tryToUnclockWallet()
	protocol := params[0]
	tokenID := params[1]
	owner := params[2]
	uri := params[3]
	txId, err := crw.client.RegisterNFToken(protocol, tokenID, owner, owner, uri)
	if err != nil {
		return rtypes.TaskResult{
			Err: rtypes.GetError(rtypes.ErrorGetNFToken, err),
		}
	}
	return rtypes.TaskResult{
		ResultValue: txId,
	}

}

func (crw *CrownResolver) NFTokenConfirmed(params ...string) rtypes.TaskResult {
	if len(params) < 1 {
		return rtypes.TaskResult{
			Err: rtypes.GetError(rtypes.MissingParams, nil),
		}
	}
	txId := params[0]
	token, err := crw.client.GetNFTokenByTxID(txId)
	if err != nil {
		return rtypes.TaskResult{
			Err: rtypes.GetError(rtypes.ErrorGetNFTokenByTxID, err),
		}
	}

	var confirmations int
	if len(token.BlockHash) == 64 {
		confirmations = 1
	}

	return rtypes.TaskResult{
		ResultValue: confirmations,
	}

}

func (crw *CrownResolver) GetNFTProtocol(params ...string) rtypes.TaskResult {
	if len(params) < 1 {
		return rtypes.TaskResult{
			Err: rtypes.GetError(rtypes.MissingParams, nil),
		}
	}
	protocol := params[0]
	proto, err := crw.client.GetNFTProtocol(protocol)
	if err != nil {
		if err.Code == -1 && strings.HasPrefix(err.Message, ProtoNotFoundPrefix) {
			return rtypes.TaskResult{
				ResultValue: nil,
			}
		}
		return rtypes.TaskResult{
			Err: rtypes.GetError(rtypes.ErrorGetNFToken, err),
		}
	}
	return rtypes.TaskResult{
		ResultValue: proto,
	}

}

func (crw *CrownResolver) NFTProtocolConfirmed(params ...string) rtypes.TaskResult {
	if len(params) < 1 {
		return rtypes.TaskResult{
			Err: rtypes.GetError(rtypes.MissingParams, nil),
		}
	}
	protocol := params[0]
	proto, err := crw.client.GetNFTProtocol(protocol)
	if err != nil {
		return rtypes.TaskResult{
			Err: rtypes.GetError(rtypes.ErrorGetNFTProtocol, err),
		}
	}
	block, err := crw.client.GetBlock(proto.BlockHash)
	if err != nil {
		return rtypes.TaskResult{
			Err: rtypes.GetError(rtypes.ErrorGetBlock, err),
		}
	}
	return rtypes.TaskResult{
		ResultValue: block.Confirmations,
	}

}

func (crw *CrownResolver) RegisterNFTProtocol(params ...string) rtypes.TaskResult {
	if len(params) < 4 {
		return rtypes.TaskResult{
			Err: rtypes.GetError(rtypes.MissingParams, nil),
		}
	}
	crw.tryToUnclockWallet()
	protocolID := params[0]
	protocolName := params[1]
	owner := params[2]
	description := params[3]
	txId, err := crw.client.RegisterNFTProtocol(protocolID, protocolName, owner, 2, "application/json", description, false, false, 255)
	if err != nil {
		return rtypes.TaskResult{
			Err: rtypes.GetError(rtypes.ErrorRegisterNFTProtocol, err),
		}
	}
	return rtypes.TaskResult{
		ResultValue: txId,
	}

}
