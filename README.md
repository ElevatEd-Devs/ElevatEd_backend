<h1> ElevatEd Backend </h1>

<h2> Getting started 🧑‍💻</h2>
<p>Make sure you have go installed: <a href="https://go.dev/doc/install">Install Go</a></p>
<p>This project would use v1.24.5</p>

<p>Clone the repository:</p>

```
git clone https://github.com/ElevatEd-Devs/ElevatEd_backend
```

<p>Make sure to set up the .env file: for more info send a message to the backend channel on discord</p>

<p>For live reloading make sure your have <a href="https://github.com/air-verse/air">air</a> installed</p>

```
go install github.com/air-verse/air@latest
```

<p>Now that your have air installed run the following in the project directory:</p>

```
air
```

<p>This would start the server and listen for changes in the program</p>

<p>You could  also just run the server normally with:</p>

```
go run main.go
```

<h4>For Ubuntu users</h4>
<p> You may face some challenges with running air on different terminals, here is a work around: </p>

```
export PATH=$PATH:$HOME/go/bin
source ~/.bashrc
air
```

<h2>Using Docker <img width="48" height="48" src="https://img.icons8.com/color/48/docker.png" alt="docker"/></h2>
<p>A quick demonstration into how <a href="https://www.docker.com/get-started/">docker</a> might be used.</p>
<p>This would build a docker image for the project called "elevated_backend", then run the image.</p>
<p>Note: in this example the host's 3000 port is mapped to the container's 3000 port.</p>

```
docker build -t elevated_backend .
docker run -p 3000:3000 --rm -v $(pwd):/app -v /app/tmp --name docker-air elevated_backend
```

<h2>Contributing ✍️</h2>
<p>For each issue that you pick up, create a new branch and work on it.</p>  
<img width="273" height="108" alt="image" src="https://github.com/user-attachments/assets/bb721a49-2d2d-44e2-b26a-d2e10ca01036" />

<p>Once done with the task, commit your changes, push, make a pull request and have your code peer reviewed.</p>

<p>Make sure to alert the team on discord so the review begins as soon as possible.</p>

<h2>Testing EndPoints🧪</h2>
<p>You could test APIs with: <a href="https://www.postman.com/">Postman</a> or <a href="https://www.usebruno.com/">Bruno</a> </p>
<p>There is also a <a href="https://learn.microsoft.com/en-us/aspnet/core/test/http-files?view=aspnetcore-9.0">.http</a> file you would find in this repository, just hit the send request button</p>
<img width="100" height="22" alt="image" src="https://github.com/user-attachments/assets/c9a77824-75f4-438a-9366-96309653422a" />

<h2>More Documentation</h2>
<p>This repository also has more detailed documentation about various features:</p>

- [Auth Docs](docs/authEndPoints.md)
