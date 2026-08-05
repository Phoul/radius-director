package renderer

// RenderedFile represents one rendered managed configuration file.
//
// RelativePath is the path of the generated file relative to the root of the
// tenant's managed configuration tree.
type RenderedFile struct {
	RelativePath string
	Content      string
}
