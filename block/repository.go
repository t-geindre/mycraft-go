package block

import (
	"fmt"
	"github.com/g3n/engine/geometry"
	"github.com/g3n/engine/graphic"
	"mycraft/file"
	"mycraft/material"
)

type Block struct {
	Id         string
	CreateMesh func() *graphic.Mesh
}

type Repository struct {
	Blocks    map[string]Block
	Materials material.Repository
}

const (
	blockFaceRight = iota * 6
	blockFaceLeft
	blockFaceTop
	blockFaceBottom
	blockFaceFront
	blockFaceBack
)

func (r *Repository) Get(id string) Block {
	if mat, ok := r.Blocks[id]; ok {
		return mat
	}
	panic(fmt.Errorf(`unknown block "%s"`, id))
}

func (r *Repository) AppendFromYAMLFile(filePath string) {
	rawYaml := file.ParseYamlFile[_YAMLBlocks](filePath)

	for id, def := range rawYaml {
		r.Blocks[id] = Block{
			Id: id,
			CreateMesh: func(def _YAMLBlock) func() *graphic.Mesh {
				return func() *graphic.Mesh {
					mesh := graphic.NewMesh(geometry.NewCube(1), nil)

					if len(def.Materials.All) > 0 {
						mesh.AddMaterial(r.Materials.Get(def.Materials.All), 0, 0)
						return mesh
					}

					if len(def.Materials.Sides) > 0 {
						mesh.AddMaterial(r.Materials.Get(def.Materials.Sides), blockFaceBack, 6)
						mesh.AddMaterial(r.Materials.Get(def.Materials.Sides), blockFaceFront, 6)
						mesh.AddMaterial(r.Materials.Get(def.Materials.Sides), blockFaceLeft, 6)
						mesh.AddMaterial(r.Materials.Get(def.Materials.Sides), blockFaceRight, 6)
						mesh.AddMaterial(r.Materials.Get(def.Materials.Bottom), blockFaceBottom, 6)
						mesh.AddMaterial(r.Materials.Get(def.Materials.Top), blockFaceTop, 6)
						return mesh
					}

					mesh.AddMaterial(r.Materials.Get(def.Materials.Back), blockFaceBack, 6)
					mesh.AddMaterial(r.Materials.Get(def.Materials.Front), blockFaceFront, 6)
					mesh.AddMaterial(r.Materials.Get(def.Materials.Left), blockFaceLeft, 6)
					mesh.AddMaterial(r.Materials.Get(def.Materials.Right), blockFaceRight, 6)
					mesh.AddMaterial(r.Materials.Get(def.Materials.Bottom), blockFaceBottom, 6)
					mesh.AddMaterial(r.Materials.Get(def.Materials.Top), blockFaceTop, 6)

					return mesh
				}
			}(def),
		}
	}
}

func CreateFromYAMLFile(filePath string, materials material.Repository) Repository {
	repository := Repository{
		Blocks:    map[string]Block{},
		Materials: materials,
	}
	repository.AppendFromYAMLFile(filePath)
	return repository
}
