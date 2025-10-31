# Railway Deployment Guide

## Prerequisites
1. Railway account at https://railway.app
2. GitHub repository connected to Railway

## Deployment Steps

### 1. Add Redis Service
In your Railway dashboard:
1. Click "New" → "Database" → "Add Redis"
2. Railway will automatically create a `REDIS_URL` environment variable

### 2. Deploy Your Go Service
1. Click "New" → "GitHub Repo" → Select your repository
2. Railway will auto-detect the Dockerfile and build

### 3. Set Environment Variables
In Railway dashboard → Your service → Variables, add:

```
JWT_SECRET=8nlq1t/cxGQx8ruOBUhIB8zUmzHeiP6KY/jv4BB21KY=
PORT=8080
```

**Important:** Railway will automatically provide:
- `PORT` (dynamic, overrides your setting)
- `REDIS_URL` (from the Redis service)

### 4. Connect Redis Service
In Railway dashboard:
1. Go to your Go service
2. Click "Settings" → "Service Variables" 
3. Click "New Variable" → "Reference" → Select Redis service → `REDIS_URL`

### 5. Health Check
Railway will ping your `/` endpoint to verify the service is healthy.

## Common Issues

### Redis Connection Errors
- Ensure Redis service is running
- Verify `REDIS_URL` environment variable is set
- Check service connectivity in Railway dashboard

### Port Issues
- Don't hardcode ports - always use `PORT` environment variable
- Railway assigns ports dynamically

### Build Failures
- Ensure all dependencies are in `go.mod`
- Check Dockerfile paths are correct
- Verify Railway has access to your repository

## Local Development
For local development, continue using Docker Compose:
```bash
docker-compose up -d
```

For Railway deployment, only the Dockerfile is used.