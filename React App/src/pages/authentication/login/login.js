import React, {useState} from 'react';
import { withRouter } from 'react-router-dom';

import { loginHTTPRequest } from "./../../../services/api-service";
import { setItemInLS } from "./../../../services/storage-service";

import './login.css'

function Login(props) {

  const [loginErrorMessage, setErrorMessage] = useState(null);
  const [username, updateUsername] = useState(null);
  const [password, updatePassword] = useState(null);


  const handleUsernameChange = (event) => {
    updateUsername(event.target.value)
  }

  const handlePasswordChange = (event) => {
    updatePassword(event.target.value)
  }

  const loginUser = async () => {
    props.displayPageLoader(true);
    const userDetails = await loginHTTPRequest(username, password);
    props.displayPageLoader(false);

    if (userDetails.code === 200) {
      setItemInLS('userDetails', userDetails.response)
      props.history.push(`/home`)
    } else {
      setErrorMessage(userDetails.message);
    }
  };

  return (
    <div className="app__login-container">
      <div className="app__form-row">
        <label>Username:</label>
        <input type="email" className="email" onChange={handleUsernameChange} />
      </div>
      <div className="app__form-row">
        <label>Password:</label>
        <input type="password" className="password" onChange={handlePasswordChange} />
      </div>
      <div className="app__form-row">
        <span className="error-message">{loginErrorMessage? loginErrorMessage : ''}</span>
      </div>
      <div className="app__form-row">
        <button onClick={loginUser}>Login</button>
      </div>
    </div>
  );
}

export default withRouter(Login);