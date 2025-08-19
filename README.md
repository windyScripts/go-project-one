execs:

[
    {
        "first_name": "Dumble",
        "last_name": "Dore",
        "email": "ddoree@example.com",
        "username": "dumbledore",
        "role": "admin",
        "password": "securepassword1"
    },
    {
        "first_name": "Douglas",
        "last_name": "Adams",
        "email": "justhiker@example.com",
        "username": "solongfish",
        "role": "exec",
        "password": "securepassword2"
    }
]

---

Generate cert.pem and key.pem using:

openssl req -x509 -newkey rsa:2048 -keyout key.pem -out cert.pem -days 365 -nodes
