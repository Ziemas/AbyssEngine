#version 330 core
#extension GL_ARB_explicit_uniform_location : require
#extension GL_ARB_shading_language_420pack: enable
in vec2 TexCoords;
out vec4 color;

layout (binding = 0) uniform sampler2D image;

void main()
{
    color = texture(image, TexCoords);
}
