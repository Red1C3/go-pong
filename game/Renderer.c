#include<Renderer.h>
GLFWwindow *pWindow;

int initRenderer()
{
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
    glfwSetWindowCloseCallback(pWindow, windowCloseCallback);
    return 0;
}
Event loop(){
    Event event = pop(getStack());
    if(event.code!=-1){
        return event;
    }
    return render();
}
Event render(){
    Event event;
    event.code = 0;
    glfwSwapBuffers(pWindow);
    glfwPollEvents();
    return event;
}
int terminateRenderer(){
    if(pWindow){
        glfwDestroyWindow(pWindow);
    }
    glfwTerminate();
    return 0;
}