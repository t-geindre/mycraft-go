package main

import (
	"github.com/g3n/engine/material"
	"github.com/g3n/engine/math32"
	"github.com/g3n/engine/texture"
)

type YAMLMaterial struct {
	Color          *string
	Color3         *math32.Color
	SpecularColor  *string
	SpecularColor3 *math32.Color
	EmissiveColor  *string
	EmissiveColor3 *math32.Color
	Shininess      *float32
	Opacity        *float32
	Texture        *string
}

type YAMLMaterials map[string]YAMLMaterial

type Materials struct {
	materials map[string]*material.Standard
	textures  map[string]*texture.Texture2D
}

func (m *Materials) LoadFromYamlFile(file string) {
	rawYaml := ParseYamlFile[YAMLMaterials](file)
	println(rawYaml["grass_green"].Color)

}

func (m Materials) Get(name string) *material.Standard {
	return nil
}

func (m *Materials) GetTexture(name string) *texture.Texture2D {
	return nil
}
