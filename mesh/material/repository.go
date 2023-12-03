package material

import (
	"fmt"
	"github.com/g3n/engine/gls"
	"github.com/g3n/engine/material"
	"github.com/g3n/engine/math32"
	"github.com/g3n/engine/texture"
)

type Repository struct {
	materials map[uint16]material.IMaterial
	textures  map[string]*texture.Texture2D
}

var RepositoryInstance *Repository

func (r *Repository) Get(id uint16) material.IMaterial {
	if mat, ok := r.materials[id]; ok {
		return mat
	}
	panic(fmt.Errorf(`unknown material "%d"`, id))
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
	text.SetWrapT(gls.REPEAT)      // Repeat texture on bigger surfaces
	text.SetWrapS(gls.REPEAT)      // Repeat texture on bigger surfaces

	r.textures[file] = text

	return r.textures[file]
}

func newRepository() *Repository {
	r := new(Repository)
	r.materials = make(map[uint16]material.IMaterial)
	r.textures = make(map[string]*texture.Texture2D)

	for id, mat := range materialReference() {
		m := material.NewStandard(&math32.Color{R: 1, G: 1, B: 1})
		if mat.TextureFile != "" {
			m.AddTexture(r.getTexture(mat.TextureFile))
		}
		if mat.Setup != nil {
			mat.Setup(m)
		}
		r.materials[id] = m
	}

	return r
}

func GetRepository() *Repository {
	if RepositoryInstance == nil {
		RepositoryInstance = newRepository()
	}

	return RepositoryInstance
}
