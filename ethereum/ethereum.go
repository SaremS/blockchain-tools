package ethereum

import (
	"context"
	"fmt"
	"strings"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

func GetChainData[T any](client *ethclient.Client, targetAddress, targetInterface, targetMethod string, writeContainer *T) error {
	targetAddressUnhexed := common.HexToAddress(targetAddress)
	parsedABI, err := abi.JSON(strings.NewReader(targetInterface))
	if err != nil {
		return fmt.Errorf("error while parsing ABI: %w", err)
	}

	callData, err := parsedABI.Pack(targetMethod)
	if err != nil {
		return fmt.Errorf("failed to pack slot0 call data: %v", err)
	}

	ctx := context.Background()
	msg := ethereum.CallMsg{To: &targetAddressUnhexed, Data: callData}
	result, err := client.CallContract(ctx, msg, nil)
	if err != nil {
		return fmt.Errorf("failed to call contract: %v", err)
	}

	err = parsedABI.UnpackIntoInterface(writeContainer, targetMethod, result)
	if err != nil {
		return fmt.Errorf("failed to unpack data: %v", err)
	}

	return nil
}
