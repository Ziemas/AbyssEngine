#version 330
#extension GL_ARB_explicit_uniform_location : require
// Input frag attributes
in vec2 TexCoords;

// Init uniform values
layout (location = 3) uniform sampler2D image;
layout (location = 4) uniform sampler2D paletteTex;
layout (location = 5) uniform float paletteOffset;

out vec4 color;

void main() {
  vec4 index = texture(image, TexCoords);
  if (index.x < (1.0 / 255)) {
      discard;
  }

  color = texture(paletteTex, vec2(index.x, paletteOffset));
}
