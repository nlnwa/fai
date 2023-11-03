# First Article Inspection (FAI)

FAI loops trough files in a given directory that match a given pattern
and for every matching file:

- creates a checksum file
- validates the file as a WARC-file
- logs the name, size, checksum and validation status
- updates a file size histogram metric and a validation error counter metric
- moves the file and the corresponding checksum file to a target directory

```text
Usage of fai:
  -concurrency int
        number of concurrent files processed (default [number of CPU cores])
  -invalid-target-dir string
        path to target directory where invalid files and their corresponding checksum files will be moved to
  -metrics-port int
        port to expose metrics on (default 8081)
  -pattern string
        glob pattern used to match filenames in source directory (default "*")
  -sleep duration
        sleep duration between directory listings, set to 0 to only do a single run (default 5s)
  -source-dir string
        path to source directory
  -tmp-dir string
        path to directory where temporary buffer files will be stored
  -valid-target-dir string
        path to target directory where valid files and their corresponding checksum files will be moved to
```