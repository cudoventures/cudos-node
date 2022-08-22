package cli

import flag "github.com/spf13/pflag"

const (
	FlagMintRoyalties   = "mint-royalties"
	FlagResaleRoyalties = "resale-royalties"

	FlagMintNftName = "name"
	FlagMintNftUri  = "uri"
	FlagMintNftData = "data"
)

var FsPublishCollection = flag.NewFlagSet("", flag.ContinueOnError)
var FsMintNFT = flag.NewFlagSet("", flag.ContinueOnError)

func init() {
	FsPublishCollection.String(FlagMintRoyalties, "", "Royalties only for NFT first sale")
	FsPublishCollection.String(FlagResaleRoyalties, "", "Royalties for NFT resale after the first sale")

	FsMintNFT.String(FlagMintNftName, "", "NFT name")
	FsMintNFT.String(FlagMintNftUri, "", "NFT uri")
	FsMintNFT.String(FlagMintNftData, "", "NFT data")
}
