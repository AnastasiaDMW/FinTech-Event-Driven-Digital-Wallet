package notifier

import (
	"fmt"
	"log/slog"
	"net/smtp"

	"github.com/AnastasiaDMW/notification-service/internal/model"
)

type SMTPConfig struct {
	Host     string
	Port     string
	Username string
	Password string
	From     string
}

type Notifier struct {
	logger *slog.Logger
	cfg    SMTPConfig
}

func New(logger *slog.Logger, cfg SMTPConfig) *Notifier {
	return &Notifier{
		logger: logger,
		cfg:    cfg,
	}
}

func (n *Notifier) Send(email string, e model.TransactionEvent) error {
	n.logger.Debug("sending notification", "user_id", e.UserID, "email", email)

	subject := "Уведомление о статусе операции"
	body := n.buildMessage(e)

	msg := fmt.Sprintf(
		"From: %s\r\nTo: %s\r\nSubject: %s\r\n\r\n%s",
		n.cfg.From,
		email,
		subject,
		body,
	)

	addr := fmt.Sprintf("%s:%s", n.cfg.Host, n.cfg.Port)

	auth := smtp.PlainAuth(
		"",
		n.cfg.Username,
		n.cfg.Password,
		n.cfg.Host,
	)

	err := smtp.SendMail(
		addr,
		auth,
		n.cfg.From,
		[]string{email},
		[]byte(msg),
	)
	if err != nil {
		n.logger.Debug("failed to send email", "error", err, "user_id", e.UserID)
	} else {
		n.logger.Debug("email sent successfully", "user_id", e.UserID, "email", email)
	}

	return nil
}

func (n *Notifier) buildMessage(e model.TransactionEvent) string {
	statusText := "успешно выполнена"
	if e.Status != "success" {
		statusText = "отклонена"
	}

	return fmt.Sprintf(
		`Уважаемый клиент,

		Сообщаем, что операция по вашему счёту была %s.

		Детали операции:
		— Идентификатор пользователя: %d
		— Сумма операции: %.2f RUB
		— Статус: %s

		Если вы не совершали данную операцию, пожалуйста, незамедлительно свяжитесь со службой поддержки банка.

		С уважением,
		Служба уведомлений банка
		`,
		statusText,
		e.UserID,
		e.Amount,
		statusText,
	)
}
