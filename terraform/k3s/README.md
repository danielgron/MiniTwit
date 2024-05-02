# terraform for setting up k3s


Required secrets/tokens:

DO read/write token - for creating droplets
DO read/update token - for use with limited rights to read/update reserved ip in case of node going down.
Simply API key - for updating dns record


Self chosen:

k3s token - used between nodes for creating cluster
rancher pw - password for the admin user in web interface
