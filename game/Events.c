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