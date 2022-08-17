package types

type DenomTrait int

const (
	NotEditable DenomTrait = iota + 1
	ManageableBySuperAdmin
)

var DenomTraitsMap = map[string]DenomTrait{
	"NotEditable":            NotEditable,
	"ManageableBySuperAdmin": ManageableBySuperAdmin,
}
