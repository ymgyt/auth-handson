export class ServerSideAuth {
    constructor(document) {
        this.dom = {}
        this.dom.debug = document.getElementById('debug')
        this.dom.email = document.getElementById('email')
        this.dom.password = document.getElementById('password')
        this.dom.signUp = document.getElementById('sign-up')
        this.dom.signIn = document.getElementById('sign-in')
    }

    printDebug(msg) {
        this.dom.debug.textContent = msg
    }

    email() {
        return this.dom.email.value
    }
    password() {
        return this.dom.password.value
    }
}

export class HandsOn {
    constructor(document) {
        this.dom = {}
        this.dom.email = document.getElementById('input-email')
        this.dom.password = document.getElementById('input-password')
        this.dom.signOut = document.getElementById('sign-out')

        this.dom.debugMsg = document.getElementById('debug-message')
        this.dom.debugData = document.getElementById('debug-data')
        this.dom.debugJWT = document.getElementById('debug-jwt')

        this.dom.customSignUp = document.getElementById('custom-sign-up')
        this.dom.customSignIn = document.getElementById('custom-sign-in')

        this.dom.firebaseSignUp = document.getElementById("firebase-sign-up")
        this.dom.firebaseSignIn = document.getElementById("firebase-sign-in")
        this.dom.firebaseGoogleSignIn = document.getElementById("firebase-google-sign-in")
        this.dom.firebaseFacebookSignIn = document.getElementById("firebase-sign-in")

        this.dom.auth0SignUp = document.getElementById('auth0-sign-up')
        this.dom.auth0SignIn = document.getElementById('auth0-sign-in')
        this.dom.auth0ProtectedResource = document.getElementById('auth0-protected-resource')
        this.dom.auth0SignOut = document.getElementById('auth0-sign-out')
    }

    email() {
        return this.dom.email.value;
    }
    password() {
        return this.dom.password.value;
    }
    oauth_flow() {
        return document.querySelector('input[name="oauth-flow-handle"]:checked').value
    }

    printMsg(msg) {
        this.dom.debugMsg.textContent = msg
    }
    appendMsg(msg) {
        const org = this.dom.debugMsg.textContent
        this.dom.debugMsg.textContent = `${org}\n${msg}`
    }
    printData(data) {
        this.dom.debugData.textContent = JSON.stringify(data, null, ' ')
    }
    print(msg, data) {
        this.printMsg(msg)
        this.printData(data)
    }
    printFail(msg, data) {
        this.appendMsg(msg)
        this.printData(data)
    }

    jwt() {
        return this.dom.debugJWT.textContent
    }

    setJWT(jwt) {
        this.dom.debugJWT.textContent = jwt
    }
}