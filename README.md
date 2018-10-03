# AWS-utilities

iam-rotation-check/

Language: Golang
AWS best practices say to rotate your keys every 90 days.
This script will check a username in your AWS account and email RECIPIENT if it needs to be rotated.
Any SMTP server should work, tested with Gmail and an app password.

```usage:
iam-rotation-check -user USERNAME -sender you@example.com -rcpt RECIPIENT 
                   -smtpServer smtp.gmail.com -smtpPort 587 -smtpPassword PASSWORD
```
