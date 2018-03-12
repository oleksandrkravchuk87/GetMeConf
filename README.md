
[![Build Status](https://travis-ci.org/YAWAL/GetMeConf.svg?branch=master)](https://travis-ci.org/YAWAL/GetMeConf)

Config service
==============


This is a simple config service, which allows basic CRUD operations for different configs. Configs are stored in a Postgres database.
gRPC is used to communicate with the service.

  


How to start

To install dep  dependency management tool run 

``````````````````
make install dep
``````````````````

To install application dependencies run


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

To to run the application in a docker container

``````````````````
make docker-build
``````````````````

