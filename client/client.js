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

var currentUser;

function handleAuthStateChange(auth) {
    auth.onAuthStateChanged(function(user){
        console.log("auth state changed")
        currentUser = user;
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
           console.log("sign out")
           document.getElementById('quickstart-sign-in-status').textContent = 'Signed out';
           document.getElementById('quickstart-sign-in').textContent = 'Signed in';
           document.getElementById('quickstart-account-details').textContent = 'null';
           document.getElementById('quickstart-account-details-provider-specific').textContent = 'null';
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

function googleSignIn(auth) {
    let redirected = false
    return () => {
        console.log("google sign in")

        let provider = new firebase.auth.GoogleAuthProvider();
        provider.addScope('https://www.googleapis.com/auth/contacts.readonly');

        auth.useDeviceLanguage();

        /*
        provider.setCustomParameters({
            'login_hint': 'user@example.com'
        })
        */
       let popup = false
       if (popup) {
           auth.signInWithPopup(provider).then(result => {
               console.log("google sign-in success", result);
           }).catch(error => {
               console.log("google sign-in fail ", error);
           })
       } else {
           // redirect
           if (!redirected) {
               auth.signInWithRedirect(provider);
               redirected = true;
           } else {
               auth.getRedirectResult().then(result => {
                   console.log("google sign-in success ", result);
               }).catch(error => {
                   console.log("google sign-in fail ", error);
               })
           }
       }
    }
}

function facebookSignIn(auth) {
    return () => {
        console.log("facebook sign in")

        let provider = new firebase.auth.FacebookAuthProvider();
        let scope = "email"

        provider.addScope(scope)

        auth.signInWithPopup(provider).then(result => {
            console.log("facebook sign-in success", result);
        }).catch(error => {
            console.log("facebook sign-in fail", error);
        })
    }
}

function db_save(db) {
    let keyIn = document.getElementById('dbKey');
    let valueIn = document.getElementById('dbValue');

    return () => {
        let key = keyIn.value;
        let value = valueIn.value;
        if (key.lengh == 0 || value.length == 0) {
           alert("key & value required !") ;
           return;
        }

        if (!currentUser) {
            alert("sign in required !");
            return;
        }
        let uid = currentUser.uid;
        if (!uid) {
            console.log("db_save I NEED uid");
            return;
        }

        console.log(`db_set ${key}:${value}`)

        db.collection("users").doc(uid).set({
            [key]: value,
        },{
            merge: true
        })
        .then(() => {
            console.log("db_save success");
        })
        .catch(error => {
            console.log("db_save fail", error);
        })

        keyIn.value = '';
        valueIn.value = '';

        db_display(db);
    }
}

function db_display(db) {
    let container = document.getElementById('dbContainer')

    container.innerHTML = ''

    if (!currentUser) {
        return ;
    }
    let uid = currentUser.uid;
    // ここの判定あやしい
    if (!uid) {
        console.log("db_display: user not signin");
        return;
    }
    let docSnapshot = db.collection("/users").doc(uid).get()
    .then(docSnapshot => {
        let data = docSnapshot.data();
        let div = document.createElement("div");
        div.innerHTML = `<p>${JSON.stringify(data, null, ' ')}</p>`;
        container.appendChild(div);
    })
    .catch(error => {
        console.log("db_diplay fail", error);
    })

}

function firestore(auth) {
    let db = firebase.firestore();

    // warningに従い以下の設定を追加
    db.settings({
        timestampsInSnapshots: true,
    });

    document.getElementById('dbSave').addEventListener('click', db_save(db), false);

    auth.onAuthStateChanged(user => {
        db_display(db);
    })
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
    document.getElementById('google-sign-in').addEventListener('click', googleSignIn(auth), false);
    document.getElementById('facebook-sign-in').addEventListener('click', facebookSignIn(auth), false);

    // cloud firestore
    firestore(auth);
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