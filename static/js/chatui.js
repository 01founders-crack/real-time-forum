export async function messagePage() {
    try {
        let content = `
        <div class="messagePage">
            <header class="header">
                <a href="/" id="logo">
                    <img src="images/OIP.jpeg" alt="Profile Picture" class="profile-pic">
                </a>

                <nav class="navbar">

                    <div class="profile-container">
                        <img src="images/profile.png" alt="Profile Picture" class="profile-pic">
                        <span id="displayName"></span> <!-- Go template placeholder -->
                    </div>

                    <div class="logout-button">
                        <button type="button" id="logout-btn">Logout</button>
                    </div>
                </nav>
            </header>
            <div class="container">
                <div class="sidebar">
                    <h2>Users</h2>
                    
                    <ul id="userListChat">
                        <!-- Dynamically updated user list -->
                    </ul>
                </div>
                <div class="chat">
                    <div class="chat-header">
                        <h2>Forum Chat</h2>
                    </div>
                    <div id="currentChat">
                        <p>Currently talking to...</p>
                    </div>
                    <div id="chatMessages"></div>
                    <form id="messageForm" class="chat-input">
                        <input type="text" id="messageInput" name="messageInput" placeholder="Type a message..." required />
                        <input type="hidden" id="messageTargetUser" name="messageTargetUser" />
                        <button type="submit" id="sendButton">Send</button>
                    </form>
                </div>
            </div>
            </div>
        `;

        document.body.innerHTML = content;

        // WebSocket setup
        const ws = new WebSocket('ws://localhost:8080/ws');
        const userListChat = document.getElementById('userListChat');
        const chatMessages = document.getElementById('chatMessages');
        const messageInput = document.getElementById('messageInput');
        const messageTargetUser = document.getElementById('messageTargetUser');
        const currentChat = document.getElementById('currentChat');
        const displayName = document.getElementById('displayName')
        let currentUser = ""
        let userList = [];
        let currentPage = 1;
        let maxPage = 1;

        // Function to update chat messages
        function updateMessages(otherUser) {
            fetch('http://localhost:8080/users')
                .then(response => {
                    if (!response.ok) {
                        throw new Error('Network response was not ok');
                    }
                    return response.json();
                })
                .then(data => {
                    userList = [...data];
                    currentUser = userList.find(user => user.currentUser);
                    
                    
                    userList = userList.filter(user => !user.currentUser);

                    const userIndex = userList.findIndex(user => user.nickname === otherUser);

                    if (userIndex === -1) {
                        console.error(`User with nickname ${otherUser} not found.`);
                        return;
                    }

                    if (currentChat.textContent === `Currently talking to: ${userList[userIndex].nickname}`) {
                        chatMessages.innerHTML = '';
                    }
                    let currentMessages = userList[userIndex].messages || [];
                    maxPage = Math.ceil(currentMessages.length / 10);
                    if (currentPage > maxPage) currentPage = maxPage;

                    let endIndex = currentMessages.length - ((currentPage - 1) * 10);
                    let startIndex = endIndex - 10;
                    if (startIndex < 0) startIndex = 0;
                    startIndex = Math.max(0, startIndex); // Ensure startIndex isn't negative
                    endIndex = Math.min(startIndex+10, currentMessages.length)
                    currentMessages = currentMessages.slice(startIndex, endIndex);

                    currentMessages.forEach(m => {
                        const messageBubble = document.createElement('div');
                        messageBubble.classList.add('messageBubble');
                        if (m.senderName === otherUser) {
                            messageBubble.classList.add('otherUser');
                        }

                        const messageBubbleHead = document.createElement('div');
                        messageBubbleHead.classList.add('messageBubbleHead');

                        const sentBy = document.createElement('p');
                        sentBy.textContent = m.senderName;

                        const sentAt = document.createElement('p');
                        sentAt.textContent = m.time;

                        messageBubbleHead.appendChild(sentBy);
                        messageBubbleHead.appendChild(sentAt);

                        const messageContent = document.createElement('div');
                        messageContent.textContent = m.content;

                        messageBubble.appendChild(messageBubbleHead);
                        messageBubble.appendChild(messageContent);

                        chatMessages.appendChild(messageBubble);
                    });
                })
                .catch(error => {
                    console.error('There was a problem with the fetch operation:', error);
                });

            // Scroll to the latest message
        }

        // Function to update the user list
        function updateUserListChat() {
            currentChat.textContent = `Currently talking to...`
            fetch('http://localhost:8080/users')
                .then(response => {
                    if (!response.ok) {
                        throw new Error('Network response was not ok');
                    }
                    return response.json();
                })
                .then(data => {
                    userList = [...data];
                    currentUser = userList.find(user => user.currentUser);
                    displayName.textContent = currentUser.nickname
                    userList = userList.filter(user => !user.currentUser); // Remove current user from list
                    userListChat.innerHTML = '';
                    userList.forEach(user => {
                        const userItem = document.createElement('li');
                        userItem.id = user.nickname;
                        userItem.innerHTML = `${user.nickname} <span class="${user.status}">${user.status}</span>`;
                        userListChat.appendChild(userItem);
                    });

                    // Add click event listener to each user in the list
                    document.querySelectorAll('#userListChat li').forEach(userItem => {
                        userItem.addEventListener('click', e => {
                            const selectedUser = e.target.id;
                            chatMessages.className = `${selectedUser}`
                            currentChat.textContent = `Currently talking to: ${selectedUser}`;
                            messageTargetUser.value = selectedUser;
                            currentPage = 1; // Reset to first page of messages
                            updateMessages(selectedUser);
                            setTimeout(function(){chatMessages.scrollTop = chatMessages.scrollTopMax;}, 250);
                        });
                    });
                })
                .catch(error => {
                    console.error('There was a problem with the fetch operation:', error);
                });
        }

        // Handle WebSocket messages (status and chat messages)
        ws.onmessage = function (event) {
            const data = JSON.parse(event.data);
            console.log('Message received:', data);

            if (data.type === 'message') {
                // Only update the chat if the message is for the currently selected user
                const currentTarget = document.getElementById('messageTargetUser').value;
                if (data.username === currentTarget || data.target === currentTarget) {
                    updateMessages(currentTarget);
                }
            } else if (data.type === 'status') {
                // Update user status and refresh user list
                // const currentTarget = document.getElementById('messageTargetUser').value;
                updateUserListChat();
                
            }
        };

        // Send a message
        document.getElementById('messageForm').addEventListener('submit', (event) => {
            event.preventDefault(); // Prevent default form submission behavior

            const message = messageInput.value;
            const targetUser = messageTargetUser.value;
            if (message && targetUser) {
                ws.send(JSON.stringify({
                    type: 'message',
                    username: currentUser.nickname, // Replace this with actual user's name
                    message: message,
                    target: targetUser
                }));
                console.log('Message sent:', message, 'to', targetUser);
                messageInput.value = ''; // Clear the input field after sending

            }
            updateUserListChat()
            chatMessages.className = `${targetUser}`
            currentChat.textContent = `Currently talking to: ${targetUser}`;
            messageTargetUser.value = targetUser;
            currentPage = 1; // Reset to first page of messages
            updateMessages(targetUser);
            
            setTimeout(function(){chatMessages.scrollTop = chatMessages.scrollTopMax;}, 250);
            
        });

        // Handle scroll for loading more messages
        
        chatMessages.addEventListener("scroll", (event) => {
            if (event.target.scrollTop <= 0 && currentPage < maxPage){
                currentPage++
                setTimeout(updateMessages(event.target.classList[0]),100)
                setTimeout(function(){chatMessages.scrollTop = chatMessages.scrollTopMax - 1;}, 250)
            }
            if (event.target.scrollTop >= event.target.scrollTopMax && currentPage > 1){
                currentPage--
                setTimeout(updateMessages(event.target.classList[0]),100)
                setTimeout(function(){chatMessages.scrollTop = 1;}, 250)
            }
        });
        document.getElementById('logo').addEventListener('click', function(e) {
            [...document.getElementsByClassName('post-card')].forEach(element => {
                element.style.display = 'flex'
            });
        })
        let logoutUser = document.getElementById('logout-btn')
        logoutUser.onclick = function(){
           window.location.href = 'http://localhost:8080/logout' 
        }

        // Initial call to update user list
        updateUserListChat();
    } catch (error) {
        console.error('Error initializing chat page:', error);
    }
}
