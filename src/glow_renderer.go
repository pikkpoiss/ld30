package main

import (
	twodee "../libs/twodee"
	"fmt"
	"github.com/go-gl/gl"
)

const GLOW_FRAGMENT = `#version 150
precision mediump float;

uniform sampler2D u_TextureUnit;
in vec2 v_TextureCoordinates;
uniform int Orientation;
uniform int BlurAmount;
uniform float BlurScale;
uniform float BlurStrength;
uniform vec2 BufferDimensions;
out vec4 v_FragData;
vec2 TexelSize = vec2(1.0 / BufferDimensions.x, 1.0 / BufferDimensions.y);

float Gaussian (float x, float deviation)
{
  return (1.0 / sqrt(2.0 * 3.141592 * deviation)) * exp(-((x * x) / (2.0 * deviation)));
}


void main()
{
  // Locals
  float halfBlur = float(BlurAmount) * 0.5;
  vec4 colour = vec4(0.0);
  vec4 texColour = vec4(0.0);

  // Gaussian deviation
  float deviation = halfBlur * 0.35;
  deviation *= deviation;
  float strength = 1.0 - BlurStrength;

  if ( Orientation == 0 ) {
    // Horizontal blur
    for (int i = 0; i < 10; ++i) {
      if ( i >= BlurAmount ) {
        break;
      }
      float offset = float(i) - halfBlur;
      texColour = texture(
        u_TextureUnit,
        v_TextureCoordinates + vec2(offset * TexelSize.x * BlurScale, 0.0)) * Gaussian(offset * strength, deviation);
      colour += texColour;
    }
  } else {
    // Vertical blur
    for (int i = 0; i < 10; ++i) {
      if ( i >= BlurAmount ) {
        break;
      }
      float offset = float(i) - halfBlur;
      texColour = texture(
        u_TextureUnit,
        v_TextureCoordinates + vec2(0.0, offset * TexelSize.y * BlurScale)) * Gaussian(offset * strength, deviation);
      colour += texColour;
    }
  }
  // Apply colour
  v_FragData = clamp(colour, 0.0, 1.0);
  v_FragData.w = 1.0;
}`

const GLOW_VERTEX = `#version 150

in vec4 a_Position;
in vec2 a_TextureCoordinates;

out vec2 v_TextureCoordinates;

void main()
{
    v_TextureCoordinates = a_TextureCoordinates;
    gl_Position = a_Position;
}`

type GlowRenderer struct {
	GlowFb              gl.Framebuffer
	GlowTex             gl.Texture
	BlurFb              gl.Framebuffer
	BlurTex             gl.Texture
	shader              gl.Program
	positionLoc         gl.AttribLocation
	textureLoc          gl.AttribLocation
	orientationLoc      gl.UniformLocation
	blurAmountLoc       gl.UniformLocation
	blurScaleLoc        gl.UniformLocation
	blurStrengthLoc     gl.UniformLocation
	bufferDimensionsLoc gl.UniformLocation
	textureUnitLoc      gl.UniformLocation
	coords              gl.Buffer
	width               int
	height              int
	oldwidth            int
	oldheight           int
}

func NewGlowRenderer(w, h int) (r *GlowRenderer, err error) {
	r = &GlowRenderer{
		width:  w,
		height: h,
	}
	_, _, r.oldwidth, r.oldheight = GetInteger4(gl.VIEWPORT)
	if r.shader, err = twodee.BuildProgram(GLOW_VERTEX, GLOW_FRAGMENT); err != nil {
		return
	}
	r.orientationLoc = r.shader.GetUniformLocation("Orientation")
	r.blurAmountLoc = r.shader.GetUniformLocation("BlurAmount")
	r.blurScaleLoc = r.shader.GetUniformLocation("BlurScale")
	r.blurStrengthLoc = r.shader.GetUniformLocation("BlurStrength")
	r.bufferDimensionsLoc = r.shader.GetUniformLocation("BufferDimensions")
	r.positionLoc = r.shader.GetAttribLocation("a_Position")
	r.textureLoc = r.shader.GetAttribLocation("a_TextureCoordinates")
	r.textureUnitLoc = r.shader.GetUniformLocation("u_TextureUnit")
	r.shader.BindFragDataLocation(0, "v_FragData")
	var size float32 = 1.0
	var rect = []float32{
		-size, -size, 0.0, 0, 0,
		-size, size, 0.0, 0, 1,
		size, -size, 0.0, 1, 0,
		size, size, 0.0, 1, 1,
	}
	if r.coords, err = twodee.CreateVBO(len(rect)*4, rect, gl.STATIC_DRAW); err != nil {
		return
	}

	if r.GlowFb, r.GlowTex, err = r.initFramebuffer(w, h); err != nil {
		return
	}
	if r.BlurFb, r.BlurTex, err = r.initFramebuffer(w, h); err != nil {
		return
	}
	return
}

