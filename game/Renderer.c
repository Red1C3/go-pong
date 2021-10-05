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
    glewExperimental = 1;
    if(glewInit()!=GLEW_OK){
        terminateRenderer();
        return -3;
    }
    glClearColor(0, 0, 0, 1);
    glfwSetWindowCloseCallback(pWindow, windowCloseCallback);
    return 0;
}
Event loop(DrawInfo drawInfo){
    Event event = pop(getStack());
    if(event.code!=-1){
        return event;
    }
    return render(drawInfo);
}
Event render(DrawInfo drawInfo){
    Event event;
    glClear(GL_COLOR_BUFFER_BIT);
    glfwSwapBuffers(pWindow);
    glfwPollEvents();
    event.code = glGetError(); //TODO reset before release
    return event;
}
int terminateRenderer(){
    if(pWindow){
        glfwDestroyWindow(pWindow);
    }
    glfwTerminate();
    return 0;
}