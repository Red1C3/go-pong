#ifndef _H_EVENTS
#define _H_EVENTS
#include<stddef.h>
#include<stdlib.h>
typedef struct {
    void *next;
    int code;
    char key;
} Event;
typedef struct{
    Event *head;
    int length;
} EventsStack;
EventsStack newStack();
void push(EventsStack *stack, Event event);
Event pop(EventsStack *stack);
#endif