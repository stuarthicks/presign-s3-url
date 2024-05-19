package main

import (
	"flag"
	"fmt"
	"log"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

func main() {
	var s3URL, region string
	var hours int
	flag.StringVar(&s3URL, "url", "", "The S3 URL to sign")
	flag.StringVar(&region, "region", "us-east-1", "The AWS region in which the object is stored")
	flag.IntVar(&hours, "hours", 168, "Number of hours the URL should be valid for")
	flag.Parse()

	if s3URL == "" || region == "" {
		fmt.Fprint(os.Stderr, "Invalid arguments")
		flag.PrintDefaults()
		os.Exit(1)
	}

	u, err := url.Parse(s3URL)
	if err != nil {
		log.Fatal("failed to parse S3 URL", err.Error())
	}

	cfg := aws.NewConfig()
	cfg.WithRegion(region)
	if strings.Contains(u.Hostname(), ".") {
		cfg.WithS3ForcePathStyle(true)
	}

	sess := session.Must(session.NewSession(cfg))

	client := s3.New(sess)

	req, _ := client.GetObjectRequest(&s3.GetObjectInput{
		Bucket: aws.String(u.Hostname()),
		Key:    aws.String(u.Path),
	})

	signedURL, err := req.Presign(time.Duration(int64(hours) * int64(time.Minute) * int64(24)))
	if err != nil {
		log.Fatal("failed to sign URL", err.Error())
	}

	fmt.Println(signedURL)
}
