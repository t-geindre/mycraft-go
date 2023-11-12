package block

import (
	"fmt"
	"github.com/g3n/engine/geometry"
	"github.com/g3n/engine/graphic"
	"github.com/g3n/engine/math32"
	"mycraft/block/file"
	"mycraft/block/material"
)

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

func (r *Repository) Get(id string) *Block {
	if block, ok := r.Blocks[id]; ok {
		return block.Clone()
	}
	panic(fmt.Errorf(`unknown block "%s"`, id))
}

func (r *Repository) AppendFromYAMLFile(filePath string) {
	rawYaml := file.ParseYamlFile[_YAMLBlocks](filePath)

	for id, def := range rawYaml {
		var meshes []*graphic.Mesh
		switch def.Type {
		case "plant":
			meshes = r.createPlantMeshes(def)
		default:
			meshes = r.createBlockMeshes(def)
		}
		r.Blocks[id] = Block{
			Id:     id,
			Type:   def.Type,
			Meshes: meshes,
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

func (r *Repository) createBlockMeshes(def _YAMLBlock) []*graphic.Mesh {
	mesh := graphic.NewMesh(geometry.NewCube(1), nil)

	if len(def.Materials.All) > 0 {
		mesh.SetMaterial(r.Materials.Get(def.Materials.All))
		return []*graphic.Mesh{mesh}
	}

	if len(def.Materials.Sides) > 0 {
		mesh.AddGroupMaterial(r.Materials.Get(def.Materials.Sides), blockFaceBack)
		mesh.AddGroupMaterial(r.Materials.Get(def.Materials.Sides), blockFaceFront)
		mesh.AddGroupMaterial(r.Materials.Get(def.Materials.Sides), blockFaceLeft)
		mesh.AddGroupMaterial(r.Materials.Get(def.Materials.Sides), blockFaceRight)
		mesh.AddGroupMaterial(r.Materials.Get(def.Materials.Bottom), blockFaceBottom)
		mesh.AddGroupMaterial(r.Materials.Get(def.Materials.Top), blockFaceTop)
		return []*graphic.Mesh{mesh}
	}

	mesh.AddGroupMaterial(r.Materials.Get(def.Materials.Back), blockFaceBack)
	mesh.AddGroupMaterial(r.Materials.Get(def.Materials.Front), blockFaceFront)
	mesh.AddGroupMaterial(r.Materials.Get(def.Materials.Left), blockFaceLeft)
	mesh.AddGroupMaterial(r.Materials.Get(def.Materials.Right), blockFaceRight)
	mesh.AddGroupMaterial(r.Materials.Get(def.Materials.Bottom), blockFaceBottom)
	mesh.AddGroupMaterial(r.Materials.Get(def.Materials.Top), blockFaceTop)

	return []*graphic.Mesh{mesh}
}

func (r *Repository) createPlantMeshes(def _YAMLBlock) []*graphic.Mesh {
	plan := geometry.NewPlane(1, 1)
	mat := r.Materials.Get(def.Materials.All)
	var meshes []*graphic.Mesh

	for i := float32(0); i < 4; i++ {
		mesh := graphic.NewMesh(plan, mat)
		mesh.SetRotationY(math32.DegToRad(90) * i)
		meshes = append(meshes, mesh)
	}

	return meshes
}
