# go-ata-nop

`go-ata-nop` is a small program that periodically sends ATA NOP command to SATA disks.

## Motivation
I bought a Synology DS223j and two WD80EAAZ drives.

The WD80EAAZ comes with IntelliPark, but it cannot be disabled using conventional methods.

To circumvent IntelliPark, I periodically send ATA NOP command to the drives.

Most of the program, of course, was written by generative AI.

This should prevent the Load_Cycle_Count from continuously increasing.

## go-ata-nop

sends ATA NOP command periodically

```
root@ds223j:~# ./go-ata-nop -verbose /dev/sata1 /dev/sata2
2024/10/20 16:20:54.324578 Sending NOP command to /dev/sata1
2024/10/20 16:20:54.435479 NOP command sent to /dev/sata1
2024/10/20 16:20:54.435576 Sending NOP command to /dev/sata2
2024/10/20 16:20:54.547928 NOP command sent to /dev/sata2
```
