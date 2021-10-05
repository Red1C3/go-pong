#ifndef _H_RENDERER
#define _H_RENDERER
#include<GL/glew.h>
#include<GLFW/glfw3.h>
#include<Events.h>
typedef struct{

} DrawInfo;
int initRenderer();
Event loop(DrawInfo drawInfo);
Event render(DrawInfo drawInfo);
int terminateRenderer();
#endif