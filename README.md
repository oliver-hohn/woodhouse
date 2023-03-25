Inspired by Woodhouse from the Archer TV series. This repo is a collection of scripts for tasks which are easier to script than to do manually.

## Woodhouse, organize my files!
With the release binary:
```
$ ./woodhouse.bin organize /path/to/directory/with/unorganized/files/ /path/to/organized/directory/
Copied /path/to/directory/with/unorganized/files/foo.md to /path/to/organized/directory/2018/jan_to_march/foo.md
Copied /path/to/directory/with/unorganized/files/bar/lorem.ipsum to /path/to/organized/directory/2018/jan_to_march/lorem.ipsum
```

In development:
```
$ go run main.go organize path/to/directory/with/unorganized/files/ path/to/organized/directory/
Copied /path/to/directory/with/unorganized/files/foo.md to /path/to/organized/directory/2018/jan_to_march/foo.md
Copied /path/to/directory/with/unorganized/files/bar/lorem.ipsum to /path/to/organized/directory/2018/jan_to_march/lorem.ipsum
```

_Optionally, the `--dryrun=1` option can be provided to **only** print the copies that would be executed, without executing them._