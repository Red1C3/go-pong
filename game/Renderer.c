#include <Renderer.h>
GLFWwindow *pWindow;
GLuint shaderProgram;
GLuint quadVAO;
GLint modelMatLocation, ballBoolLocation;
const int sideOffset = 30;
const float stickLenScale = 6, stickWidthScale = 0.5f;
static inline GLuint createShaderProgram()
{
    FILE *vsFile = fopen("./Shader/shader.vert", "rb");
    fseek(vsFile, 0, SEEK_END);
    size_t vsSize = ftell(vsFile);
    const char *vertexShaderCode = malloc(sizeof(char) * vsSize);
    fseek(vsFile, 0, SEEK_SET);
    fread(vertexShaderCode, vsSize * sizeof(char), vsSize, vsFile);
    fclose(vsFile);
    FILE *fsFile = fopen("./Shader/shader.frag", "rb");
    fseek(fsFile, 0, SEEK_END);
    size_t fsSize = ftell(fsFile);
    const char *fragmentShaderCode = malloc(sizeof(char) * fsSize);
    fseek(fsFile, 0, SEEK_SET);
    fread(fragmentShaderCode, fsSize * sizeof(char), fsSize, fsFile);
    fclose(fsFile);
    GLuint vertexShader = glCreateShader(GL_VERTEX_SHADER);
    glShaderSource(vertexShader, 1, &vertexShaderCode, &vsSize);
    glCompileShader(vertexShader);
    GLuint fragmentShader = glCreateShader(GL_FRAGMENT_SHADER);
    glShaderSource(fragmentShader, 1, &fragmentShaderCode, &fsSize);
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
static inline GLuint createQuadVAO()
{
    GLuint VAO, VBO, EBO;
    float pos[] = {-1, -1,
                   1, -1,
                   -1, 1,
                   1, 1};
    unsigned int indices[] = {0, 1, 2,
                              3, 2, 1};
    glGenVertexArrays(1, &VAO);
    glBindVertexArray(VAO);
    glGenBuffers(1, &VBO);
    glGenBuffers(1, &EBO);
    glBindBuffer(GL_ARRAY_BUFFER, VBO);
    glBindBuffer(GL_ELEMENT_ARRAY_BUFFER, EBO);
    glBufferData(GL_ARRAY_BUFFER, sizeof(float) * 8, pos, GL_STATIC_DRAW);
    glBufferData(GL_ELEMENT_ARRAY_BUFFER, sizeof(unsigned int) * 6, indices, GL_STATIC_DRAW);
    glEnableVertexAttribArray(0);
    glVertexAttribPointer(0, 2, GL_FLOAT, GL_FALSE, 0, (void *)0);
    return VAO;
}
int initRenderer()
{
    if (!glfwInit())
    {
        return -1;
    }
    glfwWindowHint(GLFW_CONTEXT_VERSION_MAJOR, 3);
    glfwWindowHint(GLFW_CONTEXT_VERSION_MINOR, 3);
    glfwWindowHint(GLFW_OPENGL_PROFILE, GLFW_OPENGL_CORE_PROFILE);
    glfwWindowHint(GLFW_RESIZABLE, GLFW_FALSE);
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
    shaderProgram = createShaderProgram();
    if (glGetError() != 0)
    {
        terminateRenderer();
        return -4;
    }
    glUseProgram(shaderProgram);
    if (glGetError() != 0)
    {
        terminateRenderer();
        return -5;
    }
    float scaleFac = 4.0f;
    mat4s ortho = glms_ortho((-16.0f / 2.0f) * scaleFac,
                             (16.0f / 2.0f) * scaleFac,
                             (-9.0f / 2.0f) * scaleFac,
                             (9.0f / 2.0f) * scaleFac,
                             0.9f, 1.1f);
    mat4s view = glms_lookat((vec3s){.x = 0, .y = 0, .z = 1},
                             (vec3s){.x = 0, .y = 0, .z = 0},
                             (vec3s){.x = 0, .y = 1, .z = 0});
    mat4s VP = glms_mat4_mul(ortho, view);
    GLint VPLocation = glGetUniformLocation(shaderProgram, "VP");
    glUniformMatrix4fv(VPLocation, 1, GL_FALSE, &VP.raw);
    modelMatLocation = glGetUniformLocation(shaderProgram, "M");
    ballBoolLocation = glGetUniformLocation(shaderProgram, "isBall");
    if (glGetError() != 0)
    {
        terminateRenderer();
        return -6;
    }
    quadVAO = createQuadVAO();
    if (glGetError() != 0)
    {
        terminateRenderer();
        return -7;
    }
    return 0;
}
static inline Event render(DrawInfo *drawInfo)
{
    Event event;
    glClear(GL_COLOR_BUFFER_BIT);
    mat4s model = glms_translate(GLMS_MAT4_IDENTITY,
                                 (vec3s){.x = drawInfo->ball[0], .y = drawInfo->ball[1]});
    glUniformMatrix4fv(modelMatLocation, 1, GL_FALSE, &model.raw);
    glUniform1i(ballBoolLocation, 1);
    glDrawElements(GL_TRIANGLES, 6, GL_UNSIGNED_INT, 0);
    glUniform1i(ballBoolLocation, 0);
    model = glms_translate(GLMS_MAT4_IDENTITY,
                           (vec3s){.x = sideOffset, .y = drawInfo->p1});
    model = glms_scale(model,
                       (vec3s){.x = stickWidthScale, .y = stickLenScale});

    glUniformMatrix4fv(modelMatLocation, 1, GL_FALSE, &model.raw);
    glDrawElements(GL_TRIANGLES, 6, GL_UNSIGNED_INT, 0);
    model = glms_translate(GLMS_MAT4_IDENTITY,
                           (vec3s){.x = -sideOffset, .y = drawInfo->p2});
    model = glms_scale(model,
                       (vec3s){.x = stickWidthScale, .y = stickLenScale});

    glUniformMatrix4fv(modelMatLocation, 1, GL_FALSE, &model.raw);
    glDrawElements(GL_TRIANGLES, 6, GL_UNSIGNED_INT, 0);
    glfwSwapBuffers(pWindow);
    glfwPollEvents();
    event.code = glGetError(); //TODO reset before release
    return event;
}
Event loop(DrawInfo drawInfo)
{
    Event event = pop(getStack());
    if (event.code != -1)
    {
        return event;
    }
    return render(&drawInfo);
}
int terminateRenderer()
{
    if (pWindow)
    {
        glfwDestroyWindow(pWindow);
    }
    glfwTerminate();
    return 0;
}