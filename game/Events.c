/*MIT License

Copyright (c) 2021 Mohammad Issawi

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
*/

#include <Events.h>
EventsStack eventsStack;
void newStack()
{
    eventsStack.head = NULL;
    eventsStack.length = 0;
}
EventsStack *getStack()
{
    return &eventsStack;
}
void push(EventsStack *stack, Event event)
{
    Event *pEvent = malloc(sizeof(Event));
    pEvent->next = stack->head;
    pEvent->code = event.code;
    pEvent->key = event.key;
    stack->head = pEvent;
    stack->length++;
}
Event pop(EventsStack *stack)
{
    if (stack->length == 0)
    {
        Event event;
        event.code = -1;
        return event;
    }
    Event event;
    event.code = stack->head->code;
    event.key = stack->head->key;
    Event *next = stack->head->next;
    Event *cur = stack->head;
    stack->head = next;
    free(cur);
    stack->length--;
    return event;
}
void windowCloseCallback(GLFWwindow *window)
{
    Event event;
    event.code = 1;
    push(getStack(), event);
}