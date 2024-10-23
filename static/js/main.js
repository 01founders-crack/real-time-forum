// const routes = [
//     {
//         path: '/',
//         pathName: 'Dashboard',
//         action: nil
//     },
//     {
//         path: '/login',
//         pathName: 'Login-Signup',
//         action: nil
//     },
//     {
//         path: '/posts',
//         pathName: 'Get-Posts',
//         action: nil
//     },
//     {
//         path: '/addPost',
//         pathName: 'Create-Post',
//         action: nil
//     },
//     {
//         path: '/addComment',
//         pathName: 'Create-Comment',
//         action: nil
//     },
//     {
//         path: '/users',
//         pathName: 'Get-Users',
//         action: nil
//     },
//     {
//         path: '/messages',
//         pathName: 'Get-Messages',
//         action: nil
//     },
//     {
//         path: '/addMessage',
//         pathName: 'Create-Message',
//         action: nil
//     }
    
// ]

import { mainPage } from "/static/js/main-page.js";
import { signinSignup } from "/static/js/signin-signup.js";
import { messagePage } from "/static/js/chatui.js"

// creeates empty routing maps
let routes = {};
let templates = {};

//generates key: value pair of route name and function eg. template('Dashboard', mainpage) => 'Dashboard': mainpage; and adds to templates map
let template = (routeName, action) => {
    return templates[routeName] = action;
}
// generates key: value pair of route path and template eg. route('/', 'Dashboard') => '/': mainpage; and adds it to route map
let route = (path, template) => {
    if (typeof template == 'function') {
        return routes[path] = template;
    } else if (typeof template == 'string') {
        return routes[path] = templates[template]
    }
}

// if route exists return route function, else return error
let resolveRoute = (route) => {
    try {
     return routes[route];
    } catch (error) {
        throw new Error("The route is not defined");
    }
};

// find current url, check if url is valid, run route function or run error
let router = (evt) => {
    const url = window.location.href.slice('http://localhost:8080'.length);
    const routeResolved = resolveRoute(url);
    routeResolved();
};
//create templates for pages
template('Dashboard', mainPage)
template('Login', signinSignup)
template('Messages', messagePage)

//create routes to pages
route('/', 'Dashboard')
route('/login', 'Login')
route('/messages', 'Messages')


//checks router on page load
window.addEventListener('load', router);
//checks router on page/url change
window.addEventListener('hashchange', router);