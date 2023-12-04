# Intro

A script to run remote scripts in amazon ec2 instance.

## How to get started

Write your scripts in internal/scripts.txt file, comma separated commands.
Like the example below:

```text
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
    go install github.com/jkarage/ec2update/cmd/ec2@latest
```

#### Usage

``` bash
    ec2updates --region eu-west-3 --instance-id i-04493d3e5d3001e32
```

Providing Custom Script file

``` bash
   ec2updates --region eu-west-3 --instance-id i-04493d3e5d3001e32 --script a.txt

```
