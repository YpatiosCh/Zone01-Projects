## Building the App's Container Image
```
cd <app directory>
docker build -t go-webapp:1.0 .
```
## Starting App Container
```
docker run --name go-webapp-container -p 8080:8080 -d go-webapp:1.0 
```

OR

## Run the build script to build app container and image
``` 
./build.sh
```

## Removing a container using the CLI
```
docker ps
# Swap out <the-container-id> with the ID from docker ps
docker stop <the-container-id>
docker rm <the-container-id>
----OR---
docker rm -f <the-container-id>
```

## View image metadata 
- View image's labels
```
docker image inspect --format='{{json .Config.Labels}}' go-webapp:1.0
```

## Clear and grabage collection
```
./vacuum.sh
```

# IMPORTANT
- Since the image uses alpine you should use 
```
/bin/sh
```