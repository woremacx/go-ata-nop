# snippets

This directory contains small programs for trial and error.

## counts

read `/proc/diskstats` and show the number of read/write operations.

```
root@ds223j:~# ./counts sata1 sata2
[2024-10-20 14:51:47] sata1 (r:     1 w:     0) sata2 (r:     1 w:     0)
[2024-10-20 14:51:52] sata2 (r:     1 w:     0) sata1 (r:     1 w:     0)
[2024-10-20 14:51:57] sata1 (r:     1 w:     0) sata2 (r:     1 w:     0)
```
