package types

import (
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

var (
	InvalidCollection = sdkerrors.Register(ModuleName, 1, "invalid NFT collection")
	UnknownCollection = sdkerrors.Register(ModuleName, 2, "unknown NFT collection")
	InvalidNFT        = sdkerrors.Register(ModuleName, 3, "invalid NFT")
	UnknownNFT        = sdkerrors.Register(ModuleName, 4, "unknown NFT")
	NFTAlreadyExists  = sdkerrors.Register(ModuleName, 5, "NFT already exists")
	EmptyMetadata     = sdkerrors.Register(ModuleName, 6, "NFT metadata can't be empty")
)
