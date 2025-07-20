# Setup Guide

This guide will help you set up the backend Go application for the job application system.

## Prerequisites

- Go 1.19 or higher
- Google Cloud Platform account with Cloud Storage enabled
- PostgreSQL database

## Configuration Files

### 1. Environment Configuration (config.env)

Create a `config.env` file in the **cmd directory** of the project with the following variables:

```env
# Server Configuration
PORT=8080

# Database Configuration
DB_HOST=localhost
DB_PORT=5432
DB_USER=your_db_user
DB_PASSWORD=your_db_password
DB_NAME=your_database_name
DB_SSLMODE=disable

# JWT Configuration
JWT_SECRET=your-super-secret-jwt-key-here-make-it-long-and-random
JWT_EXPIRY=24h

# Google Cloud Configuration
GOOGLE_PROJECT_ID=your-gcp-project-id
GOOGLE_CREDENTIALS_FILE=credentials.json
GCS_BUCKET_NAME=your-gcs-bucket-name

# Email Configuration (if using email features)
SMTP_HOST=smtp.gmail.com
SMTP_PORT=587
SMTP_USERNAME=your-email@gmail.com
SMTP_PASSWORD=your-app-password
```

### 2. Google Cloud Credentials (credentials.json)

Place your Google Cloud service account credentials file named `credentials.json` in the **root directory** of the project.

#### How to get credentials.json:

1. Go to [Google Cloud Console](https://console.cloud.google.com/)
2. Navigate to **IAM & Admin** > **Service Accounts**
3. Create a new service account or select an existing one
4. Go to the **Keys** tab
5. Click **Add Key** > **Create new key**
6. Choose **JSON** format
7. Download the file and rename it to `credentials.json`
8. Place it in the project root directory

**Important**: Never commit `credentials.json` to version control. Add it to your `.gitignore` file.

## File Structure

Your project root should look like this:

```
backend-go/
├── credentials.json        # Google Cloud credentials
├── go.mod
├── go.sum
├── cmd/
│   ├── main.go
│   └── config.env         # Environment variables
├── internal/
│   ├── config/
│   ├── db/
│   ├── gcs/
│   ├── handler/
│   ├── middleware/
│   ├── models/
│   ├── routes/
│   ├── services/
│   └── utils/
└── SETUP.md
```

## Environment Variables Explained

### Server Configuration

- `PORT`: The port on which the server will run (default: 8080)

### Database Configuration

- `DB_HOST`: PostgreSQL server hostname
- `DB_PORT`: PostgreSQL server port (default: 5432)
- `DB_USER`: Database username
- `DB_PASSWORD`: Database password
- `DB_NAME`: Database name
- `DB_SSLMODE`: SSL mode for database connection

### JWT Configuration

- `JWT_SECRET`: Secret key for signing JWT tokens (make it long and random)
- `JWT_EXPIRY`: Token expiration time (e.g., "24h", "7d")

### Google Cloud Configuration

- `GOOGLE_PROJECT_ID`: Your Google Cloud Project ID
- `GOOGLE_CREDENTIALS_FILE`: Path to credentials file (should be "credentials.json")
- `GCS_BUCKET_NAME`: Google Cloud Storage bucket name for storing files

### Email Configuration

- `SMTP_HOST`: SMTP server hostname
- `SMTP_PORT`: SMTP server port
- `SMTP_USERNAME`: Email username
- `SMTP_PASSWORD`: Email password or app password

## Setup Steps

1. **Clone the repository**

   ```bash
   git clone <repository-url>
   cd backend-go
   ```

2. **Install dependencies**

   ```bash
   go mod download
   ```

3. **Create configuration files**

   - Create `config.env` file in the cmd directory with your configuration
   - Place `credentials.json` in the root directory

4. **Set up Google Cloud Storage**

   - Create a bucket in Google Cloud Storage
   - Ensure your service account has the necessary permissions:
     - `Storage Object Admin` for the bucket
     - `Storage Object Viewer` for reading files

5. **Set up PostgreSQL database**

   - Create a database
   - Run any required migrations (if applicable)

6. **Run the application**
   ```bash
   go run cmd/main.go
   ```

## Security Notes

- **Never commit sensitive files**: Add the following to your `.gitignore`:

  ```
  cmd/config.env
  credentials.json
  *.key
  *.pem
  ```

- **Use strong secrets**: Generate a strong JWT secret:

  ```bash
  openssl rand -base64 32
  ```

- **Restrict file permissions**: Ensure only the application can read the credentials file:
  ```bash
  chmod 600 credentials.json
  ```

## Testing the Setup

Once the application is running, you can test the endpoints:

1. **Health check**: `GET http://localhost:8080/api/v1/health`
2. **Signup**: `POST http://localhost:8080/api/v1/signup`
3. **Login**: `POST http://localhost:8080/api/v1/login`

## Troubleshooting

### Common Issues

1. **"credentials.json not found"**

   - Ensure the file is in the root directory
   - Check the file permissions

2. **"Database connection failed"**

   - Verify database credentials in `config.env`
   - Ensure PostgreSQL is running
   - Check network connectivity

3. **"JWT secret not configured"**

   - Add JWT_SECRET to your `config.env` file
   - Ensure it's a strong, random string

4. **"GCS bucket not found"**
   - Verify the bucket name in `config.env`
   - Ensure the service account has access to the bucket

### Getting Help

If you encounter issues:

1. Check the application logs for error messages
2. Verify all configuration files are in the correct locations
3. Ensure all required services (PostgreSQL, Google Cloud) are accessible
