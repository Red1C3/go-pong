#ifndef _H_EVENTS
#define _H_EVENTS
#include<stddef.h>
#include<stdlib.h>
#include<GLFW/glfw3.h>
typedef struct {
    void *next;
    int code;
    char key;
} Event;
typedef struct{
    Event *head;
    int length;
} EventsStack;
void newStack();
EventsStack *getStack();
void push(EventsStack *stack, Event event);
Event pop(EventsStack *stack);
void windowCloseCallback(GLFWwindow *window);
#endif