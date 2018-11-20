function toggleSignIn(auth) {
    return function() {
        if (auth.currentUser) {
            auth.signOut();
        } else {
            let email = document.getElementById('email').value;
            let password = document.getElementById('password').value;

            auth.signInWithEmailAndPassword(email, password).catch(function(error) {
                console.log("fail auth.signInWithEmailAndPassword", error)
            })
        }
    }
}

function handleSignUp(auth) {
    return function() {
        console.log("handle SignUp");

        let email = document.getElementById('email').value;
        let password = document.getElementById('password').value;
        if (email.length < 3) {
            alert('Please enter an email address.');
            return;
        }
        if (password.length < 6) {
            alert('Password should at least 6 characters.');
            return;
        }

        // Sign in with email and pass.
        auth.createUserWithEmailAndPassword(email, password).catch(function(error) {
            console.log("fail auth.createUserWithEmailAndPassword", error);
        });
    };
}

function handleAuthStateChange(auth) {
    auth.onAuthStateChanged(function(user){
        console.log("auth state changed")
        if (user){
            // user is signed in.
           document.getElementById('quickstart-sign-in-status').textContent = 'Signed in';
           document.getElementById('quickstart-sign-in').textContent = 'Signed out';
           document.getElementById('quickstart-account-details').textContent = JSON.stringify(user, null, ' ');

           // provider specific data
           let ps = document.getElementById('quickstart-account-details-provider-specific')
           ps.textContent = '';
           user.providerData.forEach(function(profile) {
               if (profile) {
                   ps.textContent = ps.textContent + (JSON.stringify(profile, null, ' '));
               }
           })

           if (!user.emailVerified) {
               document.getElementById('quickstart-verify-email').disabled = false;
           }
        } else {
            // user is signed out.
           document.getElementById('quickstart-sign-in-status').textContent = 'Signed out';
           document.getElementById('quickstart-sign-in').textContent = 'Signed in';
           document.getElementById('quickstart-account-details').textContent = 'null';
        }

        document.getElementById('quickstart-sign-in').disabled = false;
    })
}

function handleSendEmailVerification(auth) {
    return function() {
        auth.currentUser.sendEmailVerification().then(function() {
            alert(`email verification sent to ${auth.currentUser.email}`)
        })
    }
}

function handleUpdateEmail(auth) {
    return function(){
        let user = auth.currentUser;
        let orgEmail = user.email
        let newEmail = document.getElementById('email').value;

        if (newEmail.length < 3) {
            alert("email required !");
            return;
        }

        user.updateEmail(newEmail).then(function() {
            console.log(`update email from ${orgEmail} to ${newEmail}`)
        }).catch(function(error){
            console.log(`failed auth.currentUser.updateEmail`, error);
        })
    }
}

function init() {

    // Initialize Firebase
    let config = {
        apiKey: "AIzaSyA9Ml5AttyWhy3uc7bhrwViMOrfej6OSik",
        authDomain: "id-integration-handson.firebaseapp.com",
        databaseURL: "https://id-integration-handson.firebaseio.com",
        projectId: "id-integration-handson",
        storageBucket: "id-integration-handson.appspot.com",
        messagingSenderId: "265757313676"
    };

    let auth = firebase.initializeApp(config).auth();
    console.log(`init firebase auth ${auth.app.name}`);

    handleAuthStateChange(auth);

    document.getElementById('quickstart-sign-in').addEventListener('click', toggleSignIn(auth), false);
    document.getElementById('quickstart-sign-in-status').textContent = 'Signed out';
    document.getElementById('quickstart-sign-up').addEventListener('click', handleSignUp(auth), false);
    document.getElementById('quickstart-verify-email').addEventListener('click', handleSendEmailVerification(auth), false);
    document.getElementById('quickstart-update-email').addEventListener('click', handleUpdateEmail(auth), false);

}

function main() {
    'use strict';
    switch (document.readyState) {
        case "loading":
            document.addEventListener("DOMContentLoaded", init);
            break;
        case "interactive":
        case "complete":
            init();
            break;
        default:
            alert(`unexpected document ready stage ${document.readyState}`)
    }
}

main();