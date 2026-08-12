package renderer

// RenderedFileKind identifies the type of filesystem object represented by a
// rendered file.
type RenderedFileKind int

const (
	RenderedFileKindRegular RenderedFileKind = iota
	RenderedFileKindSymlink
)

// RenderedFile represents one rendered managed configuration filesystem object.
//
// RelativePath is the path of the generated object relative to the root of the
// tenant's managed configuration tree.
//
// Regular files use Content. Symlinks use Target.
type RenderedFile struct {
	RelativePath string
	Kind         RenderedFileKind
	Content      string
	Target       string
}
