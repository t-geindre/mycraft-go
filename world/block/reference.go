package block

import (
	matRef "mycraft/mesh/material"
)

func getReference() map[uint16]*Block {
	ref := make(map[uint16]*Block)
	matRepo := matRef.GetRepository()

	MaterialsSides := []uint16{MaterialNorth, MaterialSouth, MaterialEast, MaterialWest}
	MaterialsCube := append([]uint16{MaterialTop, MaterialBottom}, MaterialsSides...)

	// GRASS
	ref[TypeGrass] = NewBlock(TypeGrass, KindCube)
	ref[TypeGrass].SetMaterial(matRepo.Get(matRef.BlockGrassTop), MaterialTop)
	ref[TypeGrass].SetMaterial(matRepo.Get(matRef.BlockDirt), MaterialBottom)
	ref[TypeGrass].SetMaterial(matRepo.Get(matRef.BlockGrassSide), MaterialsSides...)

	// DIRT
	ref[TypeDirt] = NewBlock(TypeDirt, KindCube)
	ref[TypeDirt].SetMaterial(matRepo.Get(matRef.BlockDirt), MaterialsCube...)

	// STONE
	ref[TypeStone] = NewBlock(TypeStone, KindCube)
	ref[TypeStone].SetMaterial(matRepo.Get(matRef.BlockStone), MaterialsCube...)

	// WATER
	ref[TypeWater] = NewBlock(TypeWater, KindCube)
	ref[TypeWater].SetMaterial(matRepo.Get(matRef.BlockWater), MaterialsCube...)
	ref[TypeWater].SetTransparent(true)

	// SAND
	ref[TypeSand] = NewBlock(TypeSand, KindCube)
	ref[TypeSand].SetMaterial(matRepo.Get(matRef.BlockSand), MaterialsCube...)

	// BEDROCK
	ref[TypeBedrock] = NewBlock(TypeBedrock, KindCube)
	ref[TypeBedrock].SetMaterial(matRepo.Get(matRef.Bedrock), MaterialsCube...)

	// DANDELION
	ref[TypeDandelion] = NewBlock(TypeDandelion, KindPlant)
	ref[TypeDandelion].SetMaterial(matRepo.Get(matRef.Dandelion), MaterialTop)

	// GRAVEL
	ref[TypeGravel] = NewBlock(TypeGravel, KindCube)
	ref[TypeGravel].SetMaterial(matRepo.Get(matRef.Gravel), MaterialsCube...)

	// SPRUCE LOG
	ref[TypeOakLog] = NewBlock(TypeOakLog, KindCube)
	ref[TypeOakLog].SetMaterial(matRepo.Get(matRef.OakSpruceLogTop), MaterialTop)
	ref[TypeOakLog].SetMaterial(matRepo.Get(matRef.OakSpruceLogTop), MaterialBottom)
	ref[TypeOakLog].SetMaterial(matRepo.Get(matRef.OakLogSide), MaterialsSides...)

	// OAK LEAVES
	ref[TypeOakLeaves] = NewBlock(TypeOakLeaves, KindCube)
	ref[TypeOakLeaves].SetMaterial(matRepo.Get(matRef.OakLeaves), MaterialsCube...)
	ref[TypeOakLeaves].SetTransparent(true)

	return ref
}
