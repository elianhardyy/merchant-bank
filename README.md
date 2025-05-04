# Merchant-Bank API

A RESTful API service built with Go that facilitates interactions between merchants and banks, handling customer authentication, payment processing, and transaction history.

## Features

- **Authentication**: Register, login, and logout functionality for customers
- **Payment Processing**: Secure transfer between registered customers
- **Transaction History**: Complete logging of all transactions
- **Role-Based Access Control**: Different permissions for customers and merchants
- **JWT Authentication**: Secure API endpoints

## Technology Stack

- Go 1.23.2
- Gorilla Mux (Router)
- JWT for Authentication
- JSON for data storage

## API Endpoints

| Method | Endpoint          | Description             | Access             |
| ------ | ----------------- | ----------------------- | ------------------ |
| POST   | /auth/register    | Register a new user     | Public             |
| POST   | /auth/login       | Login a user            | Public             |
| POST   | /auth/logout      | Logout a user           | Customer, Merchant |
| POST   | /trx/create       | Process a payment       | Customer           |
| GET    | /trx/history/{id} | Get transaction history | Customer, Merchant |
| GET    | /user/users       | Get list of users       | Merchant           |

## Prerequisites

- Go 1.23.2 or higher
- Git

## Setup and Installation

1. Clone the repository:

   ```bash
   git clone https://github.com/yourusername/merchant-bank-api.git
   cd merchant-bank-api
   ```

2. Install dependencies:

   ```bash
   go mod tidy
   ```

3. Create a .env file in the root directory with the following content:

   ```
   JWT_SECRET=your_jwt_secret_key
   ```

4. Run the application:

   ```bash
   go run main.go
   ```

5. For development with hot reload, you can use Air:

   ```bash
   # Install Air first if you haven't
   go install github.com/cosmtrek/air@latest

   # if Air on above didn't work install another package
   go install github.com/air-verse/air@latest
   
   # Run with Air
   air
   ```

## Project Structure Explanation

- **Controllers**: Handle HTTP requests and responses
- **Middlewares**: Implement authentication and request/response processing
- **Models**: Define data structures
- **Repositories**: Manage data access (JSON files in this implementation)
- **Services**: Implement business logic
- **Security**: Handle JWT token generation and validation

## Testing

Run the tests with:

```bash
go test ./tests... -v
```

The project implements Test-Driven Development (TDD) principles with unit tests for each component:

- Repository tests
- Service tests
- Controller tests
- Integration tests

## Security Considerations

1. **JWT Authentication**: All protected routes are secured with JWT tokens
2. **Password Hashing**: User passwords are hashed before storage
3. **Role-Based Access Control**: Routes are protected based on user roles
4. **Input Validation**: All user inputs are validated
5. **Error Handling**: Proper error handling with appropriate HTTP status codes
6. **Environment Variables**: Sensitive information stored in environment variables

## API Usage Examples

### Register a new user

```bash
curl -X POST http://localhost:8080/auth/register \
  -H "Content-Type: application/json" \
  -d '{"username":"newuser","password":"password123","role":"customer"}'
```

### Login

```bash
curl -X POST http://localhost:8080/auth/login \
  -H "Content-Type: application/json" \
  -d '{"username":"customer1","password":"password123"}'
```

### Make a payment

```bash
curl -X POST http://localhost:8080/trx/create \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer your_token_here" \
  -d '{"customer_id":"1","merhcant_id":"2","amount":100}'
```

### View transaction history

```bash
curl -X GET http://localhost:8080/trx/history/1 \
  -H "Authorization: Bearer your_token_here"
```

## Deployment

For deployment to a production environment:

1. Build the application:

   ```bash
   go build -o app
   ```

2. Set up environment variables for production.

3. Run the application:

   ```bash
   ./app
   ```

4. For containerization, a Dockerfile is provided:
   ```bash
   docker build -t merchant-bank-api .
   docker run -p 8080:8080 merchant-bank-api
   ```

## Future Improvements

- Add database support (PostgreSQL, MongoDB)
- Implement refresh token mechanism
- Add rate limiting
- Implement logging to external services
- Add metrics and monitoring
- Implement OpenAPI/Swagger documentation

## License

This project is licensed under the MIT License - see the LICENSE file for details.

## 📌 API Documentation on Postman

For the full API documentation, please visit the following link:

👉 [Postman Documentation](https://documenter.getpostman.com/view/18886846/2sB2cPjkEe)
