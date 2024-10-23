let messages = [];
let currentPage = 1;
const messagesPerPage = 10;

async function fetchAndDisplayMessages() {
    try {
        const response = await fetch('http:localhost:8080/messages');
        messages = await response.json();
        userId = getCookieValue('userId')
        displayMessages();
    } catch (error) {
        console.error('ERROR FETCHING MESSAGES', error);
    }
}

function getCookieValue(cookieName) {
    let cookies = document.cookie.split(';');
    for (let cookie of cookies) {
        let [name, value] = cookie.trim().split('=');
        if (name === cookieName) {
            return value;
        }
    }
    return null
}

function displayMessages() {
    const container = document.getElementById('messageBox');
    container.innerHTML = '';

    const startIndex = (currentPage - 1) * messagesPerPage;
    const endIndex = startIndex + messagesPerPage;
    const messagesToDisplay = messages.slice(startIndex, endIndex);

    messagesToDisplay.forEach(message => {
        const messageCard = document.createElement('div');
        messageCard.className = 'messageCard';
        if (username.senderId === userId) {
            messageCard.classList.add('currentUser')
        } else {
            messageCard.classList.add('targetUser')
        }

            const details = document.createElement('div');
            details.className = 'messageDetails'

                const username = document.createElement('p');
                username.className = 'messageSender'
                username.textContent = `${message.senderName}`

                const timeSent = document.createElement('p');
                timeSent.className = 'timeSent';
                timeSent.textContent = `${message.createdAt}`

                details.appendChild(username)
                details.appendChild(timeSent)

            const messageBody = document.createElement('p');
            messageBody.textContent = `${message.content}`
        
            messageCard.appendChild(details)
            messageCard.appendChild(messageBody)

       container.appendChild(messageCard)
    })
}