package block

import (
	"fmt"
	"github.com/g3n/engine/core"
	"github.com/g3n/engine/geometry"
	"github.com/g3n/engine/graphic"
	"github.com/g3n/engine/math32"
	"mycraft/block/file"
	"mycraft/block/material"
)

type Repository struct {
	Blocks    map[string]func() *Block
	Materials *material.Repository
}

const DefinitionFile = "assets/blocks.yaml"

var RepositoryInstance *Repository

const (
	blockFaceRight = iota
	blockFaceLeft
	blockFaceTop
	blockFaceBottom
	blockFaceFront
	blockFaceBack
)

func (r *Repository) Get(id string) *Block {
	if block, ok := r.Blocks[id]; ok {
		return block()
	}
	panic(fmt.Errorf(`unknown block "%s"`, id))
}

func (r *Repository) AppendFromYAMLFile(filePath string) {
	rawYaml := file.ParseYamlFile[_YAMLBlocks](filePath)

	for id, def := range rawYaml {
		r.Blocks[id] = func(def _YAMLBlock) func() *Block {
			return func() *Block {
				node := core.NewNode()
				switch def.Type {
				case "plant":
					r.addPlantMeshes(def, node)
				default:
					r.addCubeMeshes(def, node)
				}
				return &Block{
					Id:   id,
					Type: def.Type,
					Node: node,
				}
			}
		}(def)

	}
}

func NewFromYamlFile(filePath string, materials *material.Repository) *Repository {
	repository := Repository{
		Blocks:    make(map[string]func() *Block),
		Materials: materials,
	}
	repository.AppendFromYAMLFile(filePath)

	return &repository
}

func (r *Repository) addCubeMeshes(def _YAMLBlock, node *core.Node) {
	mesh := graphic.NewMesh(geometry.NewCube(1), nil)
	node.Add(mesh)

	if len(def.Materials.All) > 0 {
		mesh.SetMaterial(r.Materials.Get(def.Materials.All))
		return
	}

	if len(def.Materials.Sides) > 0 {
		mesh.AddGroupMaterial(r.Materials.Get(def.Materials.Sides), blockFaceBack)
		mesh.AddGroupMaterial(r.Materials.Get(def.Materials.Sides), blockFaceFront)
		mesh.AddGroupMaterial(r.Materials.Get(def.Materials.Sides), blockFaceLeft)
		mesh.AddGroupMaterial(r.Materials.Get(def.Materials.Sides), blockFaceRight)
		mesh.AddGroupMaterial(r.Materials.Get(def.Materials.Bottom), blockFaceBottom)
		mesh.AddGroupMaterial(r.Materials.Get(def.Materials.Top), blockFaceTop)
		return
	}

	mesh.AddGroupMaterial(r.Materials.Get(def.Materials.Back), blockFaceBack)
	mesh.AddGroupMaterial(r.Materials.Get(def.Materials.Front), blockFaceFront)
	mesh.AddGroupMaterial(r.Materials.Get(def.Materials.Left), blockFaceLeft)
	mesh.AddGroupMaterial(r.Materials.Get(def.Materials.Right), blockFaceRight)
	mesh.AddGroupMaterial(r.Materials.Get(def.Materials.Bottom), blockFaceBottom)
	mesh.AddGroupMaterial(r.Materials.Get(def.Materials.Top), blockFaceTop)

	return
}

func (r *Repository) addPlantMeshes(def _YAMLBlock, node *core.Node) {
	plan := geometry.NewPlane(1, 1)
	mat := r.Materials.Get(def.Materials.All)

	for i := float32(0); i < 4; i++ {
		mesh := graphic.NewMesh(plan, mat)
		mesh.SetRotationY(math32.DegToRad(90) * i)
		node.Add(mesh)
	}
}

func GetRepository() *Repository {
	if RepositoryInstance == nil {
		RepositoryInstance = NewFromYamlFile(DefinitionFile, material.GetRepository())
	}

	return RepositoryInstance
}
