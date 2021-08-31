#version 330 core
#extension GL_ARB_explicit_uniform_location : require
layout (location = 0) in vec4 vertex; // <vec2 position, vec2 texCoords>

out vec2 TexCoords;

layout (location = 1) uniform mat4 model;
layout (location = 2) uniform mat4 projection;

void main()
{
    TexCoords = vertex.zw;
    gl_Position = projection * model * vec4(vertex.xy, 0.0, 1.0);
}
