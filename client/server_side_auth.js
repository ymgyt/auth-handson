import { ServerSideAuth as UI } from '/js/module.js';

function signUp(ui) {

    return () => {
        ui.printDebug("in process ...")

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
        .then(res => res.json())
        .then(stream => {
            ui.printDebug(JSON.stringify(stream, null, 2))
        })
        .catch(err=>{
            ui.printDebug(JSON.stringify(err))
        })
    }
}

function signIn(ui, auth) {
    return () => {
        ui.printDebug("in process ...")

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
        .then(res => res.json())
        .then(data => {
            const token = data.token
            auth.signInWithCustomToken(token)
            .then(userCred => ui.printDebug(JSON.stringify(userCred, null, ' ')))
            .catch(err => ui.printDebug(JSON.stringify(err, null, ' ')))
        })
        .catch(err=>{
            ui.printDebug(JSON.stringify(err))
        })
    }
}

function init() {
    console.log("init")
    const ui = new UI(document)

    const config = {
        apiKey: "AIzaSyA9Ml5AttyWhy3uc7bhrwViMOrfej6OSik",
        authDomain: "id-integration-handson.firebaseapp.com",
        databaseURL: "https://id-integration-handson.firebaseio.com",
        projectId: "id-integration-handson",
        storageBucket: "id-integration-handson.appspot.com",
        messagingSenderId: "265757313676"
    };

    const auth = firebase.initializeApp(config).auth();
    console.log(auth)

    ui.dom.signUp.addEventListener('click', signUp(ui), false)
    ui.dom.signIn.addEventListener('click', signIn(ui, auth), false)
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