#include <Renderer.h>
GLFWwindow *pWindow;
GLuint shaderProgram;
static inline GLuint createShaderProgram()
{
    FILE *vsFile = fopen("./Shader/shader.vert", "rb");
    fseek(vsFile, 0, SEEK_END);
    size_t size = ftell(vsFile);
    const char *vertexShaderCode = malloc(sizeof(char) * size);
    fseek(vsFile, 0, SEEK_SET);
    fread(vertexShaderCode, size * sizeof(char), size, vsFile);
    fclose(vsFile);
    FILE *fsFile = fopen("./Shader/shader.frag", "rb");
    fseek(fsFile, 0, SEEK_END);
    size = ftell(fsFile);
    const char *fragmentShaderCode = malloc(sizeof(char) * size);
    fseek(fsFile, 0, SEEK_SET);
    fread(fragmentShaderCode, size * sizeof(char), size, fsFile);
    fclose(fsFile);
    GLuint vertexShader = glCreateShader(GL_VERTEX_SHADER);
    glShaderSource(vertexShader, 1, &vertexShaderCode,NULL);
    glCompileShader(vertexShader);
    GLuint fragmentShader = glCreateShader(GL_FRAGMENT_SHADER);
    glShaderSource(fragmentShader, 1, &fragmentShaderCode, NULL);
    glCompileShader(fragmentShader);
    GLuint program = glCreateProgram();
    glAttachShader(program, vertexShader);
    glAttachShader(program, fragmentShader);
    glLinkProgram(program);
    glDetachShader(program, vertexShader);
    glDetachShader(program, fragmentShader);
    glDeleteShader(vertexShader);
    glDeleteShader(fragmentShader);
    free(fragmentShaderCode);
    free(vertexShaderCode);
    return program;
}
int initRenderer()
{
    if (!glfwInit())
    {
        return -1;
    }
    pWindow = glfwCreateWindow(1280, 720, "GO-PONG", NULL, NULL);
    if (!pWindow)
    {
        terminateRenderer();
        return -2;
    }
    glfwMakeContextCurrent(pWindow);
    glfwSetWindowCloseCallback(pWindow, windowCloseCallback);
    glewExperimental = 1;
    if (glewInit() != GLEW_OK)
    {
        terminateRenderer();
        return -3;
    }
    glClearColor(0, 0, 0, 1);
    shaderProgram=createShaderProgram();
    if(glGetError()!=0){
        terminateRenderer();
        return -4;
    }
    return 0;
}
Event loop(DrawInfo drawInfo)
{
    Event event = pop(getStack());
    if (event.code != -1)
    {
        return event;
    }
    return render(drawInfo);
}
Event render(DrawInfo drawInfo)
{
    Event event;
    glClear(GL_COLOR_BUFFER_BIT);
    glfwSwapBuffers(pWindow);
    glfwPollEvents();
    event.code = glGetError(); //TODO reset before release
    return event;
}
int terminateRenderer()
{
    glDeleteProgram(shaderProgram);
    if (pWindow)
    {
        glfwDestroyWindow(pWindow);
    }
    glfwTerminate();
    return 0;
}