# immufile
CLI tool to verify if the file or line in file was tampered

## Usage

---
```shell
immufile -f <file-path>
```

Returns SHA of the file.

---
```shell
immufile -f <file-path> -l 30 -c "file content" -d <file-SHA>
```
Returns if the provided file content (-c) at line (-l) was tampered or not basing on file hash (-d).

---
```shell
immufile -f <file-path> -l 30 -c "file content" -d <file-SHA> --short
```
Returns if the provided file content (-c) at line (-l) was tampered or not basing on file hash (-d).
Results could `tampered|verified`