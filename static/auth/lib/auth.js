Vue.createApp({
    async created() {
        const authToken = localStorage.getItem(`authToken`)
        if (authToken) {
            const response = await this.requestCheck(authToken)
            if (response && response.user && response.session) {
                this.setAuthenticationFromCheckResponse(response)
                this.screen = this.screens.authenticated
            } else {
                localStorage.removeItem(`authToken`)
            }
        }
        this.loaded = true
    },
    data() {
        return {
            loaded: false,
            loginByCodeEnabled: false,
            loginByCode: {
                inputCode: false,
                retryInterval: 120,
                timer: 0
            },
            resetPassword: {
                inputHash: false,
                retryInterval: 120,
                timer: 0
            },
            request: {
                loading: false,
                error: ``,
                code: 0,
            },
            screens: {
                login: 'login',
                loginByCode: 'loginByCode',
                resetPassword: 'resetPassword',
                changePassword: 'changePassword',
                authenticated: 'authenticated'
            },
            screenDefault: 'loginByCode',
            screen: 'loginByCode',
            authenticated: {
                user: {
                    displayName: ``,
                    type: ``,
                    email: ``,
                    phone: ``,
                    phoneFormatted: ``
                },
                session: {
                    until: ``,
                    untilFormatted: ``,
                    ip: ``,
                    userAgent: ``,
                }
            },
            input: {
                username: ``,
                password: ``,
                remember: true,
                code: ``,
                passwordResetHash: ``,
                oldPassword: ``,
                newPassword: ``,
            },
            validation: {
                enabled: false,
                failed: false,
                errors: {
                    username: ``,
                    password: ``,
                    oldPassword: ``,
                    newPassword: ``,
                    code: ``,
                    passwordResetHash: ``
                },
                rules: {
                    username: {min: 8, max: 255},
                    password: {min: 8, max: 255},
                    oldPassword: {min: 8, max: 255},
                    newPassword: {min: 8, max: 255},
                    code: {min: 3, max: 8},
                    passwordResetHash: {min: 3, max: 255},
                },
                messages: {
                    username: {
                        min: `Логин не может быть короче 8 символов`,
                        max: `Логин не может быть длиннее 255 символов`,
                    },
                    password: {
                        min: `Пароль должен содержать как минимум 8 символов`,
                        max: `Пароль не может содержать больше 255 символов`,
                    },
                    oldPassword: {
                        min: `Пароль должен содержать как минимум 8 символов`,
                        max: `Пароль не может содержать больше 255 символов`,
                    },
                    newPassword: {
                        min: `Пароль должен содержать как минимум 8 символов`,
                        max: `Пароль не может содержать больше 255 символов`,
                    },
                    code: {
                        min: `Код не может быть короче 3 символов`,
                        max: `Код не может быть длиннее 8 символов`,
                    },
                    passwordResetHash: {
                        min: `Код сброса не может быть короче 3 символов`,
                        max: `Код сброса не может быть длиннее 255 символов`,
                    }
                }
            }
        }
    },
    computed: {
        toastClasses() {
            return {
                'toast-warning': this.request.code >= 400 && this.request.code < 500,
                'toast-error': this.request.code >= 500
            }
        }
    },
    methods: {
        formClasses(property) {
            return {
                'has-error': this.validation.enabled && this.validation.errors[property],
                'has-success': this.validation.enabled && !this.validation.errors[property]
            }
        },
        reload() {
            window.location.reload()
        },
        formatTimer(timer) {
            return String(Math.floor(timer / 60)).padStart(2, '0') + ':' + String(timer % 60).padStart(2, '0')
        },
        formatDate(timestamp) {
            return (new Date(Date.parse(timestamp))).toLocaleString()
        },
        onInput() {
            const screenToValidation = {
                [this.screens.login]: this.loginForm,
                [this.screens.loginByCode]: this.generateCodeForm, // primary validation for username before login
                [this.screens.resetPassword]: this.resetPasswordForm,
                [this.screens.changePassword]: this.changePasswordForm,
            }
            if (this.loginByCode.inputCode) {
                screenToValidation[this.screens.loginByCode] = this.loginByCodeForm
            }
            if (this.resetPassword.inputHash) {
                screenToValidation[this.screens.resetPassword] = this.newPasswordForm
            }
            if (this.validation.enabled) {
                const form = screenToValidation[this.screen]
                this.validation.failed = !this.validate(form())
            }
        },
        loginForm() {
            return {
                username: this.input.username,
                password: this.input.password,
                remember: this.input.remember,
            }
        },
        generateCodeForm() {
            return {
                username: this.input.username,
            }
        },
        resetPasswordForm() {
            return {
                username: this.input.username,
            }
        },
        newPasswordForm() {
            return {
                username: this.input.username,
                password: this.input.password,
                passwordResetHash: this.input.passwordResetHash,
            }
        },
        changePasswordForm() {
            return {
                username: this.input.username,
                newPassword: this.input.newPassword,
                oldPassword: this.input.oldPassword,
            }
        },
        loginByCodeForm() {
            return {
                username: this.input.username,
                code: this.input.code,
                remember: this.input.remember,
            }
        },
        async requestPlain(url, data, method = 'GET') {
            const options = {
                method: method,
                headers: {
                    'Content-Type': 'application/json'
                },
            }
            if (!['GET', 'HEAD'].includes(method)) {
                options.body = JSON.stringify(data)
            }
            return fetch(url, options)
                .then(response => response.text())
                .catch(error => {
                    console.error(error)
                    this.request.error = String(error)
                })
        },
        async requestJSON(url, data, method = 'GET') {
            const options = {
                method: method,
                headers: {
                    'Accept': 'application/json',
                    'Content-Type': 'application/json'
                },
            }
            if (!['GET', 'HEAD'].includes(method)) {
                options.body = JSON.stringify(data)
            }
            return fetch(url, options)
                .then(response => response.json())
                .catch(error => {
                    console.error(error)
                    this.request.error = String(error)
                })
        },
        requestCheck(token) {
            return this.requestJSON(`/api/1/auth/check?authToken=${token}`, {}, 'GET')
        },
        requestLogin(data) {
            return this.requestJSON(`/api/1/auth/login`, data, 'POST')
        },
        requestLoginByCode(data) {
            return this.requestJSON(`/api/1/auth/loginByCode`, data, 'POST')
        },
        requestGenerateCode(data) {
            return this.requestPlain(`/api/1/auth/generateCode`, data, 'POST')
        },
        requestResetPassword(data) {
            return this.requestPlain(`/api/1/auth/resetPassword`, data, 'POST')
        },
        requestNewPassword(data) {
            return this.requestPlain(`/api/1/auth/newPassword`, data, 'POST')
        },
        requestChangePassword(data) {
            return this.requestPlain(`/api/1/auth/changePassword`, data, 'POST')
        },
        async requestWrapper(method, form) {
            this.validation.enabled = true

            if (!this.validate(form)) {
                this.onInput()
                return null
            }

            this.request.error = ``
            this.request.loading = true
            let response = await method(form)
            this.request.loading = false

            if (response === ``) {
                return true // case for '204 No Content'
            }

            if (typeof response === `string` && response.startsWith(`{`)) {
                response = JSON.parse(response)
            }

            if (response.code) {
                this.request.code = response.code
            }

            if (response.code >= 400 && response.code < 500) {
                const validation = (response.metadata || {})
                for (let property of Object.keys(form)) {
                    if (validation[property]) {
                        this.validation.errors[property] = validation[property]
                    }
                }
                if (Object.keys(validation).length === 0) {
                    this.validation.enabled = false
                    this.request.error = response.message
                    return null
                }
                return null
            }

            if (response.code >= 500) {
                this.validation.enabled = false
                this.request.error = response.message
                return null
            }
            return response
        },
        validate(form) {
            for (const property of Object.keys(form)) {
                let fail = {}
                if (this.validation.rules[property]) {
                    const rules = this.validation.rules[property]
                    if (rules.hasOwnProperty(`min`) && form[property].length < rules[`min`]) {
                        fail[property] = 'min'
                    }
                    if (rules.hasOwnProperty(`max`) && form[property].length > rules[`max`]) {
                        fail[property] = 'max'
                    }
                }
                if (Object.keys(fail).length !== 0) {
                    const rule = fail[property]
                    this.validation.errors[property] = this.validation.messages[property][rule]
                } else {
                    this.validation.errors[property] = ``
                }
            }
            this.validation.failed = Object.values(this.validation.errors).some(value => !!value)
            return !this.validation.failed
        },
        handleResponseError(response) {
            if (response && typeof response === `string` && response.length !== 0) {
                const e = JSON.parse(response)
                if (e && e.message) {
                    console.error(`request failed`, e)
                    this.request.error = e.message
                }
                return true
            }
            return false
        },
        async onSubmitGenerateCode() {
            const response = await this.requestWrapper(this.requestGenerateCode, this.generateCodeForm())

            if (this.handleResponseError(response)) {
                return
            }

            if (response) {
                this.validation.enabled = false
                this.loginByCode.inputCode = true
                this.loginByCode.timer = this.loginByCode.retryInterval

                const interval = setInterval(() => {
                    this.loginByCode.timer--
                    if (this.loginByCode.timer <= 0) {
                        clearInterval(interval)
                    }
                }, 1000)
            }
        },
        async authenticate(token) {
            localStorage.setItem(`authToken`, token)
            const response = await this.requestCheck(token)
            if (response && response.user && response.session) {
                this.setAuthenticationFromCheckResponse(response)
                this.screen = this.screens.authenticated
            } else {
                localStorage.removeItem(`authToken`)
                this.request.error = (response || {}).message || `Неизвестная ошибка. Попробуйте авторизоваться ещё раз`
            }
        },
        async onSubmitLoginByCode() {
            this.validation.enabled = true
            const response = await this.requestWrapper(this.requestLoginByCode, this.loginByCodeForm())
            if (response) {
                await this.authenticate(response.token)
            }
        },
        async onSubmitLogin() {
            this.validation.enabled = true
            const response = await this.requestWrapper(this.requestLogin, this.loginForm())
            if (response) {
                await this.authenticate(response.token)
            }
        },
        async onSubmitResetPassword() {
            const response = await this.requestWrapper(this.requestResetPassword, this.resetPasswordForm())

            if (this.handleResponseError(response)) {
                return
            }

            if (response) {
                this.validation.enabled = false
                this.resetPassword.inputHash = true
                this.resetPassword.timer = this.resetPassword.retryInterval

                const interval = setInterval(() => {
                    this.resetPassword.timer--
                    if (this.resetPassword.timer <= 0) {
                        clearInterval(interval)
                    }
                }, 1000)
            }
        },
        async onSubmitNewPassword() {
            const response = await this.requestWrapper(this.requestNewPassword, this.newPasswordForm())

            if (this.handleResponseError(response)) {
                return
            }

            if (response) {
                this.validation.enabled = false
                this.screen = this.screens.login
            }
        },
        async onSubmitChangePassword() {
            const response = await this.requestWrapper(this.requestChangePassword, this.changePasswordForm())

            if (this.handleResponseError(response)) {
                return
            }

            if (response) {
                this.validation.enabled = false
                this.input.password = this.input.newPassword
                this.screen = this.screens.login
            }
        },
        setAuthenticationFromCheckResponse(response) {
            this.authenticated.user.displayName = response.user.displayName
            this.authenticated.user.type = response.user.type
            this.authenticated.user.email = response.user.email || ``
            this.authenticated.user.phone = response.user.phone || ``
            this.authenticated.user.phoneFormatted = response.user.phone ? `+7` + response.user.phone : ``
            this.authenticated.session.until = response.session.until
            this.authenticated.session.untilFormatted = this.formatDate(response.session.until)
            this.authenticated.session.userAgent = response.session.userAgent
            this.authenticated.session.ip = response.session.ip
        },
        signOut() {
            localStorage.removeItem(`authToken`)
            window.location.reload()
        }
    },
}).mount('#app')
