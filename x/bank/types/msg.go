package types

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (m *MsgSend) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(m.FromAddress); err != nil {
		return fmt.Errorf("invalid from address: %w", err)
	}
	if _, err := sdk.AccAddressFromBech32(m.ToAddress); err != nil {
		return fmt.Errorf("invalid to address: %w", err)
	}
	if m.Amount == "0" {
		return fmt.Errorf("amount must be greater than zero")
	}
	return nil
}
