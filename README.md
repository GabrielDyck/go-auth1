## AUTH1 BACKEND

### Development environment setup


Set MYSQL_PASS environment value to establish connection to mysql.

Set SMTP_PASS environment value to send emails.

`
export MYSQL_PASS=${{MYSQL_PASS}}
export SMTP_PASS=${{SMTP_PASS}}
`

Run 'scripts/sql/script.sql' in mysql database.

## Front Routes
	http://localhost:8080/signin
	http://localhost:8080/signup
	http://localhost:8080/edit-profile
	http://localhost:8080/profile-info/{id}
	http://localhost:8080/forgot-password
	http://localhost:8080/reset-password