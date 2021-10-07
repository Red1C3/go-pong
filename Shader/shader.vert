#version 330
layout(location=0) in vec2 pos_model;
out vec2 out_model;
uniform mat4 VP;
uniform mat4 M;
void main(){
  out_model=pos_model;
  gl_Position=VP*M*vec4(pos_model,0,1);
}