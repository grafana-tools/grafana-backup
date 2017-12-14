# Backup tool for Gafana

CLI for the simple backup/restore operations of [Grafana](https://grafana.com/) dashboards and datasources.
It uses [Grafana client SDK](https://github.com/grafana-tools/sdk).

**Work in progress. Current state: it may works, may not. Depends on build. Don't use it yet!**

## Examples

Back up all dashboards and any datasources which they require. Backups will be saved as separate files in 
the directory $CWD/backup

```$ grafana-backup -url http://127.1:3000 -key xxxxxxxx backup```

Show all dashboards which have the flags "tagone" and "tagtwo" applied and any datasources they use.

```
$ grafana-backup -tag tagone,tagtwo ls
<31> "Test Dashboard 1" [tagone tagtwo] 
<10> "prometheus-test" (prometheus) http://prometheus-test.example.com:9090
<8> "Promt2Local" (prometheus) http://prometheus-poc.example.com:9090
```

Show the dashboard titled "Test Dashboard 3" and the datasources it uses.

```
$ grafana-backup -title 'Test Dashboard 3' ls
<33> "Test Dashboard 3" [] 
<10> "prometheus-test" (prometheus) http://prometheus-test.example.com:9090
<8> "Promt2Local" (prometheus) http://prometheus-poc.example.com:9090
```

Back up all dashboards which have the flags "tag1" and "tag2" applied and any datasources they use.

```
$ grafana-backup -url http://127.1:3000 -key xxxxxxxx -tag tag1,tag2 backup
```
		
Back up a dashboard called "sample"

```
$ grafana-backup -url http://127.1:3000 -key xxxxxxxx -tag tag1,tag2 -title sample backup
```
	
Back up all dashboards, datasources and users

```
$ grafana-backup -url http://127.1:3000 -key xxxxxxxx -apply-for all
```

Show information about local backup files

```
$ grafana-backup ls-files
test-dashboard-1.db.json:	 board id:31 "Test Dashboard 1" [tagone tagtwo]
test-dashboard-2.db.json:	 board id:32 "Test Dashboard 2"
test-dashboard-3.db.json:	 board id:33 "Test Dashboard 3"

``` 

Restore all local objects which have the tag "tag1" applied to them. 

```
$ grafana-backup -url http://127.1:3000 -key xxxxxxxx -tag tag1 restore
```

View the usage

```
$ grafana-backup 
Backup tool for Grafana.
Copyright (C) 2016-2017  Alexander I.Grafov <siberian@laika.name>

This program comes with ABSOLUTELY NO WARRANTY.
This is free software, and you are welcome to redistribute it
under conditions of GNU GPL license v3.

Usage: $ grafana-backup [flags] <command>

Available commands are: backup, restore, ls, ls-files, info, config, help.
Call 'grafana-backup help <command>' for details about the command.

  -apply-for string
    	apply operation only for some kind of objects, available values are "auto", "all", "dashboards", "datasources", "users" (default "auto")
  -dir string
    	A directory to write backup files to or read them from. (default "backup")
  -force
    	force overwrite of existing objects
  -key string
    	API key of Grafana server
  -starred
    	only match starred dashboards
  -tag string
    	dashboard should match all these tags
  -timeout duration
    	read flagTimeout for interacting with Grafana in seconds (default 6m0s)
  -title string
    	dashboard title should match name
  -url string
    	URL of Grafana server
  -verbose
    	verbose output
```
    
        
## List of proposed commands, flags and args

Draft and it is subject for changes.

	# List dashboards.
	$ grafana-backup -key=xxxx -url=x.y.z -title=match-name -tag=x,y,z ls

	# List datasources.
	$ grafana-backup -key=xxxx -url=x.y.z ls-ds

	# List users.
	$ grafana-backup -key=xxxx -url=x.y.z ls-users

	# Do backup for matching dashboards.
	$ grafana-backup -key=xxxx -url=x.y.z -title=match-name -tag=x,y,z backup path/to

	# Restore objects on a server at url only for boards match tags.
	$ grafana-backup -key=xxxx -url=x.y.z -tag x,y,z restore path/from

	# List local backups for tags and file mask
	$ grafana-backup -tag x,y,z -file "backup/*/*" ls-files 

	# Save all flags to config var.
	$ grafana-backup -key=xxxx -url=x.y.z config-set confname

	# Get flag values for config variable.
	$ grafana-backup config-get confname

	# Flag applied for backup/restore
	-objects=auto,dashboards,datasources,users,all
