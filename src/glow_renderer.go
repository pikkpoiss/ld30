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
out vec4 v_FragData;

void main()
{
    vec2 texcoords = v_TextureCoordinates;
    vec4 texcolor = texture(u_TextureUnit, texcoords);
    //v_FragData = texture(u_TextureUnit, texcoords);
    //v_FragData = vec4(1.0, 1.0, 1.0, 1.0) - texture(u_TextureUnit, texcoords);
    v_FragData = vec4(0.0, 1.0, 0.0, texcolor.a);
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
	Framebuffer    gl.Framebuffer
	Glowbuffer     gl.Texture
	Renderbuffer   gl.Renderbuffer
	shader         gl.Program
	positionLoc    gl.AttribLocation
	textureLoc     gl.AttribLocation
	textureUnitLoc gl.UniformLocation
	coords         gl.Buffer
	width          int
	height         int
	oldwidth       int
	oldheight      int
}

func NewGlowRenderer(w, h int) (r *GlowRenderer, err error) {
	r = &GlowRenderer{
		width:  w,
		height: h,
	}

	if r.shader, err = twodee.BuildProgram(GLOW_VERTEX, GLOW_FRAGMENT); err != nil {
		return
	}
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

	r.Framebuffer = gl.GenFramebuffer()
	r.Framebuffer.Bind()

	r.Glowbuffer = gl.GenTexture()
	r.Glowbuffer.Bind(gl.TEXTURE_2D)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.NEAREST)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.NEAREST)
	gl.TexImage2D(gl.TEXTURE_2D, 0, gl.RGBA, w, h, 0, gl.RGBA, gl.UNSIGNED_BYTE, nil)
	gl.Disable(gl.CULL_FACE)

	gl.FramebufferTexture2D(gl.DRAW_FRAMEBUFFER, gl.COLOR_ATTACHMENT0, gl.TEXTURE_2D, r.Glowbuffer, 0)
	if err = r.GetError(); err != nil {
		return
	}
	gl.DrawBuffers(1, []gl.GLenum{gl.COLOR_ATTACHMENT0})

	r.Renderbuffer = gl.GenRenderbuffer()
	r.Renderbuffer.Bind()
	gl.RenderbufferStorage(gl.RENDERBUFFER, gl.DEPTH_COMPONENT24, w, h)
	r.Renderbuffer.FramebufferRenderbuffer(gl.FRAMEBUFFER, gl.DEPTH_ATTACHMENT, gl.RENDERBUFFER)

	if err = r.GetError(); err != nil {
		return
	}
	r.Glowbuffer.Unbind(gl.TEXTURE_2D)
	r.Framebuffer.Unbind()
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
	r.Framebuffer.Delete()
	return r.GetError()
}

func (r *GlowRenderer) Bind() error {
	r.Framebuffer.Bind()
	_, _, r.oldwidth, r.oldheight = GetInteger4(gl.VIEWPORT)
	gl.Viewport(0, 0, r.width, r.height)
	gl.ClearColor(0.0, 0.0, 0.0, 0.0)
	gl.ClearDepth(1.0)
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
	return r.GetError()
}

func (r *GlowRenderer) Draw() (err error) {
	r.shader.Use()
	if err = r.GetError(); err != nil {
		return
	}
	gl.ActiveTexture(gl.TEXTURE0)
	if err = r.GetError(); err != nil {
		return
	}
	r.Glowbuffer.Bind(gl.TEXTURE_2D)
	if err = r.GetError(); err != nil {
		return
	}
	r.textureUnitLoc.Uniform1i(0)
	if err = r.GetError(); err != nil {
		return
	}
	r.coords.Bind(gl.ARRAY_BUFFER)
	if err = r.GetError(); err != nil {
		return
	}
	r.positionLoc.AttribPointer(3, gl.FLOAT, false, 5*4, uintptr(0))
	if err = r.GetError(); err != nil {
		return
	}
	r.textureLoc.AttribPointer(2, gl.FLOAT, false, 5*4, uintptr(3*4))
	if err = r.GetError(); err != nil {
		return
	}
	r.positionLoc.EnableArray()
	if err = r.GetError(); err != nil {
		return
	}
	r.textureLoc.EnableArray()
	if err = r.GetError(); err != nil {
		return
	}
	gl.DrawArrays(gl.TRIANGLE_STRIP, 0, 4)
	if err = r.GetError(); err != nil {
		return
	}
	r.coords.Unbind(gl.ARRAY_BUFFER)
	if err = r.GetError(); err != nil {
		return
	}
	return nil
}

func (r *GlowRenderer) Unbind() error {
	r.Framebuffer.Unbind()
	gl.Viewport(0, 0, r.oldwidth, r.oldheight)
	return r.GetError()
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
