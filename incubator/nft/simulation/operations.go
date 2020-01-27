package simulation

import (
	"fmt"
	"math/rand"

	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/simulation"
	"github.com/cosmos/modules/incubator/nft"
	"github.com/cosmos/modules/incubator/nft/internal/keeper"
	"github.com/cosmos/modules/incubator/nft/internal/types"
)

const (
	OpWeightedMsgTransferNFT     = " op_weighted_msg_transfer_nft"
	OpWeightedMsgEditNFTMetadata = "op_weighted_msg_edit_nft_metadata"
	OpWeightedMsgMintNFT         = "op_weighted_msg_mint_nft"
	OpWeightedMsgBurnNFT         = "op_weighted_msg_burn_nft"
)

func WeightedOperations(appParams simulation.AppParams, cdc *codec.Codec, k keeper.Keeper) simulation.WeightedOperations {

	var (
		weightedMsgTransferNFT     int
		weightedMsgEditNFTMetadata int
		weightedMsgMintNFT         int
		weightedMsgBurnNFT         int
	)

	appParams.GetOrGenerate(cdc, OpWeightedMsgTransferNFT, &weightedMsgTransferNFT, nil,
		func(_ *rand.Rand) {
			weightedMsgTransferNFT = 100
		},
	)
	appParams.GetOrGenerate(cdc, OpWeightedMsgEditNFTMetadata, &weightedMsgEditNFTMetadata, nil,
		func(_ *rand.Rand) {
			weightedMsgEditNFTMetadata = 67
		},
	)
	appParams.GetOrGenerate(cdc, OpWeightedMsgMintNFT, &weightedMsgMintNFT, nil,
		func(_ *rand.Rand) {
			weightedMsgMintNFT = 5
		},
	)
	appParams.GetOrGenerate(cdc, OpWeightedMsgBurnNFT, &weightedMsgBurnNFT, nil,
		func(_ *rand.Rand) {
			weightedMsgBurnNFT = 5
		},
	)

	return simulation.WeightedOperations{
		simulation.NewweightedOperation(
			weightedMsgTransferNFT,
			SimulateMsgTransferNFT(k),
		),
		simulation.NewweightedOperation(
			weightedMsgEditNFTMetadata,
			SimulateMsgEditNFTMetadata(k),
		),
		simulation.NewweightedOperation(
			weightedMsgMintNFT,
			SimulateMsgMintNFT(k),
		),
		simulation.NewweightedOperation(
			weightedMsgBurnNFT,
			SimulateMsgBurnNFT(k),
		),
	}
}

// SimulateMsgTransferNFT simulates the transfer of an NFT
func SimulateMsgTransferNFT(k keeper.Keeper) simulation.Operation {
	handler := nft.GenericHandler(k)
	return func(
		r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context, accs []simulation.Account, chainID string,
	) (simulation.OperationMsg, []simulation.FutureOperation, error) {
		ownerAddr, denom, nftID := getRandomNFTFromOwner(ctx, k, r)
		if ownerAddr.Empty() {
			return simulation.NoOpMsg(types.ModuleName), nil, nil
		}

		simAccount, _ := simulation.RandomAcc(r, accs)
		msg := types.NewMsgTransferNFT(
			ownerAddr,          // sender
			simAccount.Address, // recipient
			denom,
			nftID,
		)

		if msg.ValidateBasic() != nil {
			return simulation.NoOpMsg(types.ModuleName), nil, fmt.Errorf("expected msg to pass ValidateBasic: %s", msg.GetSignBytes())
		}

		ctx, _ = ctx.CacheContext()
		_, err := handler(ctx, msg)
		if err != nil {
			return simulation.NoOpMsg(types.ModuleName), nil, err
		}

		return simulation.NewOperationMsg(msg, true, ""), nil, nil
	}
}

