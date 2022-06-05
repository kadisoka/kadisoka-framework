package ses

import (
	"github.com/alloyzeus/go-azfl/azfl/errors"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	awscreds "github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ses"
	"github.com/rez-go/stev"

	"github.com/kadisoka/kadisoka-framework/iam/pkg/iamserver/eav10n"
	"github.com/kadisoka/kadisoka-framework/volib/pkg/email"
)

type EmailDeliveryServiceConfig struct {
	Region          string `env:"REGION,required"`
	AccessKeyID     string `env:"ACCESS_KEY_ID"`
	SecretAccessKey string `env:"SECRET_ACCESS_KEY"`
}

func EmailDeliveryServiceConfigSkeleton() EmailDeliveryServiceConfig {
	return EmailDeliveryServiceConfig{}
}

func (EmailDeliveryServiceConfig) SelfDocsDescriptor() stev.SelfDocsDescriptor {
	return stev.SelfDocsDescriptor{
		ShortDesc: "Use Amazon SES to deliver the emails",
	}
}

type EmailDeliveryService struct {
	sesClient *ses.SES
}

var _ eav10n.EmailDeliveryService = &EmailDeliveryService{}

func NewEmailDeliveryService(config interface{}) eav10n.EmailDeliveryService {
	if config == nil {
		panic(errors.New("configuration required"))
	}
	conf, ok := config.(*EmailDeliveryServiceConfig)
	if !ok {
		panic(errors.New("configuration of invalid type"))
	}

	var creds *awscreds.Credentials
	if conf.AccessKeyID != "" {
		creds = awscreds.NewStaticCredentials(
			conf.AccessKeyID,
			conf.SecretAccessKey,
			"",
		)
	}
	sess, err := session.NewSession(&aws.Config{
		Region:      aws.String(conf.Region),
		Credentials: creds,
	})
	if err != nil {
		panic(err)
	}
	sesClient := ses.New(sess)

	return &EmailDeliveryService{sesClient: sesClient}

}

func (emailDeliverySvc *EmailDeliveryService) SendHTMLMessage(
	recipient email.Address,
	subjectText string,
	htmlContent string,
	opts eav10n.EmailDeliveryOptions,
) error {
	var err error

	// Note that SES supports both text and HTML body. For better
	// support, we might want to utilizes both.
	input := &ses.SendEmailInput{
		Destination: &ses.Destination{
			ToAddresses: []*string{
				aws.String(recipient.String()),
			},
		},
		Message: &ses.Message{
			Body: &ses.Body{
				Html: &ses.Content{
					Charset: aws.String(opts.MessageCharset),
					Data:    aws.String(htmlContent),
				},
			},
			Subject: &ses.Content{
				Charset: aws.String(opts.MessageCharset),
				Data:    aws.String(subjectText),
			},
		},
		Source: aws.String(opts.SenderAddress),
	}

	_, err = emailDeliverySvc.sesClient.SendEmail(input)

	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			//TODO: translate errors
			switch aerr.Code() {
			case ses.ErrCodeMessageRejected:
				return errors.Wrap("SendEmail", aerr)
			case ses.ErrCodeMailFromDomainNotVerifiedException:
				return errors.Wrap("SendEmail", aerr)
			case ses.ErrCodeConfigurationSetDoesNotExistException:
				return errors.Wrap("SendEmail", aerr)
			default:
				return errors.Wrap("SendEmail", aerr)
			}
		}
		return err
	}

	return nil
}
