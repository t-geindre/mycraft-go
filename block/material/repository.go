package material

import (
	"fmt"
	"github.com/g3n/engine/gls"
	"github.com/g3n/engine/material"
	"github.com/g3n/engine/texture"
	"mycraft/block/file"
)

type Repository struct {
	materials map[string]material.IMaterial
	textures  map[string]*texture.Texture2D
}

const DefinitionFile = "assets/materials.yaml"

var RepositoryInstance *Repository

func (r *Repository) Get(id string) material.IMaterial {
	if mat, ok := r.materials[id]; ok {
		return mat
	}
	panic(fmt.Errorf(`unknown material "%s"`, id))
}

func (r *Repository) getTexture(file string) *texture.Texture2D {

	if text, ok := r.textures[file]; ok {
		return text
	}

	text, err := texture.NewTexture2DFromImage(file)
	if err != nil {
		panic(err)
	}

	text.SetMagFilter(gls.NEAREST) // avoid blurry upscale

	r.textures[file] = text

	return r.textures[file]
}

func (r *Repository) AppendFromYamlFile(filePath string) {
	rawYaml := file.ParseYamlFile[_YAMLMaterials](filePath)

	for id, def := range rawYaml {
		mat := material.NewStandard(getYAMLMaterialColor(def))
		mat.SetSpecularColor(getYAMLMaterialSpecularColor(def))
		mat.SetEmissiveColor(getYAMLMaterialEmissiveColor(def))

		if def.DepthMask != nil {
			mat.SetDepthMask(*def.DepthMask)
		}

		if def.Opacity != nil {
			mat.SetOpacity(*def.Opacity)
		}

		if def.Transparent != nil {
			mat.SetTransparent(*def.Transparent)
		}

		if def.Texture != nil {
			mat.AddTexture(r.getTexture(*def.Texture))
		}

		r.materials[id] = mat
	}
}

func NewFromYamlFile(filePath string) *Repository {
	repository := Repository{
		map[string]material.IMaterial{},
		map[string]*texture.Texture2D{},
	}

	repository.AppendFromYamlFile(filePath)

	return &repository
}

func GetRepository() *Repository {
	if RepositoryInstance == nil {
		RepositoryInstance = NewFromYamlFile(DefinitionFile)
	}

	return RepositoryInstance
}
