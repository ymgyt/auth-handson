import { HandsOn as UI } from '/js/module.js';

function signOut(ui, auth) {
    return () => {
        ui.print("sign out", {})
        auth.signOut()
    }
}

function customSignUp(ui) {
    return () => {
        ui.printMsg("sign up start")

        const email = ui.email()
        const password = ui.password()

        fetch('/server/auth/sign-up',{
            method: 'POST',
            headers: {
                "Content-Type": "application/json; charset=utf8",
            },
            body: JSON.stringify({
                email: email,
                password: password,
            })
        })
        .then(response => {
            if (!response.ok) {
                ui.appendMsg("sign up fail")
            } else {
                ui.appendMsg("sign up success")
            }
            response.json().then(json => ui.printData(json))
        })
    }
}

function customSignIn(ui, auth) {
    return () => {
        ui.printMsg("sign in start")

        const email = ui.email()
        const password = ui.password()

        fetch('/server/auth/sign-in',{
            method: 'POST',
            headers: {
                "Content-Type": "application/json; charset=utf8",
            },
            body: JSON.stringify({
                email: email,
                password: password,
            })
        })
        .then(response => {
            if (!response.ok) {
                ui.appendMsg("sign in fail")
                response.json().then(json => ui.printData(json))
                return
            }

            response.json().then(json =>{
                ui.appendMsg(`fetch custom token ${json.token}`)

                auth.signInWithCustomToken(json.token)
                .then(userCredential => {
                    ui.appendMsg("sign in with custom token success")
                    ui.printData(userCredential)
                })
                .catch(err => console.log("firebase.auth.signInWithCustomToken", err))
            })
        })
    }
}

function googleSignIn(ui, auth) {
    return () => {
        ui.printMsg("google sign in")

        const provider = new firebase.auth.GoogleAuthProvider();
        provider.addScope('https://www.googleapis.com/auth/contacts.readonly')

        if (ui.oauth_flow() == "popup") {
            auth.signInWithPopup(provider)
            .then(result => {
                ui.appendMsg("google sign in success")
                ui.printData(result)
            })
            .catch(err => ui.printFail("google sign in fail", err))
        } else {
            auth.signInWithRedirect(provider)
            ui.appendMsg("redirect user")
            auth.getRedirectResult()
            .then(result => {
                ui.appednMsg("google sign in success")
                ui.printData(result)
            })
            .catch(err => ui.printFail("google sign in fail", err))
        }
    }
}

function auth0SignUp(ui) {
    return () => {
        ui.printMsg("auth0 sign up start")

        const email = ui.email()
        const password = ui.password()

        fetch('/auth0/sign-up',{
            method: 'POST',
            headers: {
                "Content-Type": "application/json; charset=utf8",
            },
            body: JSON.stringify({
                email: email,
                password: password,
            })
        })
        .then(response => {
            if (!response.ok) {
                ui.appendMsg("auth0 sign up fail")
            } else {
                ui.appendMsg("auth0 sign up success")
            }
            response.json().then(json => ui.printData(json))
        })
    }
}

function auth0SignIn(ui) {
    return () => {
        ui.printMsg("auth0 sign in start")

        const email = ui.email()
        const password = ui.password()

        fetch('/auth0/sign-in',{
            method: 'POST',
            headers: {
                "Content-Type": "application/json; charset=utf8",
            },
            body: JSON.stringify({
                email: email,
                password: password,
            })
        })
        .then(response => {
            if (!response.ok) {
                ui.appendMsg("auth0 sign in fail")
            } else {
                ui.appendMsg("auth0 sign in success")
            }
            response.json().then(json => {
                ui.printData(json)
                if (response.ok) {
                    ui.setJWT(json.id_token)
                }
            })
        })
    }
}

function auth0SignOut(ui) {
    return () => {
        localStorage.removeItem('access_token');
        localStorage.removeItem('id_token');
        localStorage.removeItem('expires_at');
    }
}

function auth0ProtectedResource(ui) {
    return () => {
        ui.printMsg("auth0 fetch private resource")

        const jwtToken = ui.jwt()

        fetch('/auth0/private/resource', {
            method: 'GET',
            headers: {
                "Authorization": `bearer ${jwtToken}`,
            }
        })
        .then(response => {
            if (!response.ok) {
                ui.appendMsg("auth0 fetch private resource fail")
            } else {
                ui.appendMsg("auth0 fetch private resource success")
            }
            response.json().then(json => ui.printData(json))
        })
    }
}

function handleAuth0Authentication(ui,webAuth) {
    webAuth.parseHash(function(err, authResult) {
        if (authResult && authResult.accessToken && authResult.idToken) {
            ui.printData(authResult)
            ui.setJWT(authResult.idToken)
        } else {
            ui.printData(err)
        }
        window.location.hash = ''
    })
}

function initAuth0(ui) {

    const webAuth = new auth0.WebAuth({
        domain: 'ymgyt.auth0.com',
        clientID: 'fVA2Tsbu0OT27v12MGur6E3uDB3pfMpx',
        responseType: 'token id_token',
        scope: 'openid profile email',
        redirectUri: window.location.href,
        realm: 'auth-handson',
        // leeway: 5,
    })

    const loginBtn = document.getElementById('auth0-cs-sign-in')
    loginBtn.addEventListener('click', event => {
        event.preventDefault()
        webAuth.authorize({
             // connection: 'auth-handson', // 指定するとsocialがでなくなる
        })
    })

    handleAuth0Authentication(ui,webAuth)
}

function init() {
    // firebaseのwarningが消せないので
    console.clear()
    console.log("init")
    const ui = new UI(document)

    const fbConfig = {
        apiKey: "AIzaSyA9Ml5AttyWhy3uc7bhrwViMOrfej6OSik",
        authDomain: "id-integration-handson.firebaseapp.com",
        databaseURL: "https://id-integration-handson.firebaseio.com",
        projectId: "id-integration-handson",
        storageBucket: "id-integration-handson.appspot.com",
        messagingSenderId: "265757313676"
    };

    const auth = firebase.initializeApp(fbConfig).auth();

    ui.dom.customSignUp.addEventListener('click', customSignUp(ui), false)
    ui.dom.customSignIn.addEventListener('click', customSignIn(ui, auth), false)
    ui.dom.firebaseGoogleSignIn.addEventListener('click', googleSignIn(ui, auth), false)
    ui.dom.signOut.addEventListener('click', signOut(ui, auth), false)
    ui.dom.auth0SignUp.addEventListener('click', auth0SignUp(ui), false)
    ui.dom.auth0SignIn.addEventListener('click', auth0SignIn(ui), false)
    ui.dom.auth0ProtectedResource.addEventListener('click', auth0ProtectedResource(ui), false)
    ui.dom.auth0SignOut.addEventListener('click', auth0SignOut(ui), false)

    ui.printMsg("hello")

    initAuth0(ui)
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
            alert(`unexpected document ready stage ${document.readyState}`);
    }
}

main();