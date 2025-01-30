# Mail API

![Go](https://img.shields.io/badge/Go-1.23-blue)
![Fiber](https://img.shields.io/badge/Fiber-2.0-green)
![GORM](https://img.shields.io/badge/GORM-1.25-orange)

Welcome to the API Mail Service! This API supports native SMTP, Google Cloud Gmail API over OAuth2, and Azure Outlook API over OAuth2. It provides a RESTful way to manage email services with CRUD operations for each service.

## 🚀 Supported Services
### Native SMTP
The API supports sending emails using native SMTP.

### Google Cloud Gmail API
To use the Gmail API, follow the instructions in the [OAUTH_GMAIL.md](docs/OAUTH_GMAIL.md) file.

### Azure Outlook API
To use the Outlook API, follow the instructions in the [OAUTH_AZURE.md](docs/OAUTH_AZURE.md) file.

## Getting Started

### Prerequisites

- Docker
- Docker Compose

### 🛠️ Installation

1. Clone the repository:

    ```bash
    git clone https://github.com/ArnoldPMolenaar/api-mail.git
    cd api-mail
    ```

2. Build and run the Docker containers:

    ```bash
    docker-compose up dev --build
    ```

3. The API will be available at `http://localhost:5002`.

## 🧑‍💻 API Endpoints
### Send a Mail
- `POST /v1/mail/send`: Send an email using the specified service.

### SMTP
- `POST /v1/smtps`: Create a new SMTP configuration.
- `GET /v1/smtps`: Retrieve a list of SMTP configurations.
- `GET /v1/smtps/{id}`: Retrieve a specific SMTP configuration.
- `PUT /v1/smtps/{id}`: Update a specific SMTP configuration.
- `DELETE /v1/smtps/{id}`: Delete a specific SMTP configuration.
- `PUT /v1/smtps/{id}/restore`: Restore a deleted SMTP configuration.

### Gmail
- `POST /v1/gmails`: Create a new Gmail configuration.
- `GET /v1/gmails`: Retrieve a list of Gmail configurations.
- `GET /v1/gmails/{id}`: Retrieve a specific Gmail configuration.
- `PUT /v1/gmails/{id}`: Update a specific Gmail configuration.
- `DELETE /v1/gmails/{id}`: Delete a specific Gmail configuration.
- `PUT /v1/gmails/{id}/restore`: Restore a deleted Gmail configuration.

### Outlook
- `POST /v1/azures`: Create a new Outlook configuration.
- `GET /v1/azures`: Retrieve a list of Outlook configurations.
- `GET /v1/azures/{id}`: Retrieve a specific Outlook configuration.
- `PUT /v1/azures/{id}`: Update a specific Outlook configuration.
- `DELETE /v1/azures/{id}`: Delete a specific Outlook configuration.
- `PUT /v1/azures/{id}/restore`: Restore a deleted Outlook configuration.

## 🤝 Contributing
We welcome contributions! Please fork the repository and submit a pull request.

## 📝 License

This project is licensed under the MIT License.

## 📞 Contact

For any questions or support, please contact [arnold.molenaar@webmi.nl](mailto:arnold.molenaar@webmi.nl).
<hr></hr> Made with ❤️ by Arnold Molenaar