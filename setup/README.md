## Docker discovers :

Docker discovers plugins by looking for them in the plugin directory whenever a user or container tries to use one by name.

There are three types of files which can be put in the plugin directory.

    .sock files are UNIX domain sockets.
    .spec files are text files containing a URL, such as unix:///other.sock or tcp://localhost:8080.
    .json files are text files containing a full json specification for the plugin.

### Socket dependencies :
   UNIX domain socket files must be located under  `/run/docker/plugins` , whereas spec files can be located either under 
    /etc/docker/plugins or `/usr/lib/docker/plugins`. 
    Docker always searches for unix sockets in `/run/docker/plugins` first. 
    You can define each plugin into a separated subdirectory if you want to isolate definitions from each other. 
     For example, you can create the Openebsdriver socket under ```/run/docker/plugins/openebs/openebsdriver.sock``` and 
    only mount ```/run/docker/plugins/openebs``` inside the openebsdriver container.
    

### JSON specification

This is the JSON format for a plugin:

- ![openebsdriver.json](https://github.com/maheshreddy7797/openebsdriver/blob/master/config.json)

### Json file :
  It checks for spec or json files under /etc/docker/plugins and /usr/lib/docker/plugins if the socket doesn’t exist.
  
  
### Systemd socket activation

Plugins may also be socket activated by systemd. The official Plugins helpers natively supports socket activation. In order for a plugin to be socket activated it needs a service file and a socket file.

The service file (for example /lib/systemd/system/openebsdriver.service)

- ![openebsdriver.service](https://github.com/maheshreddy7797/openebsdriver/blob/master/setup/openebsdriver.service)


The socket file (for example /lib/systemd/system/openebsdriver.socket):

- ![openebsdriver.socket](https://github.com/maheshreddy7797/openebsdriver/blob/master/setup/openebsdriver.socket)



   
