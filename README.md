# blog

Personal blog.

## Process

Live Server:

```bash
hugo server -D
```

New posts: 

```bash
hugo new posts/title.md
```

Update `kodata`:

```bash
hugo -D --destination kodata/ --minify
```

Test deployment:

```bash
FILE_PATH=~/src/n3wscott/blog/kodata/ go run blog.go
```

Publish images:

```bash
ko publish github.com/n3wscott/blog
```
