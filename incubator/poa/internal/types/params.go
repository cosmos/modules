package types

import (
	"bytes"
	"fmt"
	"time"

	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/x/params"
)

// POA params default values
const (
	// DefaultUnbondingTime reflects three weeks in seconds as the default
	// unbonding time.
	DefaultUnbondingTime time.Duration = time.Hour * 24 * 7 * 3

	// Default maximum number of bonded validators
	DefaultMaxValidators uint16 = 100

	// Default maximum entries in a UBD/RED pair
	DefaultMaxEntries uint16 = 7

	// Default to accept all validators
	DefaultAcceptAllValidators bool = true

	// Default to not allow validators to increase there weight
	DefaultIncreaseWeight bool = false

	// DefaultHistorical entries is 0 since it must only be non-zero for
	// IBC connected chains
	DefaultHistoricalEntries uint16 = 0
)

// nolint - Keys for parameter access
var (
	KeyUnbondingTime       = []byte("UnbondingTime")
	KeyMaxValidators       = []byte("MaxValidators")
	KeyMaxEntries          = []byte("KeyMaxEntries")
	KeyAcceptAllValidators = []byte("KeyAcceptAllValidators")
	KeyIncreaseWeight      = []byte("KeyIncreaseWeight")
	KeyHistoricalEntries   = []byte("HistoricalEntries")
)

var _ params.ParamSet = (*Params)(nil)

// Params defines the high level settings for staking
type Params struct {
	UnbondingTime       time.Duration `json:"unbonding_time" yaml:"unbonding_time"`               // time duration of unbonding
	MaxValidators       uint16        `json:"max_validators" yaml:"max_validators"`               // maximum number of validators (max uint16 = 65535)
	MaxEntries          uint16        `json:"max_entries" yaml:"max_entries"`                     // max entries for either unbonding delegation or redelegation (per pair/trio)
	HistoricalEntries   uint16        `json:"historical_entries" yaml:"historical_entries"`       // number of historical entries to persist
	AcceptAllValidators bool          `json:"accept_all_validators" yaml:"accept_all_validators"` // Sets the value if a network wants to accept all applicants to be validators
	IncreaseWeight      bool          `json:"increase_weight" yaml:"increase_weight"`             // Disallow validators to increase there power
}

// NewParams creates a new Params instance
func NewParams(unbondingTime time.Duration, maxValidators, maxEntries, historicalEntries uint16,
	acceptAllValidators, increaseWeight bool) Params {

	return Params{
		UnbondingTime:       unbondingTime,
		MaxValidators:       maxValidators,
		MaxEntries:          maxEntries,
		HistoricalEntries:   historicalEntries,
		AcceptAllValidators: acceptAllValidators,
		IncreaseWeight:      increaseWeight,
	}
}

// Implements params.ParamSet
func (p *Params) ParamSetPairs() params.ParamSetPairs {
	return params.ParamSetPairs{
		params.NewParamSetPair(KeyUnbondingTime, &p.UnbondingTime, validateUnbondingTime),
		params.NewParamSetPair(KeyMaxValidators, &p.MaxValidators, validateMaxValidators),
		params.NewParamSetPair(KeyMaxEntries, &p.MaxEntries, validateMaxEntries),
		params.NewParamSetPair(KeyAcceptAllValidators, &p.AcceptAllValidators, validateAcceptAllValidators),
		params.NewParamSetPair(KeyHistoricalEntries, &p.HistoricalEntries, validateHistoricalEntries),
		params.NewParamSetPair(KeyIncreaseWeight, &p.IncreaseWeight, validateKeyIncreaseWeight),
	}
}

// Equal returns a boolean determining if two Param types are identical.
// TODO: This is slower than comparing struct fields directly
func (p Params) Equal(p2 Params) bool {
	bz1 := ModuleCdc.MustMarshalBinaryLengthPrefixed(&p)
	bz2 := ModuleCdc.MustMarshalBinaryLengthPrefixed(&p2)
	return bytes.Equal(bz1, bz2)
}

// DefaultParams returns a default set of parameters.
func DefaultParams() Params {
	return NewParams(DefaultUnbondingTime, DefaultMaxValidators, DefaultMaxEntries, DefaultHistoricalEntries, DefaultAcceptAllValidators, DefaultIncreaseWeight)
}

// String returns a human readable string representation of the parameters.
func (p Params) String() string {
	return fmt.Sprintf(`Params:
  Unbonding Time:    %s
  Max Validators:    %d
	Max Entries:       %d
	Accept All Validators %t`, p.UnbondingTime,
		p.MaxValidators, p.MaxEntries, p.AcceptAllValidators)
}

// unmarshal the current staking params value from store key or panic
func MustUnmarshalParams(cdc *codec.Codec, value []byte) Params {
	params, err := UnmarshalParams(cdc, value)
	if err != nil {
		panic(err)
	}
	return params
}

// unmarshal the current staking params value from store key
func UnmarshalParams(cdc *codec.Codec, value []byte) (params Params, err error) {
	err = cdc.UnmarshalBinaryLengthPrefixed(value, &params)
	if err != nil {
		return
	}
	return
}

// validate a set of params
func (p Params) Validate() error {
	if err := validateUnbondingTime(p.UnbondingTime); err != nil {
		return err
	}
	if err := validateMaxValidators(p.MaxValidators); err != nil {
		return err
	}
	if err := validateMaxEntries(p.MaxEntries); err != nil {
		return err
	}
	return nil
}

func validateUnbondingTime(i interface{}) error {
	v, ok := i.(time.Duration)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if v <= 0 {
		return fmt.Errorf("unbonding time must be positive: %d", v)
	}

	return nil
}

func validateMaxValidators(i interface{}) error {
	v, ok := i.(uint16)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if v == 0 {
		return fmt.Errorf("max validators must be positive: %d", v)
	}

	return nil
}

func validateMaxEntries(i interface{}) error {
	v, ok := i.(uint16)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if v == 0 {
		return fmt.Errorf("max entries must be positive: %d", v)
	}

	return nil
}

func validateHistoricalEntries(i interface{}) error {
	_, ok := i.(uint16)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	return nil
}

func validateAcceptAllValidators(i interface{}) error {
	_, ok := i.(bool)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	return nil
}

func validateKeyIncreaseWeight(i interface{}) error {
	_, ok := i.(bool)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	return nil
}
