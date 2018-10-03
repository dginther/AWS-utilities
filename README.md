# AWS-utilities

## iam-rotation-check/

Language: Golang

AWS best practices say to rotate your keys every 90 days.
This script will check a USERNAME in your AWS account and email RECIPIENT if it needs to be rotated.
Any SMTP server should work, tested with Gmail and an app password.

usage:
```
iam-rotation-check -user USERNAME -sender you@example.com -rcpt RECIPIENT 
                   -smtpServer smtp.gmail.com -smtpPort 587 -smtpPassword PASSWORD
```


arguments:
```
-user USERNAME          AWS Username to check IAM Access Keys for.
-sender EMAIL           Email address for sender (SMTP server username)
-rcpt EMAIL             Email address for the recipient (User of AWS account)
-smtpServer SERVER      SMTP server hostname
-smtpPort PORT          SMTP port
-smtpPassword PASSWORD  SMTP password
-send true/false	Send email? (default: false)
```
