export async function displayUserList() {
    try {
        const response = await fetch('http://localhost:8080/users')
        
        const userList = await response.json();
        
        const displayName = document.getElementById('displayName')
        let currentUser = userList.filter(u => u.currentUser)
        const notifsnum = document.getElementById('totalNotifNum')
        notifsnum.innerHTML = currentUser[0].totalNotifications
        if (currentUser[0].totalNotifications === 0) notifsnum.style.display = 'none'
        else notifsnum.style.display = 'flex'
        displayName.textContent = currentUser[0].nickname
        const container = document.getElementById('userList');
        const online = document.createElement('div')
        online.className = 'online'
        const onlineTitle = document.createElement('h2')
        onlineTitle.textContent = 'Online'
        online.appendChild(onlineTitle)
        const offline = document.createElement('div')
        offline.className = 'offline'
        const offlineTitle = document.createElement('h2')
        offlineTitle.textContent = 'Offline'
        offline.appendChild(offlineTitle)
        let newuserList = userList.filter(u => !u.currentUser)
        newuserList.forEach(u => {
            const user = document.createElement('p')
            user.class = 'userName'
            user.textContent = `${u.nickname}`
            if (u.status === 'Online') {
                online.appendChild(user)
            } else if (u.status === 'Offline') {
                offline.appendChild(user)
            } else {
                console.log('AJFREWJKBSJ')
            }
            container.appendChild(online)
            container.appendChild(offline)


        });
    
    } catch(error) {
        console.error('Error fetching users:', error);
    }

    
}