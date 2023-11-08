package block

type _YAMLMaterial struct {
	Sides  string
	All    string
	Front  string
	Left   string
	Right  string
	Back   string
	Top    string
	Bottom string
}

type _YAMLBlock struct {
	Materials _YAMLMaterial
}

type _YAMLBlocks map[string]_YAMLBlock
