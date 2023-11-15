package block

type _YAMLMaterial struct {
	Sides string
	All   string
	South string
	West  string
	Est   string
	North string
	Up    string
	Down  string
}

type _YAMLBlock struct {
	Materials _YAMLMaterial
	Type      string
}

type _YAMLBlocks map[string]_YAMLBlock