// SimulateMsgEditNFTMetadata simulates an edit metadata transaction
func SimulateMsgEditNFTMetadata(k keeper.Keeper) simulation.Operation {
	handler := nft.GenericHandler(k)
	return func(
		r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context, accs []simulation.Account, chainID string,
	) (simulation.OperationMsg, []simulation.FutureOperation, error) {
		ownerAddr, denom, nftID := getRandomNFTFromOwner(ctx, k, r)
		if ownerAddr.Empty() {
			return simulation.NoOpMsg(types.ModuleName), nil, nil
		}

		msg := types.NewMsgEditNFTMetadata(
			ownerAddr,
			nftID,
			denom,
			simulation.RandStringOfLength(r, 45), // tokenURI
		)

		if msg.ValidateBasic() != nil {
			return simulation.NoOpMsg(types.ModuleName), nil, fmt.Errorf("expected msg to pass ValidateBasic: %s", msg.GetSignBytes())
		}

		ctx, _ = ctx.CacheContext()
		_, err := handler(ctx, msg)
		if err != nil {
			return simulation.NoOpMsg(types.ModuleName), nil, err
		}

		return simulation.NewOperationMsg(msg, true, ""), nil, nil
	}
}

// SimulateMsgMintNFT simulates a mint of an NFT
func SimulateMsgMintNFT(k keeper.Keeper) simulation.Operation {
	handler := nft.GenericHandler(k)
	return func(
		r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context, accs []simulation.Account, chainID string,
	) (simulation.OperationMsg, []simulation.FutureOperation, error) {
		a1, _ := simulation.RandomAcc(r, accs)
		a2, _ := simulation.RandomAcc(r, accs)

		msg := types.NewMsgMintNFT(
			a1.Address,                           // sender
			a2.Address,                           // recipient
			simulation.RandStringOfLength(r, 10), // nft ID
			simulation.RandStringOfLength(r, 10), // denom
			simulation.RandStringOfLength(r, 45), // tokenURI
		)

		if msg.ValidateBasic() != nil {
			return simulation.NoOpMsg(types.ModuleName), nil, fmt.Errorf("expected msg to pass ValidateBasic: %s", msg.GetSignBytes())
		}

		ctx, _ = ctx.CacheContext()
		_, err := handler(ctx, msg)
		if err != nil {
			return simulation.NoOpMsg(types.ModuleName), nil, err
		}

		return simulation.NewOperationMsg(msg, true, ""), nil, nil
	}
}

// SimulateMsgBurnNFT simulates a burn of an existing NFT
func SimulateMsgBurnNFT(k keeper.Keeper) simulation.Operation {
	handler := nft.GenericHandler(k)
	return func(
		r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context, accs []simulation.Account, chainID string,
	) (simulation.OperationMsg, []simulation.FutureOperation, error) {
		ownerAddr, denom, nftID := getRandomNFTFromOwner(ctx, k, r)
		if ownerAddr.Empty() {
			return simulation.NoOpMsg(types.ModuleName), nil, nil
		}

		msg := types.NewMsgBurnNFT(ownerAddr, nftID, denom)

		if msg.ValidateBasic() != nil {
			return simulation.NoOpMsg(types.ModuleName), nil, fmt.Errorf("expected msg to pass ValidateBasic: %s", msg.GetSignBytes())
		}

		ctx, _ = ctx.CacheContext()
		_, err := handler(ctx, msg)
		if err != nil {
			return simulation.NoOpMsg(types.ModuleName), nil, err
		}

		return simulation.NewOperationMsg(msg, true, ""), nil, nil
	}
}

func getRandomNFTFromOwner(ctx sdk.Context, k keeper.Keeper, r *rand.Rand) (address sdk.AccAddress, denom, nftID string) {
	owners := k.GetOwners(ctx)

	ownersLen := len(owners)
	if ownersLen == 0 {
		return nil, "", ""
	}

	// get random owner
	i := r.Intn(ownersLen)
	owner := owners[i]

	idCollectionsLen := len(owner.IDCollections)
	if idCollectionsLen == 0 {
		return nil, "", ""
	}

	// get random collection from owner's balance
	i = r.Intn(idCollectionsLen)
	idsCollection := owner.IDCollections[i] // nfts IDs
	denom = idsCollection.Denom

	idsLen := len(idsCollection.IDs)
	if idsLen == 0 {
		return nil, "", ""
	}

	// get random nft from collection
	i = r.Intn(idsLen)
	nftID = idsCollection.IDs[i]

	return owner.Address, denom, nftID
}
