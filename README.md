# glookbs
It's a restful task API application, which includes the following endpoints:

 - `GET` /tasks
 - `POST` /tasks
 - `PUT` /tasks/{id}
 - `DELETE` /tasks/{id}

A `task` should contain at least the following fields:
 - `name`
   - type: `string`
   - description:task name
 - `status`
    - type: `integer`
    - enum:[0,1]
    - description:
      - `0` represents an incomplete task, while 
      - `1` represents a completed task

**Requirements**:
 - Runtime environment should be Go 1.18+
 - Provides unit tests
 - Provides Dockerfile to run API in Docker
 - Manage the codebase on Github and provide us with the repository link
 - For data storage, you can use any in-memory mechanism

# Run

`docker-compose up`