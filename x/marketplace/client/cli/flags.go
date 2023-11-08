package cli

import flag "github.com/spf13/pflag"

const (
	FlagPublishCollectionMintRoyalties   = "mint-royalties"
	FlagPublishCollectionResaleRoyalties = "resale-royalties"

	FlagMintNftName = "name"
	FlagMintNftUri  = "uri"
	FlagMintNftData = "data"
	FlagMintNftUid  = "uid"

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

	FlagUpdateMintRoyalties   = "mint-royalties"
	FlagUpdateResaleRoyalties = "resale-royalties"
)

var (
	FsPublishCollection = flag.NewFlagSet("", flag.ContinueOnError)
	FsMintNFT           = flag.NewFlagSet("", flag.ContinueOnError)
	FsCreateCollection  = flag.NewFlagSet("", flag.ContinueOnError)
	FsUpdateRoyalties   = flag.NewFlagSet("", flag.ContinueOnError)
)

func init() {
	FsPublishCollection.String(FlagPublishCollectionMintRoyalties, "", "Collection mint royalties")
	FsPublishCollection.String(FlagPublishCollectionResaleRoyalties, "", "Collection resale royalties")

	FsMintNFT.String(FlagMintNftName, "", "NFT name")
	FsMintNFT.String(FlagMintNftUri, "", "NFT uri")
	FsMintNFT.String(FlagMintNftData, "", "NFT data")
	FsMintNFT.String(FlagMintNftUid, "", "NFT Uid")

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

	FsUpdateRoyalties.String(FlagUpdateMintRoyalties, "", "Collection mint royalties")
	FsUpdateRoyalties.String(FlagUpdateResaleRoyalties, "", "Collection resale royalties")
}
