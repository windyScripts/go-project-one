## Admin dashboard

Generate cert.pem and key.pem using:

openssl req -x509 -newkey rsa:2048 -keyout key.pem -out cert.pem -days 365 -nodes

For generating from openssl.cnf file:

openssl req -x509 -newkey rsa:2048 -nodes \
  -keyout key.pem \
  -out cert.pem \
  -days 365 \
  -config openssl.cnf

  `mailhog` for running mailhog server in terminal after installing with homebrew.

run project with go run ./cmd/api

MIddleware implemented within this project:

Compression

CORS

HPP

JWT

Rate Limiter

Response Time

Data Sanitization

Security Headers

Levels of Authority within Institute:
Student, Teacher, Admin, Manager, Director

