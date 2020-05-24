import React, {useState, useEffect} from 'react';
import { withRouter } from 'react-router-dom';

import { userSessionCheckHTTPRequest } from "./../../services/api-service";
import {
  connectToWebSocket,
  listenToWebSocketEvents,
  emitLogoutEvent,
} from './../../services/socket-service';
import {
  getItemInLS,
  removeItemInLS
} from "./../../services/storage-service";

import ChatList from './chat-list/chat-list';
import Conversation from './conversation/conversation';

import './home.css';

const useFetch = (props) => {
  
  const [internalError, setInternalError] = useState(null);
  const userDetails = getItemInLS('userDetails');
      
  useEffect(() => {

    (async () => {
      if (userDetails === null || userDetails === '') {
        props.history.push(`/`);
      } else {
        const isUserLoggedInResponse = await userSessionCheckHTTPRequest(
          userDetails.userID
        );
        if (!isUserLoggedInResponse.response) {
          props.history.push(`/`);
        } else {
          const webSocketConnection = connectToWebSocket(userDetails.userID);
          if (webSocketConnection.webSocketConnection === null) {
            setInternalError(webSocketConnection.message);
          } else {
            listenToWebSocketEvents()
          }
        }
      }
    })();

  }, [props, userDetails]);
  return [userDetails, internalError];
};

const getUserNameInitial = (userDetails) => {
  if(userDetails && userDetails.username) {
    return userDetails.username[0]
  }
  return '_';
}


const getUserName = (userDetails) => {
  if (userDetails && userDetails.username) {
    return userDetails.username;
  }
  return '___';
};

const logoutUser = (props, userDetails) => {
  if (userDetails.userID) {
    removeItemInLS('userDetails');
    emitLogoutEvent(userDetails.userID);
    props.history.push(`/`);
  }
};



function Home(props) {
  const [userDetails, internalError] = useFetch(props);
  const [selectedUser, updateSelectedUser] = useState(null);

  if (internalError !== null) {
    return <h1>{internalError}</h1>;
  }

  return (
    <div className='app__home-container'>
      <header className='app__header-container'>
        <nav className='app__header-user'>
          <div className='username-initial'>
            {getUserNameInitial(userDetails)}
          </div>
          <div className='user-details'>
            <h4>{getUserName(userDetails)}</h4>
          </div>
        </nav>
        <button className='logout' href='#' onClick={ () => {
          logoutUser(props, userDetails);
        }} >
          Logout
        </button>
      </header>
      <div className='app__content-container'>
        <div className='app__home-chatlist'>
          <ChatList
            updateSelectedUser={(user) => {
              updateSelectedUser(user);
            }}
            userDetails={userDetails}
          />
        </div>
        <div className='app__home-message'>
          <Conversation userDetails={userDetails} selectedUser={selectedUser} />
        </div>
      </div>
    </div>
  );
}

export default withRouter(Home);