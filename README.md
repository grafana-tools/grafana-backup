# Backup tool for Gafana

CLI for the simple backup/restore operations on Grafana dashboards and datasources.
It based on [autograf](https://github.com/grafov/autograf) library for Grafana.

**Work in progress. Current state: it works partially.**

## Examples

	$ grafana-backup -url http://127.1:3000 -key xxxxxxxx -tag tag1,tag2 -title "sample api"" ls
		
	$ grafana-backup -url http://127.1:3000 -key xxxxxxxx -tag tag1,tag2 -title sample backup
	
	$ grafana-backup ls-files
	
	$ grafana-backup -url http://127.1:3000 -key xxxxxxxx -tag tag1 restore

    $ grafana-backup

        Backup tool for Grafana.
        Copyright (C) 2016  Alexander I.Grafov <siberian@laika.name>
        
        This program comes with ABSOLUTELY NO WARRANTY.
        This is free software, and you are welcome to redistribute it
        under conditions of GNU GPL license v3.
        
        Usage: $ grafana-backup [flags] <command>
        
        Available commands are: backup, restore, list, info, config, help.
        Call 'grafana-backup help <command>' for details about the command.

          -file string
            	use only listed files (file masks allowed)
          -key string
            	API key of Grafana server
          -name string
            	dashboard should match name
          -starred
            	only match starred dashboards
          -tag string
            	dashboard should match all these tags
          -timeout duration
        	read flagTimeout for interacting with Grafana (seconds) (default 6m0s)
          -url string
            	URL of Grafana server
          -v	verbose output	
        
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
