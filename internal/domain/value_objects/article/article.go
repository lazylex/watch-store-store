package article

const (
	CaseWithoutDefect = iota
	CaseWithScratches
	CaseWithHeavyScratches
)

const (
	DisplayWithoutDefects = iota
	DisplayWithScratches
	DisplayWithHeavyScratches
)

const (
	PackageNotOpened = iota
	PackageOpened
)

const (
	PackagingWithoutDamage = iota
	PackagingWithDamage
)

type Article string
