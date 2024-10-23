// signin-signup.js

export function signinSignup() {
    let content =
        ` <div class="container">
            <!-- Login Form -->
            <div class="login" id="login-form">
                <form action="/signin" method="POST" id="Ibraheem" >
                <h1>Sign in to your account</h1>
                <p>If you havent signed up yet.<br><a href="#" id="show-signup">Register here!</a></p>
                <input type="text" name="input-login-email" id="input-login-email" placeholder="User name or Email address " required>
                <br />
                <input type="password" name="input-login-password" id="input-login-password" placeholder="Password" required>
                <div class="options">
                    <label><input type="checkbox"> Remember me</label>
                    <a href="#">Forgot password?</a>
                </div>
                <button type="submit" id="Login-Button" >Sign In </button>
                <p>Or continue with</p>
                <div class="social-login">
                    <button class="facebook">Facebook</button>
                    <button class="twitter">Twitter</button>
                    <button class="github">GitHub</button>
                </div>
            </form>
            </div>
            
            <!-- Signup Form -->
            <div class="signup-form" id="signup-form">
                <form action="/signup" method="POST" id="Lotfi">
                  <h1>Sign up to get started</h1>
                  <p>If you already have an account, <br><a href="#" id="show-login">Login here!</a></p>
                    <input type="text" id="nickname" name="nickname" placeholder="Nickname">
                    <input type="number" id="age" name="age" placeholder="Age">
                    <select id="gender" name="gender">
                        <option value="" disabled selected>Gender</option>
                        <option value="male">Male</option>
                        <option value="female">Female</option>
                        <option value="other">Other</option>
                    </select>
                    <input type="text" id="first-name" name="firstName" placeholder="First Name">
                    <input type="text" id="last-name" name="lastName" placeholder="Last Name">
                    <input type="email" id="email" name="email" placeholder="E-mail">
                    <input type="password" id="password" name="password" placeholder="Password">
                    <button type="submit">Get Started</button>
                
                  <p>Or continue with</p>
                  <div class="social-login">
                       <button class="facebook">Facebook</button>
                       <button class="twitter">Twitter</button>
                       <button class="github">GitHub</button>
                 </div>
                </form>
            </div>
            
            <div class="image"></div>
            </div>
        `;

    // function that switches between sign in and registration forms
    document.body.innerHTML = content;

    document.getElementById('show-signup').addEventListener('click', function (e) {
        e.preventDefault();
        document.getElementById('login-form').style.display = 'none';
        document.getElementById('signup-form').style.display = 'block';
    });

    document.getElementById('show-login').addEventListener('click', function (e) {
        e.preventDefault();
        document.getElementById('signup-form').style.display = 'none';
        document.getElementById('login-form').style.display = 'block';
    });

    document.getElementById('Ibraheem').addEventListener('submit', function (e) {
        console.log('submit')
    })

    // Initially show the login form and hide the signup form
    document.getElementById('login-form').style.display = 'block';
    document.getElementById('signup-form').style.display = 'none';
}
