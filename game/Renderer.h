#ifndef _H_RENDERER
#define _H_RENDERER
#include<GLFW/glfw3.h>
#include<Events.h>
int initRenderer();
Event loop();
Event render();
int terminateRenderer();
#endif