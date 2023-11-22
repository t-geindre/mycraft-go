package mesh

import (
	"github.com/g3n/engine/core"
	"github.com/g3n/engine/geometry"
	"github.com/g3n/engine/gls"
	"github.com/g3n/engine/graphic"
	"github.com/g3n/engine/material"
	"github.com/g3n/engine/math32"
)

type Skybox struct {
	graphic.Graphic             // embedded graphic object
	uniMVm          gls.Uniform // model view matrix uniform location cache
	uniMVPm         gls.Uniform // model view projection matrix uniform cache
	uniNm           gls.Uniform // normal matrix uniform cache
}

func NewSkybox() *Skybox {
	skybox := new(Skybox)

	geom := geometry.NewCube(1)
	skybox.Graphic.Init(skybox, geom, gls.TRIANGLES)
	skybox.Graphic.SetCullable(false)
	mat := material.NewStandard(&math32.Color{R: 0.2, G: 0.4, B: 0.8})
	mat.SetSide(material.SideBack)
	mat.SetUseLights(material.UseLightNone)
	mat.SetDepthMask(false)

	for i := 0; i < 6; i++ {
		skybox.AddGroupMaterial(skybox, mat, i)
	}

	skybox.uniMVm.Init("ModelViewMatrix")
	skybox.uniMVPm.Init("MVP")
	skybox.uniNm.Init("NormalMatrix")

	skybox.SetRenderOrder(100)

	return skybox
}

func (skybox *Skybox) RenderSetup(gs *gls.GLS, rinfo *core.RenderInfo) {

	mvm := *skybox.ModelViewMatrix()
	mvm[12] = 0
	mvm[13] = 0
	mvm[14] = 0

	location := skybox.uniMVm.Location(gs)
	gs.UniformMatrix4fv(location, 1, false, &mvm[0])

	var mvpm math32.Matrix4
	mvpm.MultiplyMatrices(&rinfo.ProjMatrix, &mvm)
	location = skybox.uniMVPm.Location(gs)
	gs.UniformMatrix4fv(location, 1, false, &mvpm[0])

	var nm math32.Matrix3
	nm.GetNormalMatrix(&mvm)
	location = skybox.uniNm.Location(gs)
	gs.UniformMatrix3fv(location, 1, false, &nm[0])
}
