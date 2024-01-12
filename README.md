# EC2Update

A script to run remote scripts in amazon ec2 instance.

## Prerequistes

- Get the instanceid of the machine
- Get the `AWS_ACCESS_KEY_ID` AND `AWS_SECRET_KEY` values 

## How to get started

Write your scripts in internal/scripts.txt file, comma separated commands.
Like the example below:

``` bash

    touch /home/ubuntu/a.txt,
    echo 'Hello from SSM command!' > /home/ubuntu/a.txt,
    pwd,
    ls -lah /home/ubuntu

```

Add your AWS_ACCESS_KEY_ID and AWS_SECRET_ACCESS_KEY in the environment.

```bash

    export AWS_SECRET_ACCESS_KEY=**************
    export AWS_ACCESS_KEY_ID=******************

```

### Installing Requirements

``` bash

    git clone https://github.com/Jkarage/EC2update.git

```

#### Usage

``` bash

    go run cmd/ec2/main.go --region eu-west-3 --instance-id i-04493d3e5d3001e32

```

Providing Custom Script file

``` bash

   go run cmd/ec2/main.go --region eu-west-3 --instance-id i-04493d3e5d3001e32 --script a.txt

```
