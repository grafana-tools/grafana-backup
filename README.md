# Backup tool for Gafana

CLI for the simple backup/restore operations on Grafana dashboards and datasources.

**Work in progress.**

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
        
        
        
        
        
        
        
    
