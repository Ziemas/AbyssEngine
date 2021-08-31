#version 330 core
#extension GL_ARB_explicit_uniform_location : require
in vec2 TexCoords;
out vec4 color;

layout (location = 0) uniform sampler2D image;

void main()
{
    color = texture(image, TexCoords);
}
