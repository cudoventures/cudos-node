package types

// NFT module event types
var (
	EventTypeIssueDenom    = "issue_denom"
	EventTypeTransferNft   = "transfer_nft"
	EventTypeApproveNft    = "approve_nft"
	EventTypeApproveAllNft = "approve_all_nft"
	EventTypeRevokeNft     = "revoke_nft"
	EventTypeEditNFT       = "edit_nft"
	EventTypeMintNFT       = "mint_nft"
	EventTypeBurnNFT       = "burn_nft"

	AttributeValueCategory = ModuleName

	AttributeKeySender    = "sender"
	AttributeKeyCreator   = "creator"
	AttributeKeyRecipient = "recipient"
	AttributeKeyOwner     = "owner"
	AttributeKeyOperator  = "operator"
	AttributeKeyTokenID   = "token_id"
	AttributeKeyTokenURI  = "token_uri"
	AttributeKeyDenomID   = "denom_id"
	AttributeKeyDenomName = "denom_name"
	AttributeKeyMessage   = "message"
	AttributeKeyApproved  = "approved"
	AttributeKeyFrom      = "from"
	AttributeKeyTo        = "to"
)
