package types

import codectypes "github.com/cosmos/cosmos-sdk/codec/types"

var (
	_ codectypes.UnpackInterfacesMessage = QueryGetAuctionResponse{}
	_ codectypes.UnpackInterfacesMessage = QueryAllAuctionResponse{}
)

func (msg QueryGetAuctionResponse) UnpackInterfaces(unpacker codectypes.AnyUnpacker) error {
	var a Auction
	return unpacker.UnpackAny(msg.Auction, &a)
}

func (msg QueryAllAuctionResponse) UnpackInterfaces(unpacker codectypes.AnyUnpacker) error {
	for _, any := range msg.Auctions {
		var a Auction
		err := unpacker.UnpackAny(any, &a)
		if err != nil {
			return err
		}
	}
	return nil
}
