// main-page.js
import { fetchAndDisplayPosts } from "/static/js/posts.js";
import { displayUserList } from "/static/js/users.js";
export function mainPage() {
    let content =
        `<div class="main-page">
            <header class="header">
                <a href="/" id="logo">
                    <img src="images/OIP.jpeg" alt="Profile Picture" class="profile-pic">
                </a>

                <nav class="navbar">
                    <button type="button" id="openWindowBtn" class="icon-button">
                        <span class="material-icons">message</span>
                        <span id="totalNotifNum" class="icon-button__badge"></span>
                    </button>

                    <div class="profile-container">
                        <img src="images/profile.png" alt="Profile Picture" class="profile-pic">
                        <span id="displayName"></span> <!-- Go template placeholder -->
                    </div>

                    <div class="logout-button">
                        <button type="button" id="logout-btn">Logout</button>
                    </div>
                </nav>
            </header>

            <main>
                <div class="categories">
                    <h2>Categories</h2>
                    <ul>
                        <li><a href="#" id="Cat1">Category 1</a></li>
                        <li><a href="#" id="Cat2">Category 2</a></li>
                        <li><a href="#" id="Cat3">Category 3</a></li>
                        <li><a href="#" id="Cat4">Category 4</a></li>
                        <li><a href="#" id="Cat5">Category 5</a></li>
                    </ul>
                </div>
                <div id="middle">
                    <div id="midTop">
                        <div id="addPost">
                            What's on your mind?
                        </div>
                        <div id="postsForm">
                            <form action="/addPost" method="POST" id="createPost">
                                <input type="text" id="postFormTitle" name="title" placeholder="Title">
                                <textarea rows="" cols="62" name="content" id="postFormContent" placeholder="Write your post here"></textarea>
                                <label id="postFormCategory"> Category: <select name="category">
                                    <option value="">None</option>
                                    <option value="Category1">Category 1</option>
                                    <option value="Category2">Category 2</option>
                                    <option value="Category3">Category 3</option>
                                    <option value="Category4">Category 4</option>
                                    <option value="Category5">Category 5</option>
                                </select></label>
                            <button type="submit">Post</button>
                            </form>
                        </div>
                    </div>
                    <div id="feed-container">
                        <!-- Dynamic posts will be added here -->
                    </div>
                </div>
                <div id="userList">
                    <!-- Dynamic user list will be added here -->
                </div>
            </main>
        </div>
        `;
    document.body.innerHTML = content;

    document.getElementById('addPost').addEventListener('click', function(e) {
        document.getElementById('postsForm').style.display = 'block'
    })
    document.getElementById('createPost').addEventListener('submit', function(e) {

        document.getElementById('postsForm').style.display = 'none'
    })
    // document.getElementById('totalNotifNum').style.display = 'none'

    document.getElementById('postsForm').style.display = 'none'

    fetchAndDisplayPosts();
    displayUserList();

    document.getElementById('Cat1').addEventListener('click', function(e) {
        [...document.getElementsByClassName('post-card')].forEach(element => {
            element.style.display = 'none'
        });
        [...document.getElementsByClassName('Category1')].forEach(element => {
            element.style.display = 'flex'
        });
    })
    document.getElementById('Cat2').addEventListener('click', function(e) {
        [...document.getElementsByClassName('post-card')].forEach(element => {
            element.style.display = 'none'
        });
        [...document.getElementsByClassName('Category2')].forEach(element => {
            element.style.display = 'flex'
        });
    })
    document.getElementById('Cat3').addEventListener('click', function(e) {
        [...document.getElementsByClassName('post-card')].forEach(element => {
            element.style.display = 'none'
        });
        [...document.getElementsByClassName('Category3')].forEach(element => {
            element.style.display = 'flex'
        });
    })
    document.getElementById('Cat4').addEventListener('click', function(e) {
        [...document.getElementsByClassName('post-card')].forEach(element => {
            element.style.display = 'none'
        });
        [...document.getElementsByClassName('Category4')].forEach(element => {
            element.style.display = 'flex'
        });
    })
    document.getElementById('Cat5').addEventListener('click', function(e) {
        [...document.getElementsByClassName('post-card')].forEach(element => {
            element.style.display = 'none'
        });
        [...document.getElementsByClassName('Category5')].forEach(element => {
            element.style.display = 'flex'
        });
    })
    document.getElementById('logo').addEventListener('click', function(e) {
        [...document.getElementsByClassName('post-card')].forEach(element => {
            element.style.display = 'flex'
        });
    })
    let logoutUser = document.getElementById('logout-btn')
    logoutUser.onclick = function(){
       window.location.href = 'http://localhost:8080/logout' 
    }



    document.getElementById('openWindowBtn').addEventListener('click', function() {
        // Open a new blank window
        const newTab = window.open('/messages', '_blank');
        document.getElementById('totalNotifNum').style.display = 'none'
        
        
        // Inject the ChatUI content from ChatUI.js (assuming it's stored as a string)
//         newTab.document.write(`
//             <html lang="en">
// <head>
//     <meta charset="UTF-8">
//     <meta name="viewport" content="width=device-width, initial-scale=1.0">
//     <title>Chat UI</title>
//     <link rel="stylesheet" href="styles.css">
// </head>
// <body>
//     <div class="container">
//         <div class="sidebar">
//             <h2>Users</h2>
//             <ul id="userList">
//                 <!-- Dynamically updated user list -->
//             </ul>
//         </div>
//         <div class="chat">
//             <div class="chat-header">
//                 <h2>Forum Chat</h2>
//             </div>
//             <div class="chat-messages" id="chatMessages">
                
//             </div>
//             <div class="chat-input">
//                 <input type="text" id="messageInput" placeholder="Type a message..." />
//                 <button id="sendButton">Send</button>
//             </div>
//         </div>
//     </div>
// </body>
// </html>
//         `);
        // newTab.document.close(); // Complete writing to the document
    });





















}

// document.getElementById('show-signup').addEventListener('click', function (e) {
//     e.preventDefault();
//     document.getElementById('login-form').style.display = 'none';
//     document.getElementById('signup-form').style.display = 'block';
// });

// document.getElementById('show-login').addEventListener('click', function (e) {
//     e.preventDefault();
//     document.getElementById('signup-form').style.display = 'none';
//     document.getElementById('login-form').style.display = 'block';
// });

// // Initially show the login form and hide the signup form
// document.getElementById('login-form').style.display = 'block';
// document.getElementById('signup-form').style.display = 'none';