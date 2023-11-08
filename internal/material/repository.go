package material

import (
	"fmt"
	"mycraft/internal/file"

	"github.com/g3n/engine/material"
	"github.com/g3n/engine/math32"
	"github.com/g3n/engine/texture"
)

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
}

type _YAMLMaterials map[string]_YAMLMaterial

type _Repository struct {
	materials map[string]*material.Standard
	textures  map[string]*texture.Texture2D
}

func (r _Repository) Get(name string) *material.Standard {
	return r.materials[name]
}

func (r *_Repository) getTexture(file string) *texture.Texture2D {
	fmt.Printf("Loading \"%s\"... ", file)

	if texture, ok := r.textures[file]; ok {
		fmt.Println("DONE(CACHED)")
		return texture
	}

	texture, err := texture.NewTexture2DFromImage(file)
	if err != nil {
		panic(err)
	}

	r.textures[file] = texture
	fmt.Println("DONE")

	return r.textures[file]
}

func (r *_Repository) AppendFromYamlFile(filePath string) {

	rawYaml := file.ParseYamlFile[_YAMLMaterials](filePath)

	for id, def := range rawYaml {
		material := material.NewStandard(getMaterialColor(def))
		material.SetSpecularColor(getMaterialSpecularColor(def))
		material.SetEmissiveColor(getMaterialEmissiveColor(def))

		if def.Texture != nil {
			material.AddTexture(r.getTexture(*def.Texture))
		}

		r.materials[id] = material
	}
}

func CreateFromYamlFile(filePath string) _Repository {
	repository := _Repository{
		map[string]*material.Standard{},
		map[string]*texture.Texture2D{},
	}

	repository.AppendFromYamlFile(filePath)

	return repository
}

func getMaterialColor(definition _YAMLMaterial) *math32.Color {
	if definition.Color != nil {
		return math32.NewColor(*definition.Color)
	}

	if definition.Color3 != nil {
		return definition.Color3
	}

	return &math32.Color{R: 1, G: 1, B: 1}
}

func getMaterialSpecularColor(definition _YAMLMaterial) *math32.Color {
	if definition.SpecularColor != nil {
		return math32.NewColor(*definition.SpecularColor)
	}

	if definition.SpecularColor3 != nil {
		return definition.SpecularColor3
	}

	return &math32.Color{R: 0.5, G: 0.5, B: 0.5}
}

func getMaterialEmissiveColor(definition _YAMLMaterial) *math32.Color {
	if definition.EmissiveColor != nil {
		return math32.NewColor(*definition.EmissiveColor)
	}

	if definition.EmissiveColor3 != nil {
		return definition.EmissiveColor3
	}

	return &math32.Color{R: 0, G: 0, B: 0}
}
