#include<Renderer.h>
GLFWwindow *pWindow;
EventsStack eventsStack;
int initRenderer()
{
    eventsStack = newStack();
    if (!glfwInit())
    {
        return -1;
    }
    pWindow = glfwCreateWindow(1280, 720, "GO-PONG", NULL, NULL);
    if(!pWindow){
        terminateRenderer();
        return -2;
    }
    glfwMakeContextCurrent(pWindow);
    glfwPollEvents();
    glfwSwapBuffers(pWindow);
    return 0;
}
int terminateRenderer(){
    if(pWindow){
        glfwDestroyWindow(pWindow);
    }
    glfwTerminate();
    return 0;
}