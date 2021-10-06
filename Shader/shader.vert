#version 330
layout(location=0) in vec2 pos_model;
void main(){
  gl_Position=vec4(pos_model,0,1);
}