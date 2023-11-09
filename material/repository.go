package material

import (
	"fmt"
	"github.com/g3n/engine/gls"
	"github.com/g3n/engine/material"
	"github.com/g3n/engine/texture"
	"mycraft/file"
)

type Repository struct {
	materials map[string]*material.Standard
	textures  map[string]*texture.Texture2D
}

func (r *Repository) Get(id string) *material.Standard {
	if mat, ok := r.materials[id]; ok {
		return mat
	}
	panic(fmt.Errorf(`unknown material "%s"`, id))
}

func (r *Repository) getTexture(file string) *texture.Texture2D {

	if texture, ok := r.textures[file]; ok {
		return texture
	}

	texture, err := texture.NewTexture2DFromImage(file)
	if err != nil {
		panic(err)
	}

	texture.SetMagFilter(gls.NEAREST) // avoid blurry upscale

	r.textures[file] = texture

	return r.textures[file]
}

func (r *Repository) AppendFromYamlFile(filePath string) {
	rawYaml := file.ParseYamlFile[_YAMLMaterials](filePath)

	for id, def := range rawYaml {
		material := material.NewStandard(getYAMLMaterialColor(def))
		material.SetSpecularColor(getYAMLMaterialSpecularColor(def))
		material.SetEmissiveColor(getYAMLMaterialEmissiveColor(def))

		if def.Texture != nil {
			material.AddTexture(r.getTexture(*def.Texture))
		}

		r.materials[id] = material
	}
}

func NewFromYamlFile(filePath string) Repository {
	repository := Repository{
		map[string]*material.Standard{},
		map[string]*texture.Texture2D{},
	}

	repository.AppendFromYamlFile(filePath)

	return repository
}
