# Cron Job for golang


## Quick Start

```shell
go run main.go
```

## Change schedule 

In line Code `137-139` 
- You can change the schedule for 3 times a day

```
gocron.Every(1).Day().At("xx:xx").Do(job)
gocron.Every(1).Day().At("xx:xx").Do(job)
gocron.Every(1).Day().At("xx:xx").Do(job)
```