#version 330
layout(location=0) in vec2 pos_model;
uniform mat4 VP;
void main(){
  gl_Position=VP*vec4(pos_model,0,1);
}