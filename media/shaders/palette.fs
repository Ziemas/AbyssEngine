#version 330 core
#extension GL_ARB_explicit_uniform_location : require
#extension GL_ARB_shading_language_420pack: enable
in vec2 TexCoords;

layout (binding = 0) uniform sampler2D image;
layout (binding = 1) uniform sampler2D paletteTex;
layout (location = 5) uniform float paletteOffset;

out vec4 color;

void main() {
  vec4 index = texture(image, TexCoords);
  color = texture(paletteTex, vec2(index.x, paletteOffset));
}
