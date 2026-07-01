package main

import (
	"testing"

	mgl "github.com/go-gl/mathgl/mgl32"
)

type fakeTexture struct {
	data    []byte
	invalid bool
}

func (t *fakeTexture) SetData(data []byte)                                         { t.data = append([]byte(nil), data...) }
func (t *fakeTexture) SetSubData(data []byte, x, y, width, height, stride int32)   {}
func (t *fakeTexture) SetDataG(data []byte, mag, min, ws, wt TextureSamplingParam) {}
func (t *fakeTexture) SetPixelData(data []float32)                                 {}
func (t *fakeTexture) IsValid() bool                                               { return !t.invalid }
func (t *fakeTexture) GetWidth() int32                                             { return 0 }
func (t *fakeTexture) GetHeight() int32                                            { return 0 }
func (t *fakeTexture) CopyData(src *Texture)                                       {}

type fakeRenderer struct {
	tex *fakeTexture
}

func (r *fakeRenderer) GetName() string                                                    { return "fake" }
func (r *fakeRenderer) Init()                                                              {}
func (r *fakeRenderer) Close()                                                             {}
func (r *fakeRenderer) BeginFrame(clearColor bool)                                         {}
func (r *fakeRenderer) EndFrame()                                                          {}
func (r *fakeRenderer) Await()                                                             {}
func (r *fakeRenderer) IsModelEnabled() bool                                               { return false }
func (r *fakeRenderer) IsShadowEnabled() bool                                              { return false }
func (r *fakeRenderer) LoadCustomSpriteShader(shaderName string, shaderData []byte) uint32 { return 0 }
func (r *fakeRenderer) UnloadCustomSpriteShader(shaderName string)                         {}
func (r *fakeRenderer) SetSpritePipeline(shaderName string)                                {}
func (r *fakeRenderer) SetCustomUniforms(params [16]float32)                               {}
func (r *fakeRenderer) NeedsGrabPass() bool                                                { return false }
func (r *fakeRenderer) ResolveBackBuffer() Texture                                         { return nil }
func (r *fakeRenderer) EnableBlending(eq BlendEquation, src, dst BlendFunc)                {}
func (r *fakeRenderer) DisableBlending()                                                   {}
func (r *fakeRenderer) prepareShadowMapPipeline(bufferIndex uint32)                        {}
func (r *fakeRenderer) setShadowMapPipeline(doubleSided, invertFrontFace, useUV, useNormal, useTangent, useVertColor, useJoint0, useJoint1 bool, numVertices, vertAttrOffset uint32) {
}
func (r *fakeRenderer) ReleaseShadowPipeline()                                    {}
func (r *fakeRenderer) prepareModelPipeline(bufferIndex uint32, env *Environment) {}
func (r *fakeRenderer) SetModelPipeline(eq BlendEquation, src, dst BlendFunc, depthTest, depthMask, doubleSided, invertFrontFace, useUV, useNormal, useTangent, useVertColor, useJoint0, useJoint1, useOutlineAttribute bool, numVertices, vertAttrOffset uint32) {
}
func (r *fakeRenderer) SetMeshOutlinePipeline(invertFrontFace bool, meshOutline float32) {}
func (r *fakeRenderer) ReleaseModelPipeline()                                            {}
func (r *fakeRenderer) newTexture(width, height, depth int32, filter bool) Texture       { return nil }
func (r *fakeRenderer) newPaletteTexture() Texture                                       { return r.tex }
func (r *fakeRenderer) newModelTexture(width, height, depth int32, filter bool) Texture  { return nil }
func (r *fakeRenderer) newDataTexture(width, height int32) Texture                       { return nil }
func (r *fakeRenderer) newHDRTexture(width, height int32) Texture                        { return nil }
func (r *fakeRenderer) newCubeMapTexture(widthHeight int32, mipmap bool, lowestMipLevel int32) Texture {
	return nil
}
func (r *fakeRenderer) ReadPixels(data []uint8, width, height int)                    {}
func (r *fakeRenderer) EnableScissor(x, y, width, height int32)                       {}
func (r *fakeRenderer) DisableScissor()                                               {}
func (r *fakeRenderer) SetUniformI(name string, val int)                              {}
func (r *fakeRenderer) SetUniformF(name string, values ...float32)                    {}
func (r *fakeRenderer) SetUniformFv(name string, values []float32)                    {}
func (r *fakeRenderer) SetUniformMatrix(name string, value []float32)                 {}
func (r *fakeRenderer) SetTexture(name string, tex Texture)                           {}
func (r *fakeRenderer) SetModelUniformI(name string, val int)                         {}
func (r *fakeRenderer) SetModelUniformF(name string, values ...float32)               {}
func (r *fakeRenderer) SetModelUniformFv(name string, values []float32)               {}
func (r *fakeRenderer) SetModelUniformMatrix(name string, value []float32)            {}
func (r *fakeRenderer) SetModelUniformMatrix3(name string, value []float32)           {}
func (r *fakeRenderer) SetModelTexture(name string, t Texture)                        {}
func (r *fakeRenderer) SetShadowMapUniformI(name string, val int)                     {}
func (r *fakeRenderer) SetShadowMapUniformF(name string, values ...float32)           {}
func (r *fakeRenderer) SetShadowMapUniformFv(name string, values []float32)           {}
func (r *fakeRenderer) SetShadowMapUniformMatrix(name string, value []float32)        {}
func (r *fakeRenderer) SetShadowMapUniformMatrix3(name string, value []float32)       {}
func (r *fakeRenderer) SetShadowMapTexture(name string, t Texture)                    {}
func (r *fakeRenderer) SetShadowFrameTexture(i uint32)                                {}
func (r *fakeRenderer) SetShadowFrameCubeTexture(i uint32)                            {}
func (r *fakeRenderer) SetVertexData(values ...float32)                               {}
func (r *fakeRenderer) SetModelVertexData(bufferIndex uint32, values []byte)          {}
func (r *fakeRenderer) SetModelIndexData(bufferIndex uint32, values ...uint32)        {}
func (r *fakeRenderer) RenderQuad()                                                   {}
func (r *fakeRenderer) RenderElements(mode PrimitiveMode, count, offset int)          {}
func (r *fakeRenderer) RenderShadowMapElements(mode PrimitiveMode, count, offset int) {}
func (r *fakeRenderer) RenderCubeMap(envTexture Texture, cubeTexture Texture)         {}
func (r *fakeRenderer) RenderFilteredCubeMap(distribution int32, cubeTexture Texture, filteredTexture Texture, mipmapLevel, sampleCount int32, roughness float32) {
}
func (r *fakeRenderer) RenderLUT(distribution int32, cubeTexture Texture, lutTexture Texture, sampleCount int32) {
}
func (r *fakeRenderer) PerspectiveProjectionMatrix(angle, aspect, near, far float32) mgl.Mat4 {
	return mgl.Ident4()
}
func (r *fakeRenderer) OrthographicProjectionMatrix(left, right, bottom, top, near, far float32) mgl.Mat4 {
	return mgl.Ident4()
}
func (r *fakeRenderer) SetVSync(interval int) {}
func (r *fakeRenderer) NewWorkerThread() bool { return false }

func TestNewTextureFromPalette(t *testing.T) {
	t.Parallel()

	prev := gfx
	ftex := &fakeTexture{}
	gfx = &fakeRenderer{tex: ftex}
	t.Cleanup(func() { gfx = prev })

	tx := NewTextureFromPalette(nil)
	if tx != ftex {
		t.Fatal("expected fake palette texture to be returned")
	}
	if len(ftex.data) != 0 {
		t.Fatalf("expected nil palette to clear texture data, got %d bytes", len(ftex.data))
	}

	pal := []uint32{0x11223344}
	_ = NewTextureFromPalette(pal)
	if len(ftex.data) != 1024 {
		t.Fatalf("expected packed palette data to be 1024 bytes, got %d", len(ftex.data))
	}
}
