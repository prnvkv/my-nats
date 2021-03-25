## My-NATS

My-NATS is a simple repository with the basic messaging patterns using NATS.

Example methods to call the messaging interfaces are written in `cmd` directory with all the three messaging patterns. The configuration details about the nats-server are provided via the configuration file `deploy` directory. 

As a pre-requisite, the nats-server needs to be running in the machine to test it out locally. 
Execute the below command to pull and run the image locally.

    `docker run -d --name nats-main -p 4222:4222 -p 6222:6222 -p 8222:8222 nats`

Once the above command is successful, you can try out the messaging patterns in the `cmd` directory using the binaries already available or you can create your own example/application using the interface methods in `pkg` directory.

