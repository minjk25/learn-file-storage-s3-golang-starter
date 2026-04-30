# Tubely
A file server application built with Go, AWS S3, and CloudFront. The project handles uploading, storing, streaming, and securely serving files (including video) through a CDN-backed architecture.

## Overview
This project demonstrates a cloud-based data storage workflow using AWS S3.
- Data ingestion via API
- Cloud object storage (S3)
- Scalable data storage architecture

## Key Concepts & Learnings

| Feature / Concept | What I Learned |
|------------------|---------------|
| Large File Handling | How to handle large files and multipart uploads efficiently in Go |
| Storage Architecture | The difference between filesystem storage and object storage (S3) |
| AWS Integration | How to integrate the AWS SDK for Go with S3 services |
| HTTP Caching | How cache headers and cache-busting strategies improve performance |
| Video Streaming | How HTTP range requests enable efficient video streaming |
| Security (IAM) | How to secure AWS resources using IAM roles, policies, and least-privilege access |
| Pre-signed URLs | How to generate time-limited URLs for secure file access |
| CDN (CloudFront) | How CDNs reduce latency and decrease load on origin servers |
| Resiliency | Key concepts: availability, reliability, and durability in distributed systems |

## Requirement
- [Go](https://golang.org/doc/install)
- [Goose](https://github.com/pressly/goose)
- [FFMPEG](https://ffmpeg.org/download.html)
- [SQLite 3](https://www.sqlite.org/download.html)
- [AWS CLI](https://docs.aws.amazon.com/cli/latest/userguide/getting-started-install.html)
- AWS acoount with S3 and CloudFront setup

## Installation
1. Clone the repository:
```bash
git clone https://github.com/minjk25/learn-file-storage-s3-golang-starter.git
cd learn-file-storage-s3-golang-starter
go mod tidy
```

2. Create a `.env` file at the root of the directory, here is example:
```bash
DB_PATH="./tubely.db?_foreign_keys=on"
JWT_SECRET="JKFNDKAJSDKFASFNJWIROIOTNKNFDSKNFD" # this is just random string, you can change it
PLATFORM="dev"
FILEPATH_ROOT="./app"
ASSETS_ROOT="./assets"
S3_BUCKET="your-set-up-from-s3" # set up from S3
S3_REGION="your-set-up-from-s3"  # set up from S3
S3_CF_DISTRO="your-set-up-from-s3"  # set up from S3
PORT="8091"
```
3. Run the server at the root of directory:
```bash
go run .
```

## Note

This repo is forked from [bootdotdev/learn-file-storage-s3-golang-starter](https://github.com/bootdotdev/learn-file-storage-s3-golang-starter) as a part of learning purposes on the "Learn File Servers and CDNs with S3 and CloudFront" [course](https://www.boot.dev/courses/learn-file-servers-s3-cloudfront-golang) on [boot.dev](https://www.boot.dev)
