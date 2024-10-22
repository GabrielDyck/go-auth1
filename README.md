## AUTH1 BACKEND

### Development environment setup


Set MYSQL_PASS environment value to establish connection to mysql.

Set SMTP_PASS environment value to send emails.

Set DOMAIN environment value for reset link.

Set ENVIRONMENT environment value to inform where the application is running.

Set PORT environment value to run server.
`
export MYSQL_PASS=${{MYSQL_PASS}}
export SMTP_PASS=${{SMTP_PASS}}
export DOMAIN=${{DOMAIN}}
export ENVIRONMENT=${{ENVIRONMENT}}

`

Run 'scripts/sql/script.sql' in mysql database.

## Front Routes
	http://localhost:80/signin
	http://localhost:80/signup
	http://localhost:80/edit-profile
	http://localhost:80/profile-info/{id}
	http://localhost:80/forgot-password
	http://localhost:80/reset-password
	
	
## Build Docker Image

`
docker build . -t gabrieldyck/auth1:latest
`

Publish
`
docker build . -t gabrieldyck/auth1:latest

docker push gabrieldyck/auth1:latest
`


## Run Container from Docker Image

`
sudo docker run --name auth1 --network host -d -v ${{AUTH1-BE_FOLDER}}/resources/application.json:/home/ubuntu/resources/application.json -e MYSQL_PASS=${{MYSQL_PASS}} -e SMTP_PASS=${{SMTP_PASS}} -e DOMAIN='localhost' -e ENVIRONMET='local'   gabrieldyck/auth1:$1

`

## Stop and Delete running container

`
docker container stop auth1 && docker rm auth1
`


## Log access
`
docker logs -f auth1
`
