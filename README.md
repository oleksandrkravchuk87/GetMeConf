
Config service
==============


This is a simple config service, which allows basic CRUD operations for different configs. Configs are stored in a Postgres database.
gRPC is used to communicate with the service.

  


How to start

Run to install dep

``````````````````
make install dep
``````````````````

Run to install application dependencies


``````````````````
make dependencies
``````````````````

To run tests

``````````````````
make tests
``````````````````

To to build the application

``````````````````
make build
``````````````````

To to run the application

``````````````````
make run
``````````````````

To to run the application in docker container

``````````````````
make docker-build
``````````````````
