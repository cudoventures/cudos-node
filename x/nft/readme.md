# CLI
## Queries
1. `approvals` - Get the approved addresses for the NFT.
2. `collection` - Get all the NFTs from a given collection.
3. `denom` - Query the denom by the specified denom id.
4. `denom-by-name` - Query the denom by the specified denom name.
5. `denom-by-symbol` - Query the denom by the specified symbol.
6. `denoms` - Query all denominations of all collections of NFTs.
7. `is-approved-for-all` - Query if an address is an authorized operator for another address
8. `owner` - Get the NFTs owned by an account address.
9. `supply` - Query total supply of a collection or owner of NFTs.
10. `token` - Query a single NFT from a collection by denom id and token id.

## TXs
1. `approve` - Adds an address to the approved list of a NFT.
2. `approve-all` - Adds operator address to the globally approved list of sender.
3. `burn` - Burn an NFT.
4. `edit` - Edit the token data of an NFT.
5. `issue` - Issue a new denom.
6. `mint` - Mint a NFT and set the owner to the recipient. Only the denom creator can mint a new NFT.
7. `revoke` - Revokes a previously granted permission to transfer the given an NFT.
8. `transfer` - Transfer a NFT to a recipient.

# REST
1. `/nft/denoms/{denomId}` - `GET` - query a denom by denom id
2. `/nft/denoms/name/{denomName}` - `GET` - query a denom by denom name
3. `/nft/denoms/symbol/{denomSymbol}` - `GET` - query a denom by denom symbol
4. `/nft/denoms` - `POST` - query all denoms 
```
type queryDenomsRequest struct {
	Pagination query.PageRequest `json:"pagination"`
}
```
5. `/nft/collections` - `POST` - Get all the NFTs from a given collection
```
type queryCollectionRequest struct {
	DenomId    string            `json:"denom_id"`
	Pagination query.PageRequest `json:"pagination"`
}
```
6. `/nft/collections/supply/{denomId}` - `GET` - Get the total supply of a collection or owner
7. `/nft/owners` - `POST` - Get the collections of NFTs owned by an address
```
type queryOwnerRequest struct {
	DenomId      string            `json:"denom_id"`
	OwnerAddress string            `json:"owner_address"`
	Pagination   query.PageRequest `json:"pagination"`
}
```
8. `/nft/nfts/{denomId}/{tokenId}` - `GET` - Query a single NFT
9. `/nft/approvals/{denomId}/{tokenId}` - `GET` - Query approvals for NFT
10. `/nft/is-approved-for-all` - `POST` - Query if an address is an authorized operator for another address
```
type queryIsApprovedForAllRequest struct {
	Owner    string `json:"owner"`
	Operator string `json:"operator"`
}
```

# MSGS
1. Issue denom
```
typeUrl = "/cudosnode.cudosnode.nft.MsgIssueDenom"

message MsgIssueDenom {
    option (gogoproto.equal) = true;

    string id = 1;
    string name = 2;
    string schema = 3;
    string sender = 4;
    string contractAddressSigner = 5;
    string symbol = 6;
}
```

2. Transfer NFT to a receipeint
```
typeUrl = "/cudosnode.cudosnode.nft.MsgTransferNft"

message MsgTransferNft {
    option (gogoproto.equal) = true;

    string denom_id = 1 [ (gogoproto.moretags) = "yaml:\"denom_id\"" ];
    string token_id = 2;
    string from = 3;
    string to = 4;
    string sender = 5;
    string contractAddressSigner = 6;
}
```

3. Grant approval for a denom to an operator
```
typeUrl = "/cudosnode.cudosnode.nft.MsgApproveNft"

message MsgApproveNft {
    option (gogoproto.equal) = true;

    string id = 1;
    string denom_id = 2 [ (gogoproto.moretags) = "yaml:\"denom_id\"" ];
    string sender = 3;
    string approvedAddress = 4;
    string contractAddressSigner = 5;
}
```

4. Adds an adress to the globally approved list
```
typeUrl = "/cudosnode.cudosnode.nft.MsgApproveAllNft"

message MsgApproveAllNft {
    option (gogoproto.equal) = true;

    string  operator = 1;
    string  sender = 2;
    bool  approved = 3;
    string contractAddressSigner = 4;

}
```

5. Revokes a previously granted permission to transfer the given an NFT.
```
typeUrl = "/cudosnode.cudosnode.nft.MsgRevokeNft"

message MsgRevokeNft {
    option (gogoproto.equal) = true;

    string addressToRevoke = 1;
    string denom_id = 2 [ (gogoproto.moretags) = "yaml:\"denom_id\"" ];
    string token_id = 3 [ (gogoproto.moretags) = "yaml:\"denom_id\"" ];
    string sender = 4;
    string contractAddressSigner = 5;
}
```

6. Edits a denom
```
typeUrl = "/cudosnode.cudosnode.nft.MsgEditNFT"

message MsgEditNFT {
    option (gogoproto.equal) = true;

    string id = 1;
    string denom_id = 2 [ (gogoproto.moretags) = "yaml:\"denom_id\"" ];
    string name = 3;
    string uri = 4 [ (gogoproto.customname) = "URI" ];
    string data = 5;
    string sender = 6;
    string contractAddressSigner = 7;

}
```

7. Mints a new NFT.
```
typeUrl = "/cudosnode.cudosnode.nft.MsgMintNFT"

message MsgMintNFT {
    option (gogoproto.equal) = true;

    string denom_id = 1 [ (gogoproto.moretags) = "yaml:\"denom_id\"" ];
    string name = 2;
    string uri = 3 [ (gogoproto.customname) = "URI" ];
    string data = 4;
    string sender = 5;
    string recipient = 6;
    string contractAddressSigner = 7;
}
```

8. Burns a NFT.
```
typeUrl = "/cudosnode.cudosnode.nft.MsgBurnNFT"

message MsgBurnNFT {
    option (gogoproto.equal) = true;

    string id = 1;
    string denom_id = 2 [ (gogoproto.moretags) = "yaml:\"denom_id\"" ];
    string sender = 3;
    string contractAddressSigner = 4;
}
```