func (r *GlowRenderer) initFramebuffer(w, h int) (fb gl.Framebuffer, tex gl.Texture, err error) {
	fb = gl.GenFramebuffer()
	fb.Bind()

	tex = gl.GenTexture()
	tex.Bind(gl.TEXTURE_2D)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.LINEAR)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.LINEAR)
	gl.TexImage2D(gl.TEXTURE_2D, 0, gl.RGBA, w, h, 0, gl.RGBA, gl.UNSIGNED_BYTE, nil)

	gl.FramebufferTexture2D(gl.DRAW_FRAMEBUFFER, gl.COLOR_ATTACHMENT0, gl.TEXTURE_2D, tex, 0)
	if err = r.GetError(); err != nil {
		return
	}
	gl.DrawBuffers(1, []gl.GLenum{gl.COLOR_ATTACHMENT0})

	rb := gl.GenRenderbuffer()
	rb.Bind()
	gl.RenderbufferStorage(gl.RENDERBUFFER, gl.STENCIL_INDEX8, w, h)
	rb.FramebufferRenderbuffer(gl.FRAMEBUFFER, gl.STENCIL_ATTACHMENT, gl.RENDERBUFFER)

	tex.Unbind(gl.TEXTURE_2D)
	fb.Unbind()
	rb.Unbind()
	return
}

func (r *GlowRenderer) GetError() error {
	if e := gl.GetError(); e != 0 {
		return fmt.Errorf("OpenGL error: %X", e)
	}
	var status = gl.CheckFramebufferStatus(gl.DRAW_FRAMEBUFFER)
	switch status {
	case gl.FRAMEBUFFER_COMPLETE:
		return nil
	case gl.FRAMEBUFFER_INCOMPLETE_ATTACHMENT:
		return fmt.Errorf("Attachment point unconnected")
	case gl.FRAMEBUFFER_INCOMPLETE_MISSING_ATTACHMENT:
		return fmt.Errorf("Missing attachment")
	case gl.FRAMEBUFFER_INCOMPLETE_DRAW_BUFFER:
		return fmt.Errorf("Draw buffer")
	case gl.FRAMEBUFFER_INCOMPLETE_READ_BUFFER:
		return fmt.Errorf("Read buffer")
	case gl.FRAMEBUFFER_UNSUPPORTED:
		return fmt.Errorf("Unsupported config")
	default:
		return fmt.Errorf("Unknown framebuffer error: %X", status)
	}
}

func (r *GlowRenderer) Delete() error {
	r.GlowFb.Delete()
	r.GlowTex.Delete()
	r.BlurFb.Delete()
	r.BlurTex.Delete()
	r.coords.Delete()
	return r.GetError()
}

func (r *GlowRenderer) Bind() error {
	r.GlowFb.Bind()
	gl.Enable(gl.STENCIL_TEST)
	gl.Viewport(0, 0, r.width, r.height)
	gl.ClearStencil(0)
	gl.ClearColor(0.0, 0.0, 0.0, 0.0)
	gl.StencilMask(0xFF) // Write to buffer
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.STENCIL_BUFFER_BIT)
	gl.StencilMask(0x00) // Don't write to buffer
	return nil
}

func (r *GlowRenderer) Draw() (err error) {
	r.shader.Use()
	r.textureUnitLoc.Uniform1i(0)
	r.coords.Bind(gl.ARRAY_BUFFER)
	r.positionLoc.AttribPointer(3, gl.FLOAT, false, 5*4, uintptr(0))
	r.textureLoc.AttribPointer(2, gl.FLOAT, false, 5*4, uintptr(3*4))
	r.blurAmountLoc.Uniform1i(6)
	r.blurScaleLoc.Uniform1f(1.0)
	r.blurStrengthLoc.Uniform1f(0.4)
	r.bufferDimensionsLoc.Uniform2f(float32(r.width), float32(r.height))

	r.BlurFb.Bind()
	gl.Viewport(0, 0, r.width, r.height)
	gl.ActiveTexture(gl.TEXTURE0)
	r.GlowTex.Bind(gl.TEXTURE_2D)
	r.orientationLoc.Uniform1i(0)
	gl.DrawArrays(gl.TRIANGLE_STRIP, 0, 4)
	r.BlurFb.Unbind()

	gl.Viewport(0, 0, r.oldwidth, r.oldheight)
	gl.BlendFunc(gl.ONE, gl.ONE)
	r.BlurTex.Bind(gl.TEXTURE_2D)
	r.orientationLoc.Uniform1i(1)
	gl.DrawArrays(gl.TRIANGLE_STRIP, 0, 4)

	r.coords.Unbind(gl.ARRAY_BUFFER)
	gl.BlendFunc(gl.SRC_ALPHA, gl.ONE_MINUS_SRC_ALPHA)
	return nil
}

func (r *GlowRenderer) Unbind() error {
	gl.Viewport(0, 0, r.oldwidth, r.oldheight)
	r.GlowFb.Unbind()
	gl.Disable(gl.STENCIL_TEST)
	return nil
}

func (r *GlowRenderer) DisableOutput() {
	gl.ColorMask(false, false, false, false)
	gl.StencilFunc(gl.NEVER, 1, 0xFF)                // Never pass
	gl.StencilOp(gl.REPLACE, gl.REPLACE, gl.REPLACE) // Replace to ref=1
	gl.StencilMask(0xFF)                             // Write to buffer
}

func (r *GlowRenderer) EnableOutput() {
	gl.ColorMask(true, true, true, true)
	gl.StencilMask(0x00)              // No more writing
	gl.StencilFunc(gl.EQUAL, 0, 0xFF) // Only pass where stencil is 0
}

// Convenience function for glGetIntegerv
func GetInteger4(pname gl.GLenum) (v0, v1, v2, v3 int) {
	var values = []int32{0, 0, 0, 0}
	gl.GetIntegerv(pname, values)
	v0 = int(values[0])
	v1 = int(values[1])
	v2 = int(values[2])
	v3 = int(values[3])
	return
}
