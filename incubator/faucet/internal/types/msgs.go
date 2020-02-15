package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// RouterKey is the module name router key
const RouterKey = ModuleName // this was defined in your key.go file

// MsgMint defines a mint message
type MsgMint struct {
	Sender sdk.AccAddress
	Minter sdk.AccAddress
}

// NewMsgMint is a constructor function for NewMsgMint
func NewMsgMint(sender sdk.AccAddress, minter sdk.AccAddress) MsgMint {
	return MsgMint{Sender: sender, Minter: minter}
}

// Route should return the name of the module
func (msg MsgMint) Route() string { return RouterKey }

// Type should return the action
func (msg MsgMint) Type() string { return "mint" }

// ValidateBasic runs stateless checks on the message
func (msg MsgMint) ValidateBasic() error {
	if msg.Minter.Empty() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, msg.Minter.String())
	}
	if msg.Sender.Empty() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, msg.Sender.String())
	}
	return nil
}

// GetSignBytes encodes the message for signing
func (msg MsgMint) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

// GetSigners defines whose signature is required
func (msg MsgMint) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Sender}
}
