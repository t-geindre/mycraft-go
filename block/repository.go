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
	blockFaceRight = iota
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
						mesh.SetMaterial(r.Materials.Get(def.Materials.All))
						return mesh
					}

					if len(def.Materials.Sides) > 0 {
						mesh.AddGroupMaterial(r.Materials.Get(def.Materials.Sides), blockFaceBack)
						mesh.AddGroupMaterial(r.Materials.Get(def.Materials.Sides), blockFaceFront)
						mesh.AddGroupMaterial(r.Materials.Get(def.Materials.Sides), blockFaceLeft)
						mesh.AddGroupMaterial(r.Materials.Get(def.Materials.Sides), blockFaceRight)
						mesh.AddGroupMaterial(r.Materials.Get(def.Materials.Bottom), blockFaceBottom)
						mesh.AddGroupMaterial(r.Materials.Get(def.Materials.Top), blockFaceTop)
						return mesh
					}

					mesh.AddGroupMaterial(r.Materials.Get(def.Materials.Back), blockFaceBack)
					mesh.AddGroupMaterial(r.Materials.Get(def.Materials.Front), blockFaceFront)
					mesh.AddGroupMaterial(r.Materials.Get(def.Materials.Left), blockFaceLeft)
					mesh.AddGroupMaterial(r.Materials.Get(def.Materials.Right), blockFaceRight)
					mesh.AddGroupMaterial(r.Materials.Get(def.Materials.Bottom), blockFaceBottom)
					mesh.AddGroupMaterial(r.Materials.Get(def.Materials.Top), blockFaceTop)

					return mesh
				}
			}(def),
		}
	}
}

func NewFromYAMLFile(filePath string, materials material.Repository) Repository {
	repository := Repository{
		Blocks:    map[string]Block{},
		Materials: materials,
	}
	repository.AppendFromYAMLFile(filePath)
	return repository
}
