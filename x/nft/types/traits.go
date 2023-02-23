package types

type DenomTrait int

const (
	NotEditable DenomTrait = iota + 1
)

var DenomTraitsMapStrToType = map[string]DenomTrait{
	"NotEditable": NotEditable,
}
