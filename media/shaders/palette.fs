#version 330
#extension GL_ARB_explicit_uniform_location : require
// Input frag attributes
in vec2 TexCoords;

// Init uniform values
layout (location = 0) uniform sampler2D image;
layout (location = 1) uniform sampler2D paletteTex;
layout (location = 2) uniform float paletteOffset;

out vec4 color;

void main() {
  vec4 index = texture(image, TexCoords);
  color = texture(paletteTex, vec2(index.x, paletteOffset));
}
