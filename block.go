package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/g3n/engine/geometry"
	"github.com/g3n/engine/graphic"
	"github.com/g3n/engine/material"
	"github.com/g3n/engine/math32"
	"github.com/g3n/engine/texture"
)

const BLOCK_LIBRARY = "assets/blocks.json"
const BLOCk_FACE_RIGHT = 0
const BLOCk_FACE_LEFT = 6
const BLOCk_FACE_TOP = 12
const BLOCk_FACE_BOTTOM = 18
const BLOCk_FACE_FRONT = 24
const BLOCk_FACE_BACK = 30

type JSONBlocks map[string]JSONBlock

type JSONBlock struct {
	Sides  string
	Top    string
	Bottom string
	All    string
	Left   string
	Right  string
	Back   string
	Front  string
}

var materials = map[string]*material.Standard{}

func LoadBlocks() map[string]*graphic.Mesh {
	jsonBlocks := loadJsonBlocks()
	blocks := make(map[string]*graphic.Mesh, len(jsonBlocks))

	for blockName, blockDefinition := range jsonBlocks {
		blocks[blockName] = graphic.NewMesh(geometry.NewCube(1), nil)

		if len(blockDefinition.All) > 0 {
			blocks[blockName].AddMaterial(getMaterial(blockDefinition.All), 0, 0)
			continue
		}

		if len(blockDefinition.Sides) > 0 {
			blocks[blockName].AddMaterial(getMaterial(blockDefinition.Sides), BLOCk_FACE_BACK, 6)
			blocks[blockName].AddMaterial(getMaterial(blockDefinition.Sides), BLOCk_FACE_FRONT, 6)
			blocks[blockName].AddMaterial(getMaterial(blockDefinition.Sides), BLOCk_FACE_LEFT, 6)
			blocks[blockName].AddMaterial(getMaterial(blockDefinition.Sides), BLOCk_FACE_RIGHT, 6)
			blocks[blockName].AddMaterial(getMaterial(blockDefinition.Bottom), BLOCk_FACE_BOTTOM, 6)
			blocks[blockName].AddMaterial(getMaterial(blockDefinition.Top), BLOCk_FACE_TOP, 6)
			continue
		}
	}

	return blocks
}

func loadJsonBlocks() JSONBlocks {
	rawJson, err := os.ReadFile(BLOCK_LIBRARY)
	if err != nil {
		panic(err)
	}

	jsonBlocks := JSONBlocks{}
	err = json.Unmarshal(rawJson, &jsonBlocks)
	if err != nil {
		panic(err)
	}

	return jsonBlocks
}

func getMaterial(textureFile string) *material.Standard {
	if material, ok := materials[textureFile]; ok {
		return material
	}

	fmt.Printf("Loading texture \"%s\"....", textureFile)
	texture, err := texture.NewTexture2DFromImage(textureFile)
	if err != nil {
		panic(err)
	}

	fmt.Println("DONE")

	material := material.NewStandard(math32.NewColor("White"))
	material.AddTexture(texture)
	materials[textureFile] = material

	return material
}
