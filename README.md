Go Program to learn coumunication between python module and Go using watcher folder. 

Code Flow :- 

1. Read configuration parameters from the config file using viper module
2. Set a watcher for a folder using the module fsnotify
3. Watch for any file change events
4. Read the file and check for matching values from regular expression
5. Send matching content to console for display.

Dependency modules :- 
1. fmt
2. path/filepath
3. viper
4. fsnotify

To Run :- 
Start the main program and keep it watch mode. 
Start the Python module to read the pdf
