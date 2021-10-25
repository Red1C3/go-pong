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

#include <Renderer.h>
GLFWwindow *pWindow;
GLuint shaderProgram;
GLuint quadVAO;
GLint modelMatLocation, ballBoolLocation, scoreBoolLocation;
GLint textures[10];
bool canSend = true;
bool isOnline;
double sendTime;
const int sideOffset = 30;
const float stickLenScale = 6, stickWidthScale = 0.5f;
static inline void handleInput()
{
    if (isOnline && !canSend)
    {
        return;
    }
    if (glfwGetKey(pWindow, GLFW_KEY_S) == GLFW_PRESS)
    {
        push(getStack(), (Event){.code = 2, .key = 's'});
    }
    if (glfwGetKey(pWindow, GLFW_KEY_W) == GLFW_PRESS)
    {
        push(getStack(), (Event){.code = 2, .key = 'w'});
    }
    if (glfwGetKey(pWindow, GLFW_KEY_UP) == GLFW_PRESS)
    {
        push(getStack(), (Event){.code = 2, .key = 'u'});
    }
    if (glfwGetKey(pWindow, GLFW_KEY_DOWN) == GLFW_PRESS)
    {
        push(getStack(), (Event){.code = 2, .key = 'd'});
    }
    if (isOnline)
    {
        canSend = false;
        sendTime = glfwGetTime();
    }
}
//loads shaders
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
//creates a single quad VAO to be used for all game objects
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
static inline GLuint loadTexture(const char *path)
{
    int width, height;
    unsigned char *texData = stbi_load(path, &width, &height, NULL, 0);
    if (texData == NULL)
    {
        return 0;
    }
    GLuint tex;
    glGenTextures(1, &tex);
    glBindTexture(GL_TEXTURE_2D, tex);
    glTexParameteri(GL_TEXTURE_2D, GL_TEXTURE_MIN_FILTER, GL_NEAREST);
    glTexParameteri(GL_TEXTURE_2D, GL_TEXTURE_MAG_FILTER, GL_NEAREST);
    glTexImage2D(GL_TEXTURE_2D, 0, GL_RGBA8, width, height, 0, GL_RGBA, GL_UNSIGNED_BYTE, texData);
    stbi_image_free(texData);
    if (glGetError() != 0)
    {
        return 0;
    }
    return tex;
}
int initRenderer(bool online)
{
    isOnline = online;
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
    char buffer[17];
    for (int i = 0; i < 10; ++i)
    {
        sprintf(buffer, "./Textures/%d.png", i);
        textures[i] = loadTexture(buffer);
        if (textures[i] == 0)
        {
            terminateRenderer();
            return -6;
        }
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
    scoreBoolLocation = glGetUniformLocation(shaderProgram, "isScore");
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
                           (vec3s){.x = -sideOffset, .y = drawInfo->p1});
    model = glms_scale(model,
                       (vec3s){.x = stickWidthScale, .y = stickLenScale});

    glUniformMatrix4fv(modelMatLocation, 1, GL_FALSE, &model.raw);
    glDrawElements(GL_TRIANGLES, 6, GL_UNSIGNED_INT, 0);
    model = glms_translate(GLMS_MAT4_IDENTITY,
                           (vec3s){.x = sideOffset, .y = drawInfo->p2});
    model = glms_scale(model,
                       (vec3s){.x = stickWidthScale, .y = stickLenScale});

    glUniformMatrix4fv(modelMatLocation, 1, GL_FALSE, &model.raw);
    glDrawElements(GL_TRIANGLES, 6, GL_UNSIGNED_INT, 0);
    glUniform1i(scoreBoolLocation, 1);
    for (int i = 0; i < 2; ++i)
    {
        if (drawInfo->scores[i] > 9)
            continue;
        glBindTexture(GL_TEXTURE_2D, textures[drawInfo->scores[i]]);
        model = glms_translate(GLMS_MAT4_IDENTITY,
                               (vec3s){.x = (i % 2) ? 5 : -5, .y = 15});
        model = glms_scale(model,
                           (vec3s){.x = 2, .y = 2});
        glUniformMatrix4fv(modelMatLocation, 1, GL_FALSE, &model.raw);
        glDrawElements(GL_TRIANGLES, 6, GL_UNSIGNED_INT, 0);
    }
    glUniform1i(scoreBoolLocation, 0);
    glfwSwapBuffers(pWindow);
    glfwPollEvents();
    event.code = glGetError(); //TODO reset before release
    return event;
}
Event loop(DrawInfo drawInfo)
{
    if (isOnline && !canSend && glfwGetTime() - sendTime >= 1.0 / 120.0)
    {
        canSend = true;
    }
    Event event = pop(getStack());
    if (event.code != -1)
    {
        return event;
    }
    handleInput();
    return render(&drawInfo);
}
int terminateRenderer()
{
    if (pWindow)
    {
        glfwDestroyWindow(pWindow);
    }
    glfwTerminate();
    while (pop(getStack()).code != -1)
        ;
    return 0;
}