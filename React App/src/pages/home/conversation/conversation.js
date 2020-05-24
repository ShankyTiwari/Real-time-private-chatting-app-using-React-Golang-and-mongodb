import React, { useState, useEffect, useRef } from 'react';
import {
  eventEmitter,
  sendWebSocketMessage,
} from './../../../services/socket-service';
import { getConversationBetweenUsers } from './../../../services/api-service';

import './conversation.css'

const alignMessages = (userDetails, toUserID) => {
  const { userID } = userDetails;
  return userID !== toUserID;
}

const scrollMessageContainer = (messageContainer) => {
  if (messageContainer.current !== null) {
    try {
      setTimeout(() => {
        messageContainer.current.scrollTop = messageContainer.current.scrollHeight;
      }, 100);
    } catch (error) {
      console.warn(error);
    }
  }
}

const getMessageUI = (messageContainer, userDetails, conversations) => {
  return (
    <ul ref={messageContainer} className='message-thread-container'>
      {conversations.map((conversation, index) => (
        <li
          className={`message ${
            alignMessages(userDetails, conversation.toUserID) ? 'align-right' : ''
          }`}
          key={index}
        >
          {conversation.message}
        </li>
      ))}
    </ul>
  );
}

const getInitiateConversationUI = (userDetails) =>{
  if (userDetails !== null) {
    return (
      <div className="message-thread-container start-chatting-banner">
        <p className="heading">
          You haven 't chatted with {userDetails.username} in a while,
          <span className="sub-heading"> Say Hi.</span>
        </p>			
      </div>
    )
  }    
}

function Conversation(props) {
  const selectedUser = props.selectedUser;
  const userDetails = props.userDetails;

  const messageContainer = useRef(null);
  const [conversation, updateConversation] = useState([]);
  const [messageLoading, updateMessageLoading] = useState(true);

  useEffect(() => {
    if (userDetails && selectedUser) {
      (async () => {
        const conversationsResponse = await getConversationBetweenUsers(userDetails.userID, selectedUser.userID);
        updateMessageLoading(false)
        if (conversationsResponse.response) {
          updateConversation(conversationsResponse.response);
        } else if (conversationsResponse.response === null) {
          updateConversation([]);
        }
      })();
    }
  }, [userDetails, selectedUser])

  useEffect(() => {
    const newMessageSubscription = (messagePayload) => {
      if (
        selectedUser !== null &&
        selectedUser.userID === messagePayload.fromUserID
      ) {
        updateConversation([...conversation, messagePayload]);
        scrollMessageContainer(messageContainer);
      }
    };

    eventEmitter.on('message-response', newMessageSubscription);

    return () => {
      eventEmitter.removeListener('message-response', newMessageSubscription);
    };
  }, [
    conversation,
    selectedUser
  ]);

  const sendMessage = (event) => {
    if (event.key === 'Enter') {
      const message = event.target.value;

      if (message === '' || message === undefined || message === null) {
        alert(`Message can't be empty.`);
      } else if (userDetails.userID === '') {
        this.router.navigate(['/']);
      } else if (selectedUser === undefined) {
        alert(`Select a user to chat.`);
      } else {
        event.target.value = '';

        const messagePayload = {
          fromUserID: userDetails.userID,
          message: message.trim(),
          toUserID: selectedUser.userID,
        };

        sendWebSocketMessage(messagePayload);
        updateConversation([...conversation, messagePayload]);
        scrollMessageContainer(messageContainer);
      }
    }
  }

  if (messageLoading) {
    return (
      <div
        className="message-overlay"
      >
        <h3>
          {selectedUser !== null && selectedUser.username
            ? 'Loading Messages'
            : ' Select a User to chat.'}
        </h3>
      </div>
    )
  }

  return (
    <div className='app__conversion-container'>
      
      {conversation.length > 0
        ? getMessageUI(messageContainer, userDetails, conversation)
        : getInitiateConversationUI(selectedUser)}

      <div className='app__text-container'>
        <textarea
          placeholder={`${
            selectedUser !== null ? '' : 'Select a user and'
          } Type your message here`}
          className='text-type'
          onKeyPress={sendMessage}
        ></textarea>
      </div>
    </div>
  );
}

export default Conversation;