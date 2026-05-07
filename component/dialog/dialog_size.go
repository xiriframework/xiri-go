package dialog

// Size is a typed dialog size token. Maps to a width on the Angular frontend.
type Size string

const (
	SizeSm   Size = "sm"
	SizeMd   Size = "md"
	SizeLg   Size = "lg"
	SizeXl   Size = "xl"
	SizeFull Size = "full"
)
