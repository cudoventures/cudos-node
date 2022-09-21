package types

var (
	EventPublishCollectionType       = "publish_collection"
	EventPublishNftType              = "publish_nft"
	EventBuyNftType                  = "buy_nft"
	EventMintNftType                 = "marketplace_mint_nft"
	EventRemoveNftType               = "remove_nft"
	EventTransferAdminPermissionType = "transfer_admin_permission"
	EventVerifyCollectionType        = "verify_collection"
	EventUnverifyCollectionType      = "unverify_collection"

	AttributeValueCategory = ModuleName

	AttributeKeyCollectionID = "collection_id"
	AttributeKeyDenomID      = "denom_id"
	AttributeKeyCreator      = "creator"
	AttributeKeyTokenID      = "token_id"
	AttributeKeyNftID        = "nft_id"
	AttributeKeyPrice        = "price"
	AttributeKeyBuyer        = "buyer"
	AttributeKeyOwner        = "owner"
	AttributeKeyNewAdmin     = "new_admin"
	AttributeKeyUID          = "uid"
)
