import React from 'react';
import CreditCardInformation from './CreditCards';
import cookie from 'js-cookie';
import { Modal, ModalHeader, ModalBody } from 'reactstrap';

function submitRequest(path, requestBody, handleSignedIn,handleError) {
    fetch(path, {
        method: 'POST',
        headers: {
            'Accept': 'application/json',
            'Content-Type': 'application/json',
        },
        body: JSON.stringify(requestBody)
    }).then(response => response.json())
      .then(json => {
            console.log("Response received...")
            if (json.error === undefined || !json.error) {
                //save cookie if not error
                console.log("Sign in Success...");
                cookie.set("user", json);
                handleSignedIn(json);
            } else {
                handleError(json.error);
            }
        })
        .catch(error=>console.log(error));
}


// 로그인 폼 컴포넌트
class SingInForm extends React.Component {
    constructor(props) {
        super(props);
        // 사용자가 데이터를 입력하면 호출되는 함수
        this.handleChange = this.handleChange.bind(this);
        // 폼을 제출하면 호출되는 함수
        this.handleSubmit = this.handleSubmit.bind(this);
        // 로그인 실패시 errormessage 필드에 메시지를 저장한다.        
        this.handleError = this.handleError.bind(this);
        this.state = {
            errormessage: ''
        }
    }

    // 폼의 값을 state 객체에 저장하는 방식. 리액트에서 권장하는 폼 제어 방식
    handleChange(event) {
        const name = event.target.name;
        const value = event.target.value;
        this.setState({
            [name]: value
        });
    }

    handleError(error){
        this.setState({
            errormessage: error
        });
    }

    handleSubmit(event) {
        //'users/signin'
        event.preventDefault();
        submitRequest('users/signin', this.state, this.props.handleSignedIn,this.handleError);
    }


    render() {
        // 에러 메시지
        let message = null;
        // state에 에러 메시지가 있다면 출력
        if (this.state.errormessage.length !== 0) {
            message = <h5 className="mb-4 text-danger">{this.state.errormessage}</h5>;

        }
        return (
            <div>
                {message}
                <form onSubmit={this.handleSubmit}>
                    <h5 className="mb-4">Basic Info</h5>
                    <div className="form-group">
                        <label htmlFor="email">Email:</label>
                        <input name="email" type="email" className="form-control" id="email" onChange={this.handleChange} />
                    </div>
                    <div className="form-group">
                        <label htmlFor="pass">Password:</label>
                        <input name="password" type="password" className="form-control" id="pass" onChange={this.handleChange} />
                    </div>
                    <div className="form-row text-center">
                        <div className="col-12 mt-2">
                            <button type="submit" className="btn btn-success btn-large" >Sign In</button>
                        </div>
                        <div className="col-12 mt-2">
                            <button className="btn btn-link text-info" onClick={() => this.props.handleNewUser()}> New User? Register</button>
                        </div>
                    </div>
                </form>
            </div>
        );
    }

}

// 가입 폼 컴포넌트
class RegistrationForm extends React.Component {
    constructor(props) {
        super(props);
        this.handleSubmit = this.handleSubmit.bind(this);
        // 로그인 실패시 errormessage 필드에 메시지를 저장
        this.state = {
            errormessage: ''
        }
        this.handleChange = this.handleChange.bind(this);
        this.handleSubmit = this.handleSubmit.bind(this);
        this.handleError = this.handleError.bind(this);
    }

    handleChange(event) {
        event.preventDefault();
        const name = event.target.name;
        const value = event.target.value;
        this.setState({
            [name]: value
        });
    }
    
    handleError(error){
        this.setState({
            errormessage: error
        });
    }

    handleSubmit(event) {
        event.preventDefault();
        const userInfo = this.state;
        const firstlastname = userInfo.username.split(" ");
        if (userInfo.pass1 !== userInfo.pass2) {
            alert("PASSWORDS DO NOT MATCH");
            return;
        }
        const requestBody = {
            firstname: firstlastname[0],
            lastname: firstlastname[1],
            email: userInfo.email,
            password: userInfo.pass1
        };
        submitRequest('users', requestBody, this.props.handleSignedIn,this.handleError);
        
        console.log("Registration form: " + requestBody);
    }

