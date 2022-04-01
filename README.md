# Unofficial GitHub Cli Tool

![Generic badge](https://img.shields.io/badge/Release-v0.1.0-green.svg)

`ghlog` is the unofficial GitHub CLI tool that contains the most-used functions on GitHub.

![code](https://user-images.githubusercontent.com/4322099/154396746-6ce6671b-4f7e-4668-8cf1-67f767aca1b1.png)

## Getting Started.

1. Build `ghlog` using `go` by the below command.

```bash
$ go build ghlog.go
```

2. Move your application to wherever your want, like `/usr/bin` on `Linux` or `C:/Program Files/ghlog` on `Windows`, and add the environment path of the `ghlog`.

3. Start your `ghlog` using the below command.

```bash
$ ghlog
```

4. If login is required, you have to login with your `user id` and `access token`.

It will be saved on `.auth` file in `Home Directory`.

## Commands that you can play with.

### repo/org

You can print all of your repositories using the below command.

```bash
$ ghlog repo
```

Likes `repo`, in addition, you can also print all organizations you joined using the below command.

```bash
$ ghlog org
```

### create repo

It is also possible to create a new repository using the following command.

```bash
$ ghlog create repo
```

It will ask you about `Repository name` and `Visibility` of your repository.

### heatmap

`Heatmap` prints a GitHub contribution chart that you can see on your GitHub profile in command-line interface.

You can simply print your heatmap with the below command.

```bash
$ ghlog heatmap
```

You can also check other users' heatmap just by adding `username`.

```bash
$ ghlog heatmap devappmin
```

### search

You can even search repositories from GitHub!

```bash
$ ghlog search <query> [from] [to]
```

## Libraries

[Go Query](https://github.com/PuerkitoBio/goquery)

[go-github](https://github.com/google/go-github/v41/github)

[oauth2](https://golang.org/x/oauth2)

[term](https://golang.org/x/term)

## License

```
MIT License

Copyright (c) 2022 Kim Seung Hwan

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
```
