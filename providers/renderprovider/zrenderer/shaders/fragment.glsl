#version 330 core
in vec2 TexCoords;
out vec4 color;

uniform sampler2D image;
uniform sampler2D paletteTex;
uniform int paletteOffset;
uniform bool usePalette;

void main() {
    if (usePalette) {
        float palYCoord = float(paletteOffset) / float(textureSize(paletteTex, 0).y - 1);
        vec4 index = texture(image, TexCoords);
        color = texture(paletteTex, vec2(index.x, palYCoord));
    } else {
        color = texture(image, TexCoords);
    }
}