    render() {
        let message = null;
        if (this.state.errormessage.length !== 0) {
            message = <h5 className="mb-4 text-danger">{this.state.errormessage}</h5>;
        }
        return (
            <div>
                {message}
                <form onSubmit={this.handleSubmit}>
                    <h5 className="mb-4">Registration</h5>
                    <div className="form-group">
                        <label htmlFor="username">User Name:</label>
                        <input id="username" name='username' className="form-control" placeholder='John Doe' type='text' onChange={this.handleChange} />
                    </div>

                    <div className="form-group">
                        <label htmlFor="email">Email:</label>
                        <input type="email" name='email' className="form-control" id="email" onChange={this.handleChange} />
                    </div>
                    <div className="form-group">
                        <label htmlFor="pass">Password:</label>
                        <input type="password" name='pass1' className="form-control" id="pass1" onChange={this.handleChange} />
                    </div>
                    <div className="form-group">
                        <label htmlFor="pass">Confirm password:</label>
                        <input type="password" name='pass2' className="form-control" id="pass2" onChange={this.handleChange} />
                    </div>
                    <div className="form-row text-center">
                        <div className="col-12 mt-2">
                            <button type="submit" className="btn btn-success btn-large">Register</button>
                        </div>
                    </div>
                </form>
            </div>
        );
    }
}

// 로그인 폼을 포함하는 부모 모달 윈도우
// 초기값은 로그인 페이지이고 사용자가 New User? 링크를 클릭하면 가입 페이지로 변경한다.
export class SignInModalWindow extends React.Component {
    constructor(props) {
        super(props);
        this.state = {
            showRegistrationForm: false
        };
        this.handleNewUser = this.handleNewUser.bind(this);
        this.handleModalClose = this.handleModalClose.bind(this);
    }

    handleNewUser() {
        this.setState({
          showRegistrationForm: true
        });
    }

    handleModalClose(){
        this.setState({
            showRegistrationForm: false
        });
    }
   

    render() {
        // state 객체의 값에 따라 SignInForm이나 RegistrationForm 컴포넌트를 모달 윈도우에 추가한다.
        let modalBody = <SingInForm handleNewUser={this.handleNewUser} handleSignedIn={this.props.handleSignedIn} />
        if (this.state.showRegistrationForm === true) {
            modalBody = <RegistrationForm handleSignedIn={this.props.handleSignedIn} />
        }
        return (
            <Modal id="register" tabIndex="-1" role="dialog" isOpen={this.props.showModal} toggle={this.props.toggle} onClosed={this.handleModalClose}>
                <div role="document">
                    <ModalHeader toggle={this.props.toggle} className="bg-success text-white">
                        Sign in
                    </ModalHeader>
                    <ModalBody>
                        {modalBody}
                    </ModalBody>
                </div>
            </Modal>
        );
    }
}


// 다른 파일에서 이 클래스를 사용하기 때문에 export 키워드를 사용
export function BuyModalWindow(props) {
    return (
        // Card 컴포넌트의 Buy버튼을 클릭하면 #buy ID에 해당하는 모달 윈도우를 출력한다.
        <Modal id="buy" tabIndex="-1" role="dialog" isOpen={props.showModal} toggle={props.toggle}>
            <div role="document">
                    <ModalHeader toggle={props.toggle} className="bg-success text-white">
                        Buy Item
                    </ModalHeader>
                    {/* 신용카드 결제 폼 */}
                    <ModalBody>
                        <CreditCardInformation user={props.user} seperator={false} show={true} productid={props.productid} price={props.price} operation="Charge" toggle={props.toggle} />
                    </ModalBody>
                </div>
                      
        </Modal>
    );
} 