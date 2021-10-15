# Go-Pong
A Pong video game clone made with Go lang and OpenGL 3.3 using C.
## Gameplay
### Offline
Key bindings are 'w' and 's' for the left player and 'up arrow' and 'down arrow' for the right player, score is kept in terminal (for now...)
### Multiplayer
go-pong supports multiplayer by hosting a game and joining it, hosting can be performed by running ```./go-pong -h <port> ``` , and connecting to that host can be performed by running ```./go-pong <ip:address:port> ``` , control your player with up and down arrows.
## Building
After installing the required dependencies, run ```go build```, make sure your executable is in the same folder than contains the *Shader* folder.
## Roadmap
Also a real score board might be added in the future.
## Contributions
Are welcome !
## License
Under MIT, check [License](LICENSE)

### Created with fuzzy kittens, with the help of RedDeadAlice