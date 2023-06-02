# go-git-fixtures

git repository fixtures used by [go-git](https://github.com/go-git/go-git)

## Adding new Fixtures

### Adding new pack fixtures

1. Get the `.idx` and `.pack` files from the repository:

```sh
git clone https://<repository>
cd <repository_name>
git checkout <REF>
git gc

ls .git/objects/pack/
```

2. Copy them into `/data`.
3. Add a new entry in `fixtures.go`:

```
{
	Tags:         []string{"packfile", "<TAG_TO_REFER_TO>"},
	PackfileHash: "<PACK_HASH>",
}
```

4. Run `make generate`.
