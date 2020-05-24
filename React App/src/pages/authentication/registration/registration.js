import React, {useState} from 'react';
import { withRouter } from 'react-router-dom';

import { isUsernameAvailableHTTPRequest, registerHTTPRequest } from "./../../../services/api-service";
import { setItemInLS } from "./../../../services/storage-service";

import './registration.css'

function Registration(props) {
  const [registrationErrorMessage, setErrorMessage] = useState(null);
  const [username, updateUsername] = useState(null);
  const [password, updatePassword] = useState(null);

  let typingTimer = null;


  const handlePasswordChange = async (event) => {
    updatePassword(event.target.value);
  }

  const handleKeyDownChange = (event) => {
    clearTimeout(typingTimer);
  }

  const handleKeyUpChange = (event) => {
    const username = event.target.value;
    typingTimer = setTimeout( () => {
      checkIfUsernameAvailable(username);
    }, 1200);
  }

  const checkIfUsernameAvailable = async (username) => {  
    props.displayPageLoader(true);
    const isUsernameAvailableResponse = await isUsernameAvailableHTTPRequest(username);
    props.displayPageLoader(false);
    if (!isUsernameAvailableResponse.response.isUsernameAvailable) {
      setErrorMessage(isUsernameAvailableResponse.message);
    } else {
      setErrorMessage(null);
    }
    updateUsername(username);
  }

  const registerUser = async () => {
    props.displayPageLoader(true);
    const userDetails = await registerHTTPRequest(username, password);
    props.displayPageLoader(false);

    if (userDetails.code === 200) {
      setItemInLS('userDetails', userDetails.response)
      props.history.push(`/home`)
    } else {
      setErrorMessage(userDetails.message);
    }
  };

  return (
    <div className="app__register-container">
       <div className="app__form-row">
        <label>Username:</label>
        <input type="email" className="email" onKeyDown={handleKeyDownChange}  onKeyUp={handleKeyUpChange}/>
      </div>
      <div className="app__form-row">
        <label>Password:</label>
        <input type="password" className="password" onChange={handlePasswordChange}/>
      </div>
      <div className="app__form-row">
        <span className="error-message">{registrationErrorMessage? registrationErrorMessage : ''}</span>
      </div>
      <div className="app__form-row">
        <button onClick={registerUser}>Registration</button>
      </div>
    </div>
  );
}

export default withRouter(Registration);