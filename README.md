# Tubely
A file server application built with Go, AWS S3, and CloudFront. The project handles uploading, storing, streaming, and securely serving files (including video) through a CDN-backed architecture.

## What I Learned

- How to handle large file and multipart uploads in Go
- The difference between filesystem storage and object storage
- How to integrate the AWS SDK for Go with S3
- How HTTP caching works (cache headers, cache-busting)
- How to stream video using HTTP range requests
- How to secure AWS resources using IAM roles, policies, and least-privilege access
- How to generate pre-signed URLs for time-limited, secure file access
- How CDNs like CloudFront reduce latency and origin load
- Key resiliency concepts: availability, reliability, and durability

## Note

This repo is forked from [bootdotdev/learn-file-storage-s3-golang-starter](https://github.com/bootdotdev/learn-file-storage-s3-golang-starter) as a part of learning purposes on the "Learn File Servers and CDNs with S3 and CloudFront" [course](https://www.boot.dev/courses/learn-file-servers-s3-cloudfront-golang) on [boot.dev](https://www.boot.dev)