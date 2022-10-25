package cli

import flag "github.com/spf13/pflag"

const (
	FlagMintRoyalties   = "mint-royalties"
	FlagResaleRoyalties = "resale-royalties"

	FlagMintNftName = "name"
	FlagMintNftUri  = "uri"
	FlagMintNftData = "data"

	FlagCreateCollectionName            = "name"
	FlagCreateCollectionSymbol          = "symbol"
	FlagCreateCollectionSchema          = "schema"
	FlagCreateCollectionTraits          = "traits"
	FlagCreateCollectionDescription     = "description"
	FlagCreateCollectionMinter          = "minter"
	FlagCreateCollectionData            = "data"
	FlagCreateCollectionMintRoyalties   = "mint-royalties"
	FlagCreateCollectionResaleRoyalties = "resale-royalties"
	FlagCreateCollectionVerified        = "verified"
)

var FsPublishCollection = flag.NewFlagSet("", flag.ContinueOnError)
var FsMintNFT = flag.NewFlagSet("", flag.ContinueOnError)
var FsCreateCollection = flag.NewFlagSet("", flag.ContinueOnError)

func init() {
	FsPublishCollection.String(FlagMintRoyalties, "", "Royalties only for NFT first sale")
	FsPublishCollection.String(FlagResaleRoyalties, "", "Royalties for NFT resale after the first sale")

	FsMintNFT.String(FlagMintNftName, "", "NFT name")
	FsMintNFT.String(FlagMintNftUri, "", "NFT uri")
	FsMintNFT.String(FlagMintNftData, "", "NFT data")

	FsCreateCollection.String(FlagCreateCollectionName, "", "Denom name")
	FsCreateCollection.String(FlagCreateCollectionSymbol, "", "Denom symbol name")
	FsCreateCollection.String(FlagCreateCollectionSchema, "", "Denom schema")
	FsCreateCollection.String(FlagCreateCollectionTraits, "", "Denom traits")
	FsCreateCollection.String(FlagCreateCollectionDescription, "", "Denom description")
	FsCreateCollection.String(FlagCreateCollectionMinter, "", "Denom minter")
	FsCreateCollection.String(FlagCreateCollectionData, "", "Denom data")
	FsCreateCollection.String(FlagCreateCollectionMintRoyalties, "", "Collection mint royalties")
	FsCreateCollection.String(FlagCreateCollectionResaleRoyalties, "", "Collection resale royalties")
	FsCreateCollection.String(FlagCreateCollectionVerified, "", "Collection verified flag")
}
