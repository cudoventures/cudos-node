package cli

import flag "github.com/spf13/pflag"

const (
	FlagMintRoyalties   = "mint-royalties"
	FlagResaleRoyalties = "resale-royalties"
)

var FsPublishCollection = flag.NewFlagSet("", flag.ContinueOnError)

func init() {
	FsPublishCollection.String(FlagMintRoyalties, "", "Royalties only for NFT first sale")
	FsPublishCollection.String(FlagResaleRoyalties, "", "Royalties for NFT resale after the first sale")
}
