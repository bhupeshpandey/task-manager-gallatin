# task-manager-gallatin
This project handles the task grpc requests and response from and to nashville

The code contained in the project starts a grpc server.
This grpc server will provide the base for the nashville service
to create, update, get, list and delete tasks.

Following are the steps to run the service:

1. Generate the pb files using:
   follow: https://grpc.io/docs/protoc-installation/
   protoc --go_out=. --go-grpc_out=. ./internal/proto/task_service.proto
2. Ensure the code is compiling by running the go build on the cmd folder.
   If everything is fine, then go to make file and run the command 'docker-compose-up'
   This will ensure the local build for the linux os is run along with the docker image
   generation for the application and finally the application is run on docker-compose.
   There is slight issue, that sometimes the application servie doesn't start right away.
   A manual start maybe needed.
3. Once, everything is good in above step, the grpc is running on the localhost:50051.
   The api's exposed are CreateTask, UpdateTask, GetTask, ListTasks and DeleteTask.
   Postman can be used to verify the API's.
   Assignment-Service can be looked up online on the postman for the grpc api.