package main

import (
	"crypto/tls"
	"flag"
	"fmt"
	"net/smtp"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/iam"
)

type Mail struct {
	senderId string
	toId     string
	subject  string
	body     string
}

type SmtpServer struct {
	host string
	port string
}

func (s *SmtpServer) ServerName() string {
	return (s.host + ":" + s.port)
}

func getAccessKeys(user string) ([]*iam.AccessKeyMetadata, error) {
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String("us-east-1")},
	)
	if err != nil {
		return []*iam.AccessKeyMetadata{}, err
	}

	svc := iam.New(sess)

	result, err := svc.ListAccessKeys(&iam.ListAccessKeysInput{
		MaxItems: aws.Int64(5),
		UserName: aws.String(user),
	})
	if err != nil {
		panic("Unable to get access keys.")
	}

	return result.AccessKeyMetadata, nil
}

func checkCreateDate(date time.Time) time.Duration {
	return time.Since(date)
}

func (mail *Mail) BuildMessage() string {
	message := ""
	message += fmt.Sprintf("From: %s\r\n", mail.senderId)
	message += fmt.Sprintf("To: %s\r\n", mail.toId)
	message += fmt.Sprintf("Subject: %s\r\n", mail.subject)
	message += "\r\n" + mail.body

	return message
}

func sendSmtpEmail(server string, port string, password string, mail Mail) error {
	smtpServer := SmtpServer{host: server, port: port}
	m := mail
	messageBody := m.BuildMessage()
	tlsConfig := &tls.Config{
		InsecureSkipVerify: true,
		ServerName:         smtpServer.host,
	}
	auth := smtp.PlainAuth("", m.senderId, password, smtpServer.host)

	c, err := smtp.Dial(smtpServer.ServerName())
	if err != nil {
		return err
	}

	c.StartTLS(tlsConfig)

	if err = c.Auth(auth); err != nil {
		return err
	}

	if err = c.Mail(m.senderId); err != nil {
		return err
	}

	if err = c.Rcpt(m.toId); err != nil {
		return err
	}

	wc, err := c.Data()
	if err != nil {
		return err
	}

	if _, err = wc.Write([]byte(messageBody)); err != nil {
		return err
	}

	err = wc.Close()
	if err != nil {
		return err
	}

	c.Quit()
	return nil
}

func main() {
	userPtr := flag.String("user", "", "Username")
	senderPtr := flag.String("sender", "you@example.com", "Email From")
	rcptPtr := flag.String("rcpt", "you@example.com", "Email Recipient")
	smtpServerPtr := flag.String("smtpServer", "smtp.gmail.com", "SMTP Server")
	smtpPortPtr := flag.String("smtpPort", "587", "SMTP Port")
	smtpPassword := flag.String("smtpPassword", "", "SMTP Password")
	sendPtr := flag.Bool("send", false, "Send Email")
	flag.Parse()

	rotationScriptURL := "https://github.com/605data/aws_scripts/blob/master/aws-iam-rotate-keys.sh"

	result, err := getAccessKeys(*userPtr)
	if err != nil {
		panic("Unable to get keys.")
	}

	for date := range result {
		mail := Mail{}
		mail.senderId = *senderPtr
		mail.toId = *rcptPtr

		diff := checkCreateDate(*result[date].CreateDate)
		// 90 days = 2160 Hours
		// 85 Days = 2040 Hours
		if *sendPtr == true {
			if hours := diff.Hours(); hours >= 2040 {
				if hours >= 2160 {
					mail.subject = "Your IAM Access Keys are at least 90 days old."
					mail.body = "Hello, your IAM Access Keys are at least 90 days old.\n\n"
				} else {
					mail.subject = "Your IAM Access Keys are nearing 90 days old."
					mail.body = "Hello, your IAM Access Keys are nearing 90 days old.\n\n"
				}
				mail.body += "Please rotate your access keys. You can use the script located at\n" +
					rotationScriptURL +
					"\n\nThank you for doing your part to keep our accounts more secure!\n"
				if err := sendSmtpEmail(*smtpServerPtr, *smtpPortPtr, *smtpPassword, mail); err != nil {
					fmt.Println("Error sending email:", err)
				} else {
				  fmt.Println("Success: Sent email to", *rcptPtr)
				}
			} else {
				fmt.Println("No email sent.", *userPtr, "IAM Access Keys were", diff.Hours()/24, "days old.")
			}
		} else {
			fmt.Println("No email sent.", *userPtr, "IAM Access Keys were", diff.Hours()/24, "days old.")
		}
	}
}
