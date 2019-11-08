#version 410 core

in vec2 vert;

out vec2 fragTexCoord;

void main() {
    fragTexCoord = vert;
    gl_Position = vec4(vert.x * 2 - 1, (1 - vert.y) * 2 - 1, 0, 1);
}
