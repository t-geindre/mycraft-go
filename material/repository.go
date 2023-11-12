package material

import (
	"fmt"
	"github.com/g3n/engine/gls"
	"github.com/g3n/engine/material"
	"github.com/g3n/engine/texture"
	"mycraft/file"
)

type Repository struct {
	materials map[string]material.IMaterial
	textures  map[string]*texture.Texture2D
}

func (r *Repository) Get(id string) material.IMaterial {
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
		mat := material.NewStandard(getYAMLMaterialColor(def))
		mat.SetSpecularColor(getYAMLMaterialSpecularColor(def))
		mat.SetEmissiveColor(getYAMLMaterialEmissiveColor(def))

		if def.Opacity != nil {
			mat.SetOpacity(*def.Opacity)
		}

		mat.SetTransparent(def.Transparent)

		if def.Texture != nil {
			mat.AddTexture(r.getTexture(*def.Texture))
		}

		r.materials[id] = mat
	}
}

func NewFromYamlFile(filePath string) Repository {
	repository := Repository{
		map[string]material.IMaterial{},
		map[string]*texture.Texture2D{},
	}

	repository.AppendFromYamlFile(filePath)

	return repository
}
