#version 330
out vec4 color;
in vec2 out_model;
uniform int isBall;
void main(){
  color=vec4(1,1,1,1);
  if(isBall==1){
    float vecLength=length(out_model);
    if(vecLength>0.5){
      discard;
    }
  }
}