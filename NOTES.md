## To release

```bash
git tag v1.0.0
git push origin v1.0.0
gh release create v1.0.0 --title "v1.0.0" --notes "Initial release"
```

To warm the Go module proxy, run:

```bash
curl https://proxy.golang.org/github.com/schmidtw/civil/@v/v1.0.0.info
```