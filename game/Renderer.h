#ifndef _H_RENDERER
#define _H_RENDERER
#include<stdio.h>
#include<GL/glew.h>
#include<GLFW/glfw3.h>
#include<cglm/struct.h>
#include<Events.h>
typedef struct{
    float p1, p2;
    float ball[2];
} DrawInfo;
int initRenderer();
Event loop(DrawInfo drawInfo);
int terminateRenderer();
#endif