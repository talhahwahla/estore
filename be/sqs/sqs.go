package sqs

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/aws/aws-sdk-go-v2/service/sqs/types"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/joho/godotenv"
)

type SQSService struct {
	SqsClient *sqs.Client
	QueueURL  string
}

func NewSQSService(queueName string) *SQSService {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file:", err)
	}

	accessKeyID := os.Getenv("AWS_ACCESS_KEY_ID")
	secretAccessKey := os.Getenv("AWS_SECRET_ACCESS_KEY")
	region := os.Getenv("AWS_REGION")

	if accessKeyID == "" || secretAccessKey == "" || region == "" {
		log.Fatal("AWS credentials or region not provided in environment variables.")
		return nil
	}

	creds := credentials.NewStaticCredentialsProvider(accessKeyID, secretAccessKey, "")
	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithCredentialsProvider(creds),
		config.WithRegion(region),
	)
	if err != nil {
		log.Fatal("Error loading configuration:", err)
		return nil
	}

	sqsClient := sqs.NewFromConfig(cfg)

	res, err := sqsClient.GetQueueUrl(
		context.TODO(),
		&sqs.GetQueueUrlInput{
			QueueName: &queueName,
		})
	if err != nil {
		log.Fatal("Error getting queue:", err)
	}

	return &SQSService{
		SqsClient: sqsClient,
		QueueURL:  *res.QueueUrl,
	}
}

func (s *SQSService) SendMessage(messageBody string) error {
	input := &sqs.SendMessageInput{
		MessageBody: aws.String(messageBody),
		QueueUrl:    aws.String(s.QueueURL),
	}

	_, err := s.SqsClient.SendMessage(context.Background(), input)
	if err != nil {
		return fmt.Errorf("failed to send message to SQS: %v", err)
	}

	return nil
}

func (s *SQSService) ReceiveMessage() ([]types.Message, error) {
	input := &sqs.ReceiveMessageInput{
		QueueUrl:            aws.String(s.QueueURL),
		MaxNumberOfMessages: *aws.Int32(10), // Adjust as needed
	}

	result, err := s.SqsClient.ReceiveMessage(context.Background(), input)
	if err != nil {
		return nil, fmt.Errorf("failed to receive messages from SQS: %v", err)
	}

	return result.Messages, nil
}

func (s *SQSService) DeleteMessage(receiptHandle string) error {
	input := &sqs.DeleteMessageInput{
		QueueUrl:      aws.String(s.QueueURL),
		ReceiptHandle: aws.String(receiptHandle),
	}

	_, err := s.SqsClient.DeleteMessage(context.Background(), input)
	if err != nil {
		return fmt.Errorf("failed to delete message from SQS: %v", err)
	}

	return nil
}
