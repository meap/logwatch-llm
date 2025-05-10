package presenter

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"io"
	"mime/quotedprintable"
	"os/exec"

	"github.com/sirupsen/logrus"
)

type EmailAttachment struct {
	Name        string
	Content     []byte
	ContentType string // optional, defaults to application/octet-stream
}

type EmailMessage struct {
	To          string
	Subject     string
	PlainBody   string
	Attachments []EmailAttachment // support multiple attachments
}

// ComposePlainEmailWithAttachments builds the raw email message as bytes.
func ComposePlainEmailWithAttachments(msgData EmailMessage) ([]byte, error) {
	if msgData.To == "" {
		return nil, fmt.Errorf("recipient email address is required")
	}

	var msg bytes.Buffer

	if len(msgData.Attachments) == 0 {
		// No attachments, just send plain text email
		msg.WriteString(fmt.Sprintf("To: %s\r\n", msgData.To))
		msg.WriteString(fmt.Sprintf("Subject: %s\r\n", msgData.Subject))
		msg.WriteString("MIME-Version: 1.0\r\n")
		msg.WriteString("Content-Type: text/plain; charset=UTF-8\r\n")
		msg.WriteString("Content-Transfer-Encoding: quoted-printable\r\n\r\n")
		qp := quotedprintable.NewWriter(&msg)
		qp.Write([]byte(msgData.PlainBody))
		qp.Close()
	} else {
		boundary := "BOUNDARY-1234567890"
		// Write headers
		msg.WriteString(fmt.Sprintf("To: %s\r\n", msgData.To))
		msg.WriteString(fmt.Sprintf("Subject: %s\r\n", msgData.Subject))
		msg.WriteString("MIME-Version: 1.0\r\n")
		msg.WriteString(fmt.Sprintf("Content-Type: multipart/mixed; boundary=%q\r\n", boundary))
		msg.WriteString("\r\n")

		// Plain text part
		msg.WriteString(fmt.Sprintf("--%s\r\n", boundary))
		msg.WriteString("Content-Type: text/plain; charset=UTF-8\r\n")
		msg.WriteString("Content-Transfer-Encoding: quoted-printable\r\n\r\n")
		qp := quotedprintable.NewWriter(&msg)
		qp.Write([]byte(msgData.PlainBody))
		qp.Close()
		msg.WriteString("\r\n")

		// Attachments
		for _, attachment := range msgData.Attachments {
			ctype := attachment.ContentType
			if ctype == "" {
				ctype = "application/octet-stream"
			}
			msg.WriteString(fmt.Sprintf("--%s\r\n", boundary))
			msg.WriteString(fmt.Sprintf("Content-Type: %s; name=\"%s\"\r\n", ctype, attachment.Name))
			msg.WriteString("Content-Transfer-Encoding: base64\r\n")
			msg.WriteString(fmt.Sprintf("Content-Disposition: attachment; filename=\"%s\"\r\n\r\n", attachment.Name))
			b64 := base64.NewEncoder(base64.StdEncoding, &msg)
			b64.Write(attachment.Content)
			b64.Close()
			msg.WriteString("\r\n")
		}

		// End boundary
		msg.WriteString(fmt.Sprintf("--%s--\r\n", boundary))
	}

	return msg.Bytes(), nil
}

// SendRawEmail pipes the raw email message to sendmail.
func SendRawEmail(rawMsg []byte) error {
	cmd := exec.Command("sendmail", "-t")
	stdin, err := cmd.StdinPipe()
	if err != nil {
		logrus.Errorf("failed to get sendmail stdin pipe: %v", err)
		return err
	}
	// Start the process before writing to stdin to avoid deadlock
	if err := cmd.Start(); err != nil {
		stdin.Close()
		logrus.Errorf("failed to start sendmail: %v", err)
		return err
	}
	// Write the message
	_, writeErr := io.Copy(stdin, bytes.NewReader(rawMsg))
	stdin.Close() // Close after writing to signal EOF
	if writeErr != nil {
		logrus.Errorf("failed to write message to sendmail stdin: %v", writeErr)
		cmd.Wait() // Clean up process
		return writeErr
	}
	if err := cmd.Wait(); err != nil {
		logrus.Errorf("sendmail process failed: %v", err)
		return err
	}

	return nil
}
