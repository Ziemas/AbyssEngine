#version 330 core
in vec2 TexCoords;
out vec4 color;

uniform sampler2D image;
uniform sampler2D paletteTex;
uniform int paletteOffset;
uniform bool usePalette;

void main() {
    if (usePalette == false) {
        color = texture(image, TexCoords);
    } else {
        vec4 index = texture(image, TexCoords);
        color = texture(paletteTex, vec2(index.x, paletteOffset));
    }
}
