# Faucet Module

This module will enable mint function. Every address can mint 100(bonded tokens) in every 24 hours by sending a mint message. 

For security consideration, you can add this module to your project as you want, but this module would *NOT* be active by default. unless you active it manually by adding `"-tags faucet"` when you build or install. 

## Usage

Step 1: import to your app.go
```go
import (
	"github.com/cosmos/modules/incubator/faucet"
)
```

Step 2: add module and permission
```go

	ModuleBasics = module.NewBasicManager(
		..., // the official basic modules

		faucet.AppModule{},
	)
	// account permissions
	maccPerms = map[string][]string{
		staking.BondedPoolName:    {supply.Burner, supply.Staking},
		staking.NotBondedPoolName: {supply.Burner, supply.Staking},
		faucet.ModuleName          {supply.Minter}, // add permissions for faucet
	}

```

Step 3: add to module manager
```go

	app.faucetKeeper = faucet.NewKeeper(app.supplyKeeper, app.stakingKeeper, keys[faucet.StoreKey], app.cdc,)

	app.mm = module.NewManager(
		genutil.NewAppModule(app.accountKeeper, app.stakingKeeper, app.BaseApp.DeliverTx),
		auth.NewAppModule(app.accountKeeper),
		bank.NewAppModule(app.bankKeeper, app.accountKeeper),
		nameservice.NewAppModule(app.nsKeeper, app.bankKeeper),
		
		faucet.NewAppModule(app.faucetKeeper), // faucet module
		
		supply.NewAppModule(app.supplyKeeper, app.accountKeeper),
		distr.NewAppModule(app.distrKeeper, app.accountKeeper, app.supplyKeeper, app.stakingKeeper),
		slashing.NewAppModule(app.slashingKeeper, app.accountKeeper, app.stakingKeeper),
		staking.NewAppModule(app.stakingKeeper, app.accountKeeper, app.supplyKeeper),
	)
```

Step 4: enable faucet in Makefile
```
installWithFaucet: go.sum
		go install -mod=readonly $(BUILD_FLAGS) -tags faucet ./cmd/nsd
		go install -mod=readonly $(BUILD_FLAGS) -tags faucet ./cmd/nscli
```

Step 5: build your app
```
make installWithFaucet
```

Stop 6: mint some coins
``` 
iMac:~ liangping$ nscli tx faucet mint --from ping --chain-id test -y
{
  "height": "0",
  "txhash": "A2AD3CA07949FD41CDFCD01678AE9E97946F7DD36DB008F16932C9E231CCFF85",
  "raw_log": "[]"
}
iMac:~ liangping$ nscli query account cosmos1ww6g4pdr3nzlyw7d2zcndx4jkrugkjucskvgsl --chain-id test 
{
  "type": "cosmos-sdk/Account",
  "value": {
    "address": "cosmos1ww6g4pdr3nzlyw7d2zcndx4jkrugkjucskvgsl",
    "coins": [
      {
        "denom": "stake",
        "amount": "11100000000"
      }
    ],
    "public_key": "cosmospub1addwnpepqwqcupvu4z6euqzja7hx354cky3h2vl8ht247rj92n3x3d86l5zlzpslzjx",
    "account_number": 3,
    "sequence": 20
  }
}

```


## Contact Us

Author: liangping from [Ping.pub](https://ping.pub)

18786721@qq.com

If you like this module, welcome to delegate to Ping.pub on [Cosmoshub](https://cosmos.ping.pub), [IRISHub](https://iris.ping.pub), [KAVA](https://kava.ping.pub).
