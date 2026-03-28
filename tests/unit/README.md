"""
auth-gateway README

This is a simple authentication gateway written in Python.

It provides a basic API for managing users and authentication tokens.

Usage
-----

To run the gateway, execute the following command:

    python gateway.py

API Endpoints
------------

### Users

*   **GET /users**: Retrieve a list of all users.
*   **POST /users**: Create a new user.
*   **GET /users/{id}**: Retrieve a user by their ID.
*   **PUT /users/{id}**: Update a user.
*   **DELETE /users/{id}**: Delete a user.

### Authentication

*   **POST /auth**: Authenticate a user and receive an authentication token.
*   **GET /auth/{token}**: Verify an authentication token.

Configuration
------------

The gateway uses a simple configuration file located at `config.json`.

Example configuration:

    {
        "database": {
            "host": "localhost",
            "port": 5432,
            "username": "username",
            "password": "password"
        },
        "secret_key": "secret"
    }

Requirements
------------

*   Python 3.7+
*   `flask`
*   `flask_sqlalchemy`

"""

from flask import Flask, jsonify, request
from flask_sqlalchemy import SQLAlchemy

app = Flask(__name__)
app.config.from_file('config.json', load=lambda f: json.load(f))
db = SQLAlchemy(app)

from models import User

@app.route('/users', methods=['GET'])
def get_users():
    users = User.query.all()
    return jsonify([user.to_dict() for user in users])

@app.route('/users', methods=['POST'])
def create_user():
    data = request.json
    user = User(**data)
    db.session.add(user)
    db.session.commit()
    return jsonify(user.to_dict())

@app.route('/users/<int:id>', methods=['GET'])
def get_user(id):
    user = User.query.get(id)
    return jsonify(user.to_dict())

@app.route('/users/<int:id>', methods=['PUT'])
def update_user(id):
    user = User.query.get(id)
    data = request.json
    for key, value in data.items():
        setattr(user, key, value)
    db.session.commit()
    return jsonify(user.to_dict())

@app.route('/users/<int:id>', methods=['DELETE'])
def delete_user(id):
    user = User.query.get(id)
    db.session.delete(user)
    db.session.commit()
    return jsonify({'message': 'User deleted'})

@app.route('/auth', methods=['POST'])
def authenticate():
    data = request.json
    user = User.query.filter_by(username=data['username']).first()
    if user and user.check_password(data['password']):
        token = user.generate_token()
        return jsonify({'token': token})
    return jsonify({'message': 'Invalid credentials'}), 401

@app.route('/auth/<token>', methods=['GET'])
def verify_token(token):
    user = User.verify_token(token)
    if user:
        return jsonify({'message': 'Token is valid'})
    return jsonify({'message': 'Token is invalid'}), 401

if __name__ == '__main__':
    app.run(debug=True)