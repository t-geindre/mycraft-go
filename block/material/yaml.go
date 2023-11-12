package material

import "github.com/g3n/engine/math32"

type _YAMLMaterial struct {
	Color          *string
	Color3         *math32.Color
	SpecularColor  *string
	SpecularColor3 *math32.Color
	EmissiveColor  *string
	EmissiveColor3 *math32.Color
	Shininess      *float32
	Opacity        *float32
	Texture        *string
	Transparent    *bool
	A              *bool
}

type _YAMLMaterials map[string]_YAMLMaterial

func getYAMLMaterialColor(definition _YAMLMaterial) *math32.Color {
	if definition.Color != nil {
		return math32.NewColor(*definition.Color)
	}

	if definition.Color3 != nil {
		return definition.Color3
	}

	return &math32.Color{R: 1, G: 1, B: 1}
}

func getYAMLMaterialSpecularColor(definition _YAMLMaterial) *math32.Color {
	if definition.SpecularColor != nil {
		return math32.NewColor(*definition.SpecularColor)
	}

	if definition.SpecularColor3 != nil {
		return definition.SpecularColor3
	}

	return &math32.Color{R: 0.5, G: 0.5, B: 0.5}
}

func getYAMLMaterialEmissiveColor(definition _YAMLMaterial) *math32.Color {
	if definition.EmissiveColor != nil {
		return math32.NewColor(*definition.EmissiveColor)
	}

	if definition.EmissiveColor3 != nil {
		return definition.EmissiveColor3
	}

	return &math32.Color{R: 0, G: 0, B: 0}
}
