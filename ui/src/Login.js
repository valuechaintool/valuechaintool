import React from "react";
import {
    Container,
    Row,
    Col,
    Button,
    Form,
} from 'react-bootstrap';

class Login extends React.Component {
    constructor() {
        super();
        this.state = {
            username: '',
            password: '',
        }
    }

    onChange = (e) => {
        this.setState({ [e.target.name]: e.target.value });
    }

    onSubmit = (e) => {
        e.preventDefault();
        const { username, password } = this.state;
        fetch('http://localhost:8080/api/v1/login', {
            method: 'POST',
            headers: {
                "Content-Type": "application/json"
            },
            body: JSON.stringify({ username, password }),
        })
            .then(
                function(response) {
                    if (response.status !== 200) {
                        console.log('Looks like there was a problem. Status Code: ' +
                            response.status);
                        return;
                    }
                    response.json().then(function(data) {
                        console.log(data);
                        localStorage.setItem('session', data)
                    });
                }
            )
            .catch(function(err) {
                console.log('Fetch Error :-S', err);
            });
    }

    render() {
        const { username, password } = this.state;
        return (
            <Container>
                <Row>
                    <Col sm={4}></Col>
                    <Col sm={4}>
                        <Form onSubmit={this.onSubmit}>
                            <Form.Group controlId="username">
                                <Form.Label>Username</Form.Label>
                                <Form.Control name="username" type="text" placeholder="Username" value={username} onChange={this.onChange} />
                            </Form.Group>
                            <Form.Group controlId="password">
                                <Form.Label>Password</Form.Label>
                                <Form.Control name="password" type="password" placeholder="Password" value={password} onChange={this.onChange}/>
                            </Form.Group>
                            <Button variant="primary" type="submit">Submit</Button>
                        </Form>
                    </Col>
                </Row>
            </Container>
        );
    }
}

export default Login;